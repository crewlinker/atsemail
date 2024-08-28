import {
  Body,
  Button,
  Container,
  Column,
  Head,
  Heading,
  Hr,
  Html,
  Img,
  Link,
  Preview,
  Row,
  Section,
  Text,
  Tailwind,
} from "@react-email/components";
import * as React from "react";

export const JobApplicationNotification = () => {
  return (
    <Html>
      <Head />
      <Preview>Preview Text</Preview>
      <Tailwind>
        <Body className="bg-white my-auto mx-auto font-sans px-2">
          <Container className="border border-solid border-[#eaeaea] rounded my-[40px] mx-auto p-[20px] max-w-[465px]">
            <Text>Hello world</Text>
          </Container>
        </Body>
      </Tailwind>
    </Html>
  );
};

export default JobApplicationNotification;
