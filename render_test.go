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
		example    atsemail.JobApplicationNotification
		expectHTML func(Gomega, *bytes.Buffer)
		expectText func(Gomega, *bytes.Buffer)
	}{
		{
			example: atsemail.JobApplicationNotification{
				JobApplicantGivenName:  "Elon",
				JobApplicantFamilyName: "Musk",
				JobPostingTitle:        "Janitor",
				JobPostingHref:         "http://dash.sterndesk.com/posting",
				JobApplicationHref:     "http://dash.sterndesk.com/application",
			},
			expectHTML: func(g Gomega, buf *bytes.Buffer) {
				g.Expect(buf.String()).To(HavePrefix("<!DOCTYPE"))
				g.Expect(buf.String()).To(ContainSubstring("Elon"))
				g.Expect(buf.String()).To(ContainSubstring("Musk"))
				g.Expect(buf.String()).To(ContainSubstring("Janitor"))
			},
			expectText: func(g Gomega, buf *bytes.Buffer) {
				g.Expect(buf.String()).To(ContainSubstring("---"))
				g.Expect(buf.String()).To(ContainSubstring("Elon"))
				g.Expect(buf.String()).To(ContainSubstring("Musk"))
				g.Expect(buf.String()).To(ContainSubstring("Janitor"))
			},
		},
	} {
		t.Run(fmt.Sprintf("example %d", idx), func(t *testing.T) {
			t.Parallel()
			g := NewWithT(t)

			render, err := atsemail.New(entry.example)
			g.Expect(err).ToNot(HaveOccurred())

			var txtbuf, htbuf bytes.Buffer
			if err := render.Render(&txtbuf, &htbuf); err != nil {
				t.Errorf("failed to render example %d: %v", idx, err)
			}

			entry.expectHTML(g, &htbuf)
			entry.expectText(g, &txtbuf)
		})
	}
}
