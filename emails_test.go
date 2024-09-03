package atsemail_test

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	emailsv1 "github.com/crewlinker/atsemail/emails/v1"
)

func TestRenderJobApplicationNotification(t *testing.T) {
	t.Parallel()

	for idx, entry := range []struct {
		data *emailsv1.JobApplicationNotification
		exp  func(Gomega, *bytes.Buffer, *bytes.Buffer)
	}{
		{
			data: &emailsv1.JobApplicationNotification{
				JobApplicantGivenName:  "Elon",
				JobApplicantFamilyName: "Musk",
				JobPostingTitle:        "Janitor",
				JobPostingHref:         "http://dash.sterndesk.com/posting",
				JobApplicationHref:     "http://dash.sterndesk.com/application",
				OrganizationName:       "Sterndesk",
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
	tosLogo, _ := base64.StdEncoding.DecodeString("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")

	for idx, entry := range []struct {
		data *emailsv1.JobApplicationConfirm
		exp  func(Gomega, *bytes.Buffer, *bytes.Buffer)
	}{
		{
			data: &emailsv1.JobApplicationConfirm{
				JobPostingTitle:        "Janitor",
				JobPostingHref:         "http://demo.site.test.sterndesk.com/job-posting/1123",
				CareerSiteHomepageHref: "http://demo.site.test.sterndesk.com",
				OrganizationName:       "Sterndesk",
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
		{
			data: &emailsv1.JobApplicationConfirm{
				JobPostingTitle:        "Themed",
				JobPostingHref:         "http://tos.site.test.sterndesk.com/job-posting/1123",
				CareerSiteHomepageHref: "http://tos.site.test.sterndesk.com",
				OrganizationName:       "TOS Crew and Ship Delivery",

				ThemeOverwrites: &emailsv1.ThemeOverwrites{
					BorderRadius: emailsv1.BorderRadius_BORDER_RADIUS_NONE,
					ButtonBackgroundColor: &emailsv1.Color{
						Red:   0,
						Green: 30,
						Blue:  56,
					},
					LinkTextColor: &emailsv1.Color{
						Red:   0,
						Green: 91,
						Blue:  169,
					},
					HeadingImage: &emailsv1.Image{
						ContentType: "image/png",
						Data:        tosLogo,
					},
				},
			},
			exp: func(g Gomega, htbuf, txtbuf *bytes.Buffer) {
			},
		},
		{
			data: &emailsv1.JobApplicationConfirm{
				JobPostingTitle:        "Themed",
				JobPostingHref:         "http://demo.site.test.sterndesk.com/job-posting/1123",
				CareerSiteHomepageHref: "http://demo.site.test.sterndesk.com",
				OrganizationName:       "Sterndesk",
				ThemeOverwrites: &emailsv1.ThemeOverwrites{
					BorderRadius: emailsv1.BorderRadius_BORDER_RADIUS_SMALL,
					ButtonBackgroundColor: &emailsv1.Color{
						Red:   230,
						Green: 230,
						Blue:  230,
					},
					ButtonTextColor: &emailsv1.Color{
						Red:   100,
						Green: 100,
						Blue:  100,
					},
				},
			},
			exp: func(g Gomega, htbuf, txtbuf *bytes.Buffer) {
			},
		},
		{
			data: &emailsv1.JobApplicationConfirm{
				JobPostingTitle:        "Themed",
				JobPostingHref:         "http://demo.site.test.sterndesk.com/job-posting/1123",
				CareerSiteHomepageHref: "http://demo.site.test.sterndesk.com",
				OrganizationName:       "Sterndesk",
				ThemeOverwrites: &emailsv1.ThemeOverwrites{
					BorderRadius: emailsv1.BorderRadius_BORDER_RADIUS_MEDIUM,
					LinkTextColor: &emailsv1.Color{
						Red:   255,
						Green: 0,
						Blue:  0,
					},
				},
			},
			exp: func(g Gomega, htbuf, txtbuf *bytes.Buffer) {
			},
		},
		{
			data: &emailsv1.JobApplicationConfirm{
				JobPostingTitle:        "Themed",
				JobPostingHref:         "http://demo.site.test.sterndesk.com/job-posting/1123",
				CareerSiteHomepageHref: "http://demo.site.test.sterndesk.com",
				OrganizationName:       "Sterndesk",
				ThemeOverwrites: &emailsv1.ThemeOverwrites{
					BorderRadius: emailsv1.BorderRadius_BORDER_RADIUS_LARGE,
				},
			},
			exp: func(g Gomega, htbuf, txtbuf *bytes.Buffer) {
			},
		},
	} {
		t.Run(fmt.Sprintf("example %d", idx), func(t *testing.T) {
			t.Parallel()
			AssertEmailRender(t, "job-application-confirm", idx, entry.data, entry.exp)
		})
	}
}
