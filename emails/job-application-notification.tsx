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
          <Container className="border border-solid border-[#eaeaea] rounded-2xl my-[40px] mx-auto p-[20px] max-w-[465px]">
            <Heading as="h2">New job application</Heading>
            <Text>
              <em>$.JobApplicantGivenName$ $.JobApplicantFamilyName$</em> has
              applied for the job opening:&nbsp;
              <Link href="$.JobPostingHref$">$.JobPostingTitle$</Link>. Open
              the&nbsp;<b>Sterndesk</b> dashboard or view the application
              directly by clicking the button below.
            </Text>
            <Button
              className="bg-blue-500 hover:bg-blue-700 mb-5 text-white font-bold py-2 px-4 rounded"
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
