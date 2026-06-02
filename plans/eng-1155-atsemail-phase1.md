# Phase 1 — atsemail changes

> Repo: `/Users/development/Documents/work/atsemail`
> Branch: create from `main`, e.g. `romansytsykhovskyi/eng-1155-atsemail-editable-body`

---

## Overview

Make `JobApplicationConfirm` and `JobApplicationDecline` accept a dynamic
TipTap/ProseMirror JSON body and a candidate name. `EmailEditor` (from
`@react-email/editor`) renders the TipTap JSON to email-safe HTML as part of
the template. Variable placeholders (`$.candidate_name$`, `$.job_title$`,
`$.company_name$`) are resolved in `render.go` before the JSON reaches the
template.

### Architectural shift for these two templates

Currently `render.go` embeds the **exported static HTML** files and uses Go's
`html/template` package with `$...$` delimiters to substitute proto fields at
render time.

`EmailEditor` requires real TipTap JSON at render time — the JSON cannot be a
raw string placeholder inside static HTML. Therefore, `JobApplicationConfirm`
and `JobApplicationDecline` are moved to **dynamic rendering**: `render.go`
calls `@react-email/render`'s Node.js renderer at request time, passing the
proto data as a JSON argument.

The other two templates (`job-application-notification`,
`document-request-notification`) are **not changed** — they keep the static
export + Go template approach.

---

## Step 1 — Upgrade npm packages

`EmailEditor` ships in a separate package that requires react-email ≥ 6.

```bash
cd /Users/development/Documents/work/atsemail
npm install react-email@6.5.0 @react-email/components@latest @react-email/editor@latest
```

Verify `package.json` shows:

```json
"react-email": "6.5.0",
"@react-email/components": "1.0.12",   // or latest
"@react-email/editor": "1.5.3"         // or latest
```

Run `npm run check` to confirm types still pass after the upgrade.
Breaking changes between react-email v3 and v6: the CLI flags for `export`
changed; adjust `package.json` scripts if the export command errors.

---

## Step 2 — `emails/v1/emails.proto`

Add two fields to **both** `JobApplicationConfirm` and `JobApplicationDecline`:

```protobuf
// TipTap/ProseMirror JSON body; $.candidate_name$, $.job_title$,
// $.company_name$ placeholders are resolved by render.go before rendering.
string body_json = 5 [(buf.validate.field).string.min_len = 1];

// Candidate's full name, substituted for $.candidate_name$ in body_json.
// May be empty (no min_len constraint).
string candidate_name = 6;
```

> `$.job_title$` maps to existing field 1 (`job_posting_title`).
> `$.company_name$` maps to existing field 4 (`organization_name`).

Regenerate the Go bindings:

```bash
buf generate
```

Verify that `GetBodyJson()` and `GetCandidateName()` are generated on both
message types in `emails/v1/emails.pb.go`.

---

## Step 3 — `render.go`: resolve `$.var$` placeholders before rendering

Before passing `body_json` to the TSX template, `render.go` replaces the
`$.var$` placeholders with the actual proto field values. This keeps all
variable substitution inside `atsemail` — callers (atsback) pass raw values.

### New interface

```go
// BodyVarsProvider is implemented by email data types whose body_json
// may contain $.candidate_name$, $.job_title$, $.company_name$ placeholders.
type BodyVarsProvider interface {
    GetCandidateName()  string
    GetJobPostingTitle() string
    GetOrganizationName() string
    GetBodyJson()        string
}
```

### Helper function

```go
// resolveBodyVars resolves $.candidate_name$, $.job_title$, $.company_name$
// placeholders in body_json using Go's text/template with custom delimiters.
func resolveBodyVars(vars BodyVarsProvider) (string, error) {
    tmpl, err := template.New("body").
        Delims("$.", "$").
        Funcs(template.FuncMap{
            "candidate_name": func() string { return vars.GetCandidateName() },
            "job_title":      func() string { return vars.GetJobPostingTitle() },
            "company_name":   func() string { return vars.GetOrganizationName() },
        }).
        Parse(vars.GetBodyJson())
    if err != nil {
        return "", fmt.Errorf("failed to parse body_json template: %w", err)
    }

    var buf strings.Builder
    if err := tmpl.Execute(&buf, nil); err != nil {
        return "", fmt.Errorf("failed to execute body_json template: %w", err)
    }
    return buf.String(), nil
}
```

