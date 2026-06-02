import { render } from "@react-email/render";
import { JobApplicationConfirm } from "../emails/job-application-confirm.tsx";
import { JobApplicationDecline } from "../emails/job-application-decline.tsx";

const templates: Record<string, React.FC<Record<string, unknown>>> = {
  "job-application-confirm": JobApplicationConfirm as React.FC<Record<string, unknown>>,
  "job-application-decline": JobApplicationDecline as React.FC<Record<string, unknown>>,
};

async function main() {
  const [, , name] = process.argv;
  if (!name || !(name in templates)) {
    process.stderr.write(`Unknown template: ${name}\n`);
    process.exit(1);
  }

  const data: Record<string, unknown> = JSON.parse(
    await new Promise<string>((res) => {
      let buf = "";
      process.stdin.on("data", (d) => (buf += d));
      process.stdin.on("end", () => res(buf));
    }),
  );

  const Component = templates[name]!;
  const html = await render(<Component {...data} />);
  const text = await render(<Component {...data} />, { plainText: true });
  process.stdout.write(JSON.stringify({ html, text }));
}

main();
