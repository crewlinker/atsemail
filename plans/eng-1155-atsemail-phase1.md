# Phase 1 — atsemail changes

> Repo: `https://github.com/crewlinker/atsemail`
> Branch: `romansytsykhovskyi/eng-1155-atsemail-editable-body`

---

## Overview

Make `JobApplicationConfirm` and `JobApplicationDecline` accept a dynamic
HTML body and a candidate name. The caller passes pre-rendered, stringified
HTML in the `body_html` field. Variable placeholders (`$.candidate_name$`,
`$.job_title$`, `$.company_name$`) are resolved in `render.go` before the
HTML is inserted into the template.

### Architectural approach

Both templates stay on the **static export + Go template** path — the same as
`job-application-notification` and `document-request-notification`. No dynamic
Node.js rendering is used.

The TSX templates expose a `$.BodyHtml$` placeholder via
`dangerouslySetInnerHTML`. `render.go` exports them as static HTML files,
embeds them, and at render time substitutes `$.BodyHtml$` (and the other
`$...$` fields) using Go's `html/template`.

---

## Step 1 — `emails/v1/emails.proto`

Add two fields to **both** `JobApplicationConfirm` and `JobApplicationDecline`:

```protobuf
// Stringified HTML body; $.candidate_name$, $.job_title$,
// $.company_name$ placeholders are resolved by render.go before rendering.
string body_html = 5 [(buf.validate.field).string.min_len = 1];

// Candidate's full name, substituted for $.candidate_name$ in body_html.
string candidate_name = 6;
```

Regenerate Go bindings:

```bash
buf generate
```

Verify that `GetBodyHtml()` and `GetCandidateName()` are generated on both
message types in `emails/v1/emails.pb.go`.

---

## Step 2 — `render.go`: resolve `$.var$` placeholders

### `BodyVarsProvider` interface

```go
type BodyVarsProvider interface {
    GetCandidateName() string
    GetJobPostingTitle() string
    GetOrganizationName() string
    GetBodyHtml() string
    GetJobPostingHref() string
    GetCareerSiteHomepageHref() string
}
```

### `resolveBodyVars`

Resolves `$.candidate_name$`, `$.job_title$`, `$.company_name$`,
`$.job_posting_href$`, `$.career_site_homepage_href$` placeholders inside
`body_html` using Go's `text/template` with `$.` / `$` delimiters.

> `html/template` does not re-execute content inserted as `htemplate.HTML`,
> so placeholders inside `body_html` would survive verbatim without this
> separate pass.

### `bodyTemplateData`

```go
type bodyTemplateData struct {
    JobPostingTitle        string
    JobPostingHref         string
    CareerSiteHomepageHref string
    OrganizationName       string
    CandidateName          string
    BodyHtml               htemplate.HTML
}
```

`BodyHtml` is typed as `htemplate.HTML` so `html/template` inserts it without
escaping.

### `Render[E].Render` method

`Render[E]` handles both plain and body-var templates via a runtime type
assertion:

```go
if bvp, ok := any(data).(BodyVarsProvider); ok {
    // resolve body vars, build bodyTemplateData, execute templates with it
} else {
    // execute templates with raw proto data
}
```

Full sequence:

1. Validates the proto message.
2. If `data` implements `BodyVarsProvider`: calls `resolveBodyVars`, builds
   `bodyTemplateData`, and executes templates against it.
3. Otherwise: executes templates against the raw proto message.
4. Applies theme overwrites.

---

## Step 3 — TSX templates

Both `emails/job-application-confirm.tsx` and
`emails/job-application-decline.tsx` replace their hardcoded body content with
a static `$.BodyHtml$` placeholder:

```tsx
<div dangerouslySetInnerHTML={{ __html: "$.BodyHtml$" }} />
```

No props are needed on the React component — all substitution happens in
`render.go` after export.

---

## Step 4 — `npm run export`

Export the updated templates to update the static HTML/text files:

```bash
npm run export
```

The exported files are embedded by `render.go` and used by `Render[E]`.

---

## Step 5 — Update `emails_test.go`

### Shared body constant

```go
const testBodyHTML = `<h2>Test heading</h2><p>Hello <strong>$.candidate_name$</strong>, your application for <em>$.job_title$</em> at $.company_name$ has been processed.</p>`
```

### Test cases

All `JobApplicationConfirm` and `JobApplicationDecline` test cases set
`BodyHtml: testBodyHTML` and `CandidateName: "Jane Doe"`. Assertions verify
that `$.candidate_name$` is not present in the output and that resolved values
appear.

---

## Step 6 — Run tests

```bash
go test ./... -v -count=1
```

---

## Step 7 — Publish new version

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

| File                                         | Change                                                                                                    |
| -------------------------------------------- | --------------------------------------------------------------------------------------------------------- |
| `emails/v1/emails.proto`                     | Add `body_html` (field 5) + `candidate_name` (field 6) to both messages                                   |
| `emails/v1/emails.pb.go`                     | Auto-generated by `buf generate`                                                                          |
| `emails/job-application-confirm.tsx`         | Replace hardcoded body with `dangerouslySetInnerHTML={{ __html: "$.BodyHtml$" }}`                         |
| `emails/job-application-decline.tsx`         | Same                                                                                                      |
| `exported/html/job-application-confirm.html` | Re-exported with `$.BodyHtml$` placeholder                                                                |
| `exported/html/job-application-decline.html` | Same                                                                                                      |
| `exported/text/job-application-confirm.txt`  | Re-exported                                                                                               |
| `exported/text/job-application-decline.txt`  | Same                                                                                                      |
| `render.go`                                  | Add `BodyVarsProvider` interface, `resolveBodyVars`, `bodyTemplateData`; extend `Render[E].Render` with runtime `BodyVarsProvider` check |
| `emails_test.go`                             | Replace `testBodyJSON` with `testBodyHTML`; add `BodyHtml` + `CandidateName` to all confirm/decline cases |
