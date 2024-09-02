package atsemail_test

import (
	"bytes"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/crewlinker/atsemail"
)

func TestRenderJobApplicationNotification(t *testing.T) {
	t.Parallel()

	for idx, entry := range []struct {
		data atsemail.JobApplicationNotification
		exp  func(Gomega, *bytes.Buffer, *bytes.Buffer)
	}{
		{
			data: atsemail.JobApplicationNotification{
				JobApplicantGivenName:  "Elon",
				JobApplicantFamilyName: "Musk",
				JobPostingTitle:        "Janitor",
				JobPostingHref:         "http://dash.sterndesk.com/posting",
				JobApplicationHref:     "http://dash.sterndesk.com/application",
			},
			exp: func(g Gomega, htbuf, txtbuf *bytes.Buffer) {
				g.Expect(htbuf.String()).To(HavePrefix("<!DOCTYPE"))
				g.Expect(htbuf.String()).To(ContainSubstring("Elon"))
				g.Expect(htbuf.String()).To(ContainSubstring("Musk"))
				g.Expect(htbuf.String()).To(ContainSubstring("Janitor"))

				g.Expect(txtbuf.String()).To(ContainSubstring("---"))
				g.Expect(txtbuf.String()).To(ContainSubstring("Elon"))
				g.Expect(txtbuf.String()).To(ContainSubstring("Musk"))
				g.Expect(txtbuf.String()).To(ContainSubstring("Janitor"))
			},
		},
	} {
		t.Run(fmt.Sprintf("example %d", idx), func(t *testing.T) {
			t.Parallel()
			AssertEmailRender(t, "job-application-notification", idx, entry.data, entry.exp)
		})
	}
}

func TestRenderJobApplicationConfirm(t *testing.T) {
	t.Parallel()

	for idx, entry := range []struct {
		data atsemail.JobApplicationConfirm
		exp  func(Gomega, *bytes.Buffer, *bytes.Buffer)
	}{
		{
			data: atsemail.JobApplicationConfirm{
				JobPostingTitle:        "Janitor",
				CareerSiteHomepageHref: "http://demo.site.test.sterndesk.com",
			},
			exp: func(g Gomega, htbuf, txtbuf *bytes.Buffer) {
				g.Expect(htbuf.String()).To(HavePrefix("<!DOCTYPE"))
				g.Expect(htbuf.String()).To(ContainSubstring("Janitor"))
				g.Expect(htbuf.String()).To(ContainSubstring("demo.site.test.sterndesk.com"))

				g.Expect(txtbuf.String()).To(ContainSubstring("---"))
				g.Expect(txtbuf.String()).To(ContainSubstring("Janitor"))
				g.Expect(txtbuf.String()).To(ContainSubstring("demo.site.test.sterndesk.com"))
			},
		},
	} {
		t.Run(fmt.Sprintf("example %d", idx), func(t *testing.T) {
			t.Parallel()
			AssertEmailRender(t, "job-application-confirm", idx, entry.data, entry.exp)
		})
	}
}
