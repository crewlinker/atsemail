package atsemail

type JobApplicationNotification struct {
	JobApplicantGivenName  string
	JobApplicantFamilyName string
	JobPostingTitle        string
	JobPostingHref         string
	JobApplicationHref     string
}

func (JobApplicationNotification) Name() string {
	return "job-application-notification"
}