> Uses Go's `text/template` with `$.` / `$` delimiters. Each placeholder
> (e.g. `$.candidate_name$`) becomes a template function call resolved via
> `FuncMap`. This is consistent with the master plan's approach and is more
> extensible than `strings.ReplaceAll` — adding a new variable requires only
> a new `FuncMap` entry.

---

## Step 4 — `render.go`: dynamic rendering for confirm / decline

### How dynamic rendering works

`@react-email/render` exports a `render(Component, props)` function that
renders a React component to an HTML string and a plain-text string entirely
in Node.js, with no browser needed. `render.go` calls a small Node.js script
that imports the TSX module and calls `render(...)` with the data passed on
stdin as JSON, printing the results on stdout.

Add a new `RenderDynamic[E]` type (or extend `Render[E]`) that:

1. Serializes the proto message to JSON.
2. Substitutes `$.var$` placeholders via `resolveBodyVars`.
3. Spawns `node scripts/render.mjs <template-name>` (or equivalent), writes
   the serialized data to stdin, reads HTML + text from stdout.
4. Applies `ApplyTheme` to the returned HTML buffer.

#### `scripts/render.mjs` (new file)

```js
import { render } from "@react-email/render";
import { JobApplicationConfirm } from "../emails/job-application-confirm.tsx";
import { JobApplicationDecline } from "../emails/job-application-decline.tsx";

const templates = {
  "job-application-confirm": JobApplicationConfirm,
  "job-application-decline": JobApplicationDecline,
};

const [, , name] = process.argv;
const data = JSON.parse(
  await new Promise((res) => {
    let buf = "";
    process.stdin.on("data", (d) => (buf += d));
    process.stdin.on("end", () => res(buf));
  }),
);

const Component = templates[name];
const html = await render(<Component {...data} />);
const text = await render(<Component {...data} />, { plainText: true });
process.stdout.write(JSON.stringify({ html, text }));
```

> The exact implementation (subprocess vs. long-lived sidecar, JSX transform,
> etc.) depends on how the project runs Go tests in CI. The simplest starting
> point is `exec.Command("node", ...)` in `render.go`.

#### `render.go` dynamic render path

```go
type RenderDynamic[E EmailData] struct {
    name string
}

func NewDynamic[E EmailData](name string) *RenderDynamic[E] {
    return &RenderDynamic[E]{name: name}
}

func (r *RenderDynamic[E]) Render(val *protovalidate.Validator, txtw, htmw io.Writer, data E) error {
    if err := val.Validate(data); err != nil {
        return fmt.Errorf("invalid email data: %w", err)
    }

    // Resolve $.var$ placeholders in body_json before rendering.
    payload := buildPayload(data)   // proto → map[string]any with resolved body_json

    b, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal payload: %w", err)
    }

    cmd := exec.Command("node", "scripts/render.mjs", r.name)
    cmd.Stdin = bytes.NewReader(b)
    out, err := cmd.Output()
    if err != nil {
        return fmt.Errorf("node render failed: %w", err)
    }

    var result struct {
        HTML string `json:"html"`
        Text string `json:"text"`
    }
    if err := json.Unmarshal(out, &result); err != nil {
        return fmt.Errorf("failed to parse render result: %w", err)
    }

    if _, err := io.WriteString(txtw, result.Text); err != nil {
        return fmt.Errorf("failed to write text: %w", err)
    }

    var htmBuf bytes.Buffer
    if _, err := htmBuf.WriteString(result.HTML); err != nil {
        return fmt.Errorf("failed to write html: %w", err)
    }

    if theme := data.GetThemeOverwrites(); theme != nil {
        if err := r.ApplyTheme(&htmBuf, theme); err != nil {
            return fmt.Errorf("failed to apply theme: %w", err)
        }
    }

    _, err = io.Copy(htmw, &htmBuf)
    return err
}
```

`buildPayload` converts the proto message to a plain `map[string]any` with
camelCase keys matching the TSX component props, and substitutes `$.var$`
placeholders in `body_json` via `resolveBodyVars`.

---

## Step 5 — TSX templates

### `emails/job-application-confirm.tsx`

Remove all hardcoded body content (Heading, Text, bullet Section, Button).
Import `EmailEditor` and render the TipTap JSON from the `bodyJson` prop.

