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
  Row,
  Section,
  Column,
} from "@react-email/components";
import * as React from "react";

export const JobApplicationConfirm = () => {
  return (
    <Html>
      <Head />
      <Preview>Application received for $.JobPostingTitle$</Preview>
      <Tailwind>
        <Body className="bg-white my-auto mx-auto font-sans px-2">
          <Container className="border border-solid border-[#eaeaea] rounded-2xl my-[40px] mx-auto p-[20px] max-w-[465px]">
            <Heading as="h2">Application received</Heading>
            <Text>
              We have succesfully received your application for the job
              posting:&nbsp;<a href="$.JobPostingHref$">$.JobPostingTitle$</a>.
              This is what will happen next:
            </Text>
            <Section>
              <Row>
                <Column align="center" className="pl-2" width="20">
                  &bull;
                </Column>
                <Column className="pl-2">
                  <Text>
                    We will review your resume and any other documents you've
                    provided
                  </Text>
                </Column>
              </Row>
              <Row>
                <Column align="center" className="pl-2" width="20">
                  &bull;
                </Column>
                <Column className="pl-2">
                  <Text>You'll hear back from us in the coming days</Text>
                </Column>
              </Row>
            </Section>
            <Text>
              If you want to take a look at some more job postings. You can
              click the button below or copy it in your browser:
            </Text>
            <Button
              className="bg-blue-500 hover:bg-blue-700 mb-5 text-white font-bold py-2 px-4 rounded"
              href="$.CareerSiteHomepageHref$"
            >
              View other job postings
            </Button>
            <Hr />
            <Text className="faded text-gray-400">
              If you didn't apply for this job posting, it could be that someone
              applied with the wrong email and it accidentally arrived in your
              inbox. If that is the case you can ignore this email.
            </Text>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};

export default JobApplicationConfirm;
