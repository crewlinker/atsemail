package atsemail

// DefaultConfirmBodyHTML is the default body_html for job-application-confirm emails.
// It uses $.job_posting_href$, $.job_title$, and $.career_site_homepage_href$ placeholders
// resolved via resolveBodyVars at render time.
const DefaultConfirmBodyHTML = `<hr>` +
	`<h2>Application received</h2>` +
	`<p>We have succesfully received your application for the job posting:&nbsp;` +
	`<a href="$.job_posting_href$">$.job_title$</a>. This is what will happen next:</p>` +
	`<table>` +
	`<tr>` +
	`<td align="center" class="pl-2" width="20">&bull;</td>` +
	`<td class="pl-2"><p>We will review your resume and any other documents you&#39;ve provided</p></td>` +
	`</tr>` +
	`<tr>` +
	`<td align="center" class="pl-2" width="20">&bull;</td>` +
	`<td class="pl-2"><p>You&#39;ll hear back from us in the coming days</p></td>` +
	`</tr>` +
	`</table>` +
	`<p>If you want to take a look at some more job postings. You can click the button below or copy it in your browser:</p>` +
	`<a href="$.career_site_homepage_href$" class="bg-blue-500 hover:bg-blue-700 mb-5 text-white font-bold py-2 px-4 rounded sd-theme-button">View other job postings</a>` +
	`<hr>`

// DefaultDeclineBodyHTML is the default body_html for job-application-decline emails.
// It uses $.job_posting_href$, $.job_title$, and $.career_site_homepage_href$ placeholders
// resolved via resolveBodyVars at render time.
const DefaultDeclineBodyHTML = `<hr>` +
	`<h2>Update on your application</h2>` +
	`<p>We&#39;ve reviewed your application for ` +
	`<a href="$.job_posting_href$">$.job_title$</a> and, ` +
	`unfortunately, you won&#39;t be moving forward for this role. ` +
	`We appreciate the time you took to apply and encourage you to explore ` +
	`other opportunities that might be a great fit for you.</p>` +
	`<a href="$.career_site_homepage_href$" class="bg-blue-500 hover:bg-blue-700 mb-5 text-white font-bold py-2 px-4 rounded sd-theme-button">View Job Postings</a>`
