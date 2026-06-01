package atsemail_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/crewlinker/atsemail"
	emailsv1 "github.com/crewlinker/atsemail/emails/v1"
)

func TestRenderCustomTemplate(t *testing.T) {
	t.Parallel()

	for idx, entry := range []struct {
		subject string
		body    string
		params  atsemail.CustomTemplateParams
		theme   *emailsv1.ThemeOverwrites
		exp     func(Gomega, atsemail.CustomTemplateResult)
	}{
		{
			subject: "Thank you for applying for: {job_title}",
			body: "Dear {candidate_name},\n\n" +
				"Thank you for applying for the {job_title} position at {company_name}.\n\n" +
				"Best regards,\n{company_name}",
			params: atsemail.CustomTemplateParams{
				CandidateName: "Jane Doe",
				JobTitle:      "Software Engineer",
				CompanyName:   "Acme Corp",
			},
			exp: func(g Gomega, res atsemail.CustomTemplateResult) {
				g.Expect(res.Subject).To(Equal("Thank you for applying for: Software Engineer"))

				g.Expect(res.TextBody).To(ContainSubstring("Jane Doe"))
				g.Expect(res.TextBody).To(ContainSubstring("Software Engineer"))
				g.Expect(res.TextBody).To(ContainSubstring("Acme Corp"))

				g.Expect(res.HTMLBody).To(HavePrefix("<!DOCTYPE"))
				g.Expect(res.HTMLBody).To(ContainSubstring("Jane Doe"))
				g.Expect(res.HTMLBody).To(ContainSubstring("Software Engineer"))
				g.Expect(res.HTMLBody).To(ContainSubstring("Acme Corp"))
				// default accent color
				g.Expect(res.HTMLBody).To(ContainSubstring("#4F46E5"))
			},
		},
		{
			subject: "Update on your application – {job_title}",
			body:    "Dear {candidate_name},\n\nWe regret to inform you.",
			params: atsemail.CustomTemplateParams{
				CandidateName: "John Smith",
				JobTitle:      "Designer",
				CompanyName:   "DesignCo",
			},
			theme: &emailsv1.ThemeOverwrites{
				ButtonBackgroundColor: &emailsv1.Color{
					Red: 0, Green: 30, Blue: 56,
				},
			},
			exp: func(g Gomega, res atsemail.CustomTemplateResult) {
				g.Expect(res.Subject).To(ContainSubstring("Designer"))
				g.Expect(res.HTMLBody).To(ContainSubstring("rgb(0,30,56)"))
				g.Expect(res.HTMLBody).NotTo(ContainSubstring("#4F46E5"))
			},
		},
		{
			subject: "No variables here",
			body:    "Plain body without variables.",
			params:  atsemail.CustomTemplateParams{},
			exp: func(g Gomega, res atsemail.CustomTemplateResult) {
				g.Expect(res.Subject).To(Equal("No variables here"))
				g.Expect(res.TextBody).To(Equal("Plain body without variables."))
			},
		},
	} {
		t.Run(fmt.Sprintf("example %d", idx), func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			res := atsemail.RenderCustomTemplate(entry.subject, entry.body, entry.params, entry.theme)
			entry.exp(g, res)
		})
	}
}
