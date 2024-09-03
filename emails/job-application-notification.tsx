import {
  Body,
  Button,
  Container,
  Head,
  Heading,
  Hr,
  Html,
  Link,
  Preview,
  Text,
  Tailwind,
  Img,
} from "@react-email/components";
import * as React from "react";

export const JobApplicationNotification = () => {
  return (
    <Html>
      <Head />
      <Preview>
        $.JobApplicantGivenName$ $.JobApplicantFamilyName$ has applied for:
        $.JobPostingTitle$
      </Preview>
      <Tailwind>
        <Body className="bg-white my-auto mx-auto font-sans px-2">
          <Container className="border border-solid border-[#eaeaea] rounded-2xl my-[40px] mx-auto p-[20px] max-w-[465px] sd-theme-container">
            <Img
              src={`data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEsAAABLCAYAAAA4TnrqAAAEtGlUWHRYTUw6Y29tLmFkb2JlLnhtcAAAAAAAPD94cGFja2V0IGJlZ2luPSLvu78iIGlkPSJXNU0wTXBDZWhpSHpyZVN6TlRjemtjOWQiPz4KPHg6eG1wbWV0YSB4bWxuczp4PSJhZG9iZTpuczptZXRhLyIgeDp4bXB0az0iWE1QIENvcmUgNS41LjAiPgogPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4KICA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0iIgogICAgeG1sbnM6dGlmZj0iaHR0cDovL25zLmFkb2JlLmNvbS90aWZmLzEuMC8iCiAgICB4bWxuczpleGlmPSJodHRwOi8vbnMuYWRvYmUuY29tL2V4aWYvMS4wLyIKICAgIHhtbG5zOnBob3Rvc2hvcD0iaHR0cDovL25zLmFkb2JlLmNvbS9waG90b3Nob3AvMS4wLyIKICAgIHhtbG5zOnhtcD0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wLyIKICAgIHhtbG5zOnhtcE1NPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvbW0vIgogICAgeG1sbnM6c3RFdnQ9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZUV2ZW50IyIKICAgdGlmZjpJbWFnZUxlbmd0aD0iNzUiCiAgIHRpZmY6SW1hZ2VXaWR0aD0iNzUiCiAgIHRpZmY6UmVzb2x1dGlvblVuaXQ9IjIiCiAgIHRpZmY6WFJlc29sdXRpb249IjcyLzEiCiAgIHRpZmY6WVJlc29sdXRpb249IjcyLzEiCiAgIGV4aWY6UGl4ZWxYRGltZW5zaW9uPSI3NSIKICAgZXhpZjpQaXhlbFlEaW1lbnNpb249Ijc1IgogICBleGlmOkNvbG9yU3BhY2U9IjEiCiAgIHBob3Rvc2hvcDpDb2xvck1vZGU9IjMiCiAgIHBob3Rvc2hvcDpJQ0NQcm9maWxlPSJzUkdCIElFQzYxOTY2LTIuMSIKICAgeG1wOk1vZGlmeURhdGU9IjIwMjQtMDkt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`}
              className="mb-5 sd-theme-heading-image"
              height="75"
              alt="$.OrganizationName$"
            />
            <Hr />
            <Heading as="h2">New job application</Heading>
            <Text>
              <em>$.JobApplicantGivenName$ $.JobApplicantFamilyName$</em> has
              applied for the job opening:&nbsp;
              <Link href="$.JobPostingHref$">$.JobPostingTitle$</Link>. Open
              the&nbsp;<b>Sterndesk</b> dashboard or view the application
              directly by clicking the button below.
            </Text>
            <Button
              className="bg-blue-500 hover:bg-blue-700 mb-5 text-white font-bold py-2 px-4 rounded sd-theme-button"
              href="$.JobApplicationHref$"
            >
              View application
            </Button>
            <Hr />
            <Text className="faded text-gray-400">
              If you don't want to receive these notifications, unassign
              yourself from the job posting or close it.
            </Text>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};

export default JobApplicationNotification;
