package atsemail

import (
	"fmt"
	"html"
	"strings"

	emailsv1 "github.com/crewlinker/atsemail/emails/v1"
)

// CustomTemplateParams holds the dynamic variables for rendering a custom email template.
type CustomTemplateParams struct {
	CandidateName string
	JobTitle      string
	CompanyName   string
}

// CustomTemplateResult contains the fully rendered email output ready to be sent.
type CustomTemplateResult struct {
	Subject  string
	HTMLBody string
	TextBody string
}

// RenderCustomTemplate substitutes {candidate_name}, {job_title}, and {company_name} in the
// subject and body templates, then wraps the resolved body in minimal branded HTML. The plain-text
// body is returned as-is (after substitution) for the text/plain MIME part.
func RenderCustomTemplate(
	subjectTemplate, bodyTemplate string,
	params CustomTemplateParams,
	theme *emailsv1.ThemeOverwrites,
) CustomTemplateResult {
	r := strings.NewReplacer(
		"{candidate_name}", params.CandidateName,
		"{job_title}", params.JobTitle,
		"{company_name}", params.CompanyName,
	)

	subject := r.Replace(subjectTemplate)
	textBody := r.Replace(bodyTemplate)
	htmlBody := renderCustomEmailHTML(textBody, theme)

	return CustomTemplateResult{
		Subject:  subject,
		HTMLBody: htmlBody,
		TextBody: textBody,
	}
}

// renderCustomEmailHTML wraps a plain-text body in a minimal branded HTML email.
// Newlines become <br>, double newlines become paragraph breaks. The accent bar at
// the top of the email uses the org's button background color when available.
func renderCustomEmailHTML(textBody string, theme *emailsv1.ThemeOverwrites) string {
	accentColor := "#4F46E5"
	if theme != nil && theme.GetButtonBackgroundColor() != nil {
		buttonColor := theme.GetButtonBackgroundColor()

		accentColor = fmt.Sprintf(
			"rgb(%d,%d,%d)",
			buttonColor.GetRed(),
			buttonColor.GetGreen(),
			buttonColor.GetBlue(),
		)
	}

	escaped := html.EscapeString(textBody)
	// double newline → paragraph break; single newline → line break
	escaped = strings.ReplaceAll(escaped, "\n\n", "</p><p style=\"margin:0 0 16px 0\">")
	escaped = strings.ReplaceAll(escaped, "\n", "<br>")

	return fmt.Sprintf(
		`
			<!DOCTYPE html>
			<html lang="en">
				<head>
				  <meta charset="UTF-8">
				  <meta name="viewport" content="width=device-width,initial-scale=1">
				</head>
				<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,sans-serif">
				  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:24px 0">
				    <tr>
				      <td align="center">
				        <table
				          width="600"
				          cellpadding="0"
				          cellspacing="0"
				          style="max-width:600px;width:100%%;background:#ffffff;border-radius:6px;overflow:hidden"
				        >
				          <tr>
				            <td style="background-color:%s;height:4px;font-size:0;line-height:0">&nbsp;</td>
				          </tr>
				          <tr>
				            <td style="padding:32px 40px;color:#333333;font-size:15px;line-height:1.6">
				              <p style="margin:0 0 16px 0">%s</p>
				            </td>
				          </tr>
				        </table>
				      </td>
				    </tr>
				  </table>
				</body>
			</html>
		`,
		accentColor,
		escaped,
	)
}