```tsx
import { Body, Container, Head, Hr, Html, Img, Preview, Tailwind } from "@react-email/components";
import { EmailEditor } from "@react-email/editor";

interface Props {
  jobPostingTitle: string;
  jobPostingHref: string;
  careerSiteHomepageHref: string;
  organizationName: string;
  bodyJson: string; // TipTap JSON; $.var$ already resolved by render.go
  candidateName: string;
}

export const JobApplicationConfirm = ({ jobPostingTitle, organizationName, bodyJson }: Props) => {
  const doc = JSON.parse(bodyJson);

  return (
    <Html>
      <Head />
      <Preview>Application received for {jobPostingTitle}</Preview>
      <Tailwind>
        <Body className="bg-white my-auto mx-auto font-sans px-2">
          <Container className="border border-solid border-[#eaeaea] rounded-2xl my-[40px] mx-auto p-[20px] max-w-[465px] sd-theme-container">
            <Img
              src={`data:image/png;base64,<EXISTING_LOGO_BASE64>`}
              className="mb-5 sd-theme-heading-image"
              height="75"
              alt={organizationName}
            />
            <Hr />
            <EmailEditor content={doc} editable={false} />
            <Hr />
            <p style={{ color: "rgb(156,163,175)", fontSize: "14px", lineHeight: "24px" }}>
              If you didn&#39;t apply for this job posting, it could be that someone applied with the wrong email. You
              can ignore this email.
            </p>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};

export default JobApplicationConfirm;
```

