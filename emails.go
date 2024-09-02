package atsemail

type JobApplicationNotification struct {
	JobApplicantGivenName  string
	JobApplicantFamilyName string
	JobPostingTitle        string
	JobPostingHref         string
	JobApplicationHref     string
}

type JobApplicationConfirm struct {
	JobPostingTitle        string
	JobPostingHref         string
	CareerSiteHomepageHref string
}