Keep the existing logo `base64` string verbatim (don't re-encode).

### `emails/job-application-decline.tsx`

Same pattern. Remove hardcoded Heading, Text, Button. Use `EmailEditor`.

```tsx
import { Body, Container, Head, Hr, Html, Img, Preview, Tailwind } from "@react-email/components";
import { EmailEditor } from "@react-email/editor";

interface Props {
  jobPostingTitle: string;
  jobPostingHref: string;
  careerSiteHomepageHref: string;
  organizationName: string;
  bodyJson: string;
  candidateName: string;
}

export const JobApplicationDecline = ({ jobPostingTitle, organizationName, bodyJson }: Props) => {
  const doc = JSON.parse(bodyJson);

  return (
    <Html>
      <Head />
      <Preview>Update on your application – {jobPostingTitle}</Preview>
      <Tailwind>
        <Body className="bg-white my-auto mx-auto font-sans px-2">
          <Container className="border border-solid border-[#eaeaea] rounded-2xl my-[40px] mx-auto p-[20px] max-w-[465px] sd-theme-container">
            <Img
              src={`data:image/png;base64,<EXISTING_LOGO_BASE64>`}
              className="mb-5 sd-theme-heading-image"
              height="75"
              alt={organizationName}
            />
            <Hr />
            <EmailEditor content={doc} editable={false} />
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};

export default JobApplicationDecline;
```

> **Note on the import path**: `EmailEditor` lives in the `@react-email/editor`
> package, not in the `react-email` CLI package. Import from
> `@react-email/editor`. If the master plan referenced `from "react-email"`,
> that import path does not exist in either v3 or v6 — use
> `from "@react-email/editor"` instead.

---

## Step 6 — `npm run export`

The two changed templates no longer participate in the static export (they are
rendered dynamically). Run export anyway to update the other two templates and
to verify the TSX compiles without type errors:

```bash
npm run export
```

If the export fails because `email export` tries to export all TSX files (and
the two new ones fail without real TipTap JSON), exclude them via the export
config, or provide a wrapper with hardcoded default JSON for export-only use.

Check `exported/html/job-application-confirm.html` and `…-decline.html` — these
files will still be overwritten by `npm run export` but are no longer embedded
by `render.go`. They serve as static previews only.

---

## Step 7 — Update `emails_test.go`

### Shared body constant

```go
const testBodyJSON = `{"type":"doc","content":[{"type":"heading","attrs":{"level":2},"content":[{"type":"text","text":"Test heading"}]},{"type":"paragraph","content":[{"type":"text","text":"Hello "},{"type":"text","marks":[{"type":"bold"}],"text":"$.candidate_name$"},{"type":"text","text":", your application for "},{"type":"text","marks":[{"type":"italic"}],"text":"$.job_title$"},{"type":"text","text":" at $.company_name$ has been processed."}]}]}`
```

Both confirm and decline tests use this constant — satisfying the requirement
that "two screenshots use the same body".

### `TestRenderJobApplicationConfirm` — update all cases

Every existing test case adds `BodyJson` and `CandidateName`. Case 0 gets
assertions; the rest produce screenshots only:

```go
data: &emailsv1.JobApplicationConfirm{
    JobPostingTitle:        "Janitor",
    JobPostingHref:         "http://demo.site.test.sterndesk.com/job-posting/1123",
    CareerSiteHomepageHref: "http://demo.site.test.sterndesk.com",
    OrganizationName:       "Sterndesk",
    BodyJson:               testBodyJSON,
    CandidateName:          "Jane Doe",
},
exp: func(g Gomega, htbuf, txtbuf *bytes.Buffer) {
    g.Expect(htbuf.String()).To(HavePrefix("<!DOCTYPE"))
    g.Expect(htbuf.String()).To(ContainSubstring("Test heading"))
    g.Expect(htbuf.String()).To(ContainSubstring("Jane Doe"))   // $.candidate_name$ resolved
    g.Expect(htbuf.String()).To(ContainSubstring("Janitor"))    // $.job_title$ resolved
    g.Expect(htbuf.String()).To(ContainSubstring("Sterndesk"))  // $.company_name$ resolved
    g.Expect(htbuf.String()).ToNot(ContainSubstring("$.candidate_name$"))
},
```

### `TestRenderJobApplicationDecline` — update both cases

Add `BodyJson: testBodyJSON, CandidateName: "Jane Doe"` to both cases.
Case 1 (the simple, non-themed case) gets the same assertions as confirm
case 0 above.

### Note on `AssertEmailRender` / `render_test.go`

`AssertEmailRender` is generic and calls `atsemail.New[T](templateName)`.
Once `render.go` introduces `RenderDynamic`, update `AssertEmailRender` to
accept either `Render[T]` or `RenderDynamic[T]` (e.g., via an interface), or
add a parallel `AssertDynamicEmailRender` helper. The screenshot-saving logic
in `SaveScreenshot` does not change.

---

## Step 8 — Run tests to verify

```bash
cd /Users/development/Documents/work/atsemail
go test ./... -v -count=1
```

Expected outcome:

- All tests pass.
- `screenshots/job-application-confirm_0.png` — branded shell + `EmailEditor`
  rendering "Test heading / Hello Jane Doe, your application for Janitor at
  Sterndesk has been processed."
- `screenshots/job-application-decline_1.png` — same body, decline shell.
- No literal `$.candidate_name$` visible in any screenshot.

---

## Step 9 — Publish new version

```bash
git tag v0.0.7
git push origin v0.0.7
```

Then in atsback:

```bash
go get github.com/crewlinker/atsemail@v0.0.7
go mod tidy
```

---

## File change summary

| File                                         | Change                                                                                      |
| -------------------------------------------- | ------------------------------------------------------------------------------------------- |
| `package.json`                               | Upgrade `react-email` → 6.5.0, add `@react-email/editor`, upgrade `@react-email/components` |
| `emails/v1/emails.proto`                     | Add `body_json` (field 5) + `candidate_name` (field 6) to both messages                     |
| `emails/v1/emails.pb.go`                     | Auto-generated by `buf generate`                                                            |
| `emails/job-application-confirm.tsx`         | Replace hardcoded body with `<EmailEditor content={doc} editable={false} />`                |
| `emails/job-application-decline.tsx`         | Same                                                                                        |
| `exported/html/job-application-confirm.html` | Re-exported (static preview only, no longer embedded)                                       |
| `exported/html/job-application-decline.html` | Same                                                                                        |
| `exported/text/job-application-confirm.txt`  | Re-exported                                                                                 |
| `exported/text/job-application-decline.txt`  | Re-exported                                                                                 |
| `scripts/render.mjs`                         | New — Node.js renderer script called by `render.go`                                         |
| `render.go`                                  | Add `BodyVarsProvider` interface, `resolveBodyVars`, `RenderDynamic[E]` type                |
| `emails_test.go`                             | Add `BodyJson` + `CandidateName` to all confirm/decline test cases                          |
| `render_test.go`                             | Update `AssertEmailRender` to handle dynamic render path                                    |

---

## Open questions

1. **Node.js availability in CI**: The dynamic rendering path spawns `node`.
   Confirm that the CI environment has Node.js available (it currently runs
   `npm run export`, so it should). If not, a long-lived sidecar approach
   (start `node scripts/render.mjs` once per test run) is more efficient.

2. **`npm run export` for confirm/decline**: If react-email's `email export`
   command fails for templates that require dynamic data (because it calls
   them without props), either skip exporting those two files or provide
   a wrapper component with hardcoded default content for export-only use.

3. **Text email for dynamic templates**: The Node.js renderer can produce
   plain text via `render(Component, { plainText: true })`. Confirm this
   output is acceptable (TipTap's text extraction vs. current hand-written
   template).
