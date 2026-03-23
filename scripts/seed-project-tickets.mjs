#!/usr/bin/env node

const args = process.argv.slice(2);

function getArg(name, fallback = "") {
  const idx = args.indexOf(name);
  if (idx === -1) return fallback;
  return args[idx + 1] ?? fallback;
}

function fail(message) {
  console.error(message);
  process.exit(1);
}

const baseUrl = getArg("--base-url", "http://localhost:8080/api/v1").replace(
  /\/$/,
  "",
);
const email = getArg("--email");
const password = getArg("--password");
const projectId = Number(getArg("--project-id"));
const count = Number(getArg("--count", "20"));

if (!email || !password || !Number.isInteger(projectId) || projectId <= 0) {
  fail(
    [
      "Usage:",
      "  node scripts/seed-project-tickets.mjs --project-id <id> --email <email> --password <password> [--count 20] [--base-url http://localhost:8080/api/v1]",
    ].join("\n"),
  );
}

if (!Number.isInteger(count) || count <= 0 || count > 500) {
  fail("--count must be an integer between 1 and 500");
}

async function request(path, options = {}, cookie = "") {
  const headers = new Headers(options.headers || {});
  if (!headers.has("Content-Type"))
    headers.set("Content-Type", "application/json");
  if (cookie) headers.set("Cookie", cookie);

  const response = await fetch(`${baseUrl}${path}`, {
    ...options,
    headers,
  });

  const contentType = response.headers.get("content-type") || "";
  const body = contentType.includes("application/json")
    ? await response.json().catch(() => null)
    : await response.text().catch(() => "");

  return { response, body };
}

function pick(list) {
  return list[Math.floor(Math.random() * list.length)];
}

function maybe(probability) {
  return Math.random() < probability;
}

const titleTemplates = [
  "Fix broken {area} rendering on {surface}",
  "Implement {area} filter for {surface}",
  "Refactor {area} data flow in {surface}",
  "Improve {area} latency for {surface}",
  "Add validation for {area} creation",
  "Investigate flaky {area} behavior",
  "Stabilize {area} sync edge case",
  "Document {area} lifecycle expectations",
  "Add tests for {area} state transitions",
  "Handle empty state in {surface} {area} panel",
];

const areas = [
  "ticket",
  "comment",
  "tag",
  "priority",
  "status",
  "project member",
  "session",
  "activity log",
];

const surfaces = [
  "table view",
  "kanban",
  "ticket detail",
  "dashboard",
  "api client",
  "auth flow",
];

const descriptionTemplates = [
  "Observed in staging. Reproducible after page refresh and quick navigation between project views.",
  "Happens when request fails once and client retries. Need deterministic handling.",
  "Affects internal users frequently; should be addressed before next release.",
  "Current behavior is inconsistent with specification. Align with expected project-scoped behavior.",
  "Please include regression coverage for this path and update API error copy if needed.",
  "Scope this change to the current module and avoid broad refactor in this pass.",
];

function buildTitle() {
  return pick(titleTemplates)
    .replace("{area}", pick(areas))
    .replace("{surface}", pick(surfaces));
}

function buildDescription() {
  const first = pick(descriptionTemplates);
  const second = pick(descriptionTemplates.filter((item) => item !== first));
  return `${first}\n\n${second}`;
}

function randomSubset(values, max = 2) {
  if (values.length === 0) return [];
  const copy = [...values];
  const size = Math.min(values.length, Math.floor(Math.random() * (max + 1)));
  const out = [];
  for (let i = 0; i < size; i += 1) {
    const idx = Math.floor(Math.random() * copy.length);
    out.push(copy[idx]);
    copy.splice(idx, 1);
  }
  return out;
}

async function main() {
  console.log(`Logging in as ${email}...`);
  const login = await request("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });

  if (!login.response.ok) {
    fail(
      `Login failed (${login.response.status}): ${JSON.stringify(login.body)}`,
    );
  }

  const cookieHeader = login.response.headers.get("set-cookie");
  if (!cookieHeader) {
    fail("No session cookie received from login.");
  }
  const cookie = cookieHeader.split(";")[0];

  console.log(`Loading taxonomy for project ${projectId}...`);
  const [statusesRes, prioritiesRes, tagsRes] = await Promise.all([
    request(`/projects/${projectId}/statuses`, {}, cookie),
    request(`/projects/${projectId}/priorities`, {}, cookie),
    request(`/projects/${projectId}/tags`, {}, cookie),
  ]);

  if (!statusesRes.response.ok)
    fail(`Failed statuses (${statusesRes.response.status})`);
  if (!prioritiesRes.response.ok)
    fail(`Failed priorities (${prioritiesRes.response.status})`);
  if (!tagsRes.response.ok) fail(`Failed tags (${tagsRes.response.status})`);

  const statuses = Array.isArray(statusesRes.body) ? statusesRes.body : [];
  const priorities = Array.isArray(prioritiesRes.body)
    ? prioritiesRes.body
    : [];
  const tags = Array.isArray(tagsRes.body) ? tagsRes.body : [];

  if (statuses.length === 0) fail("No statuses found for project.");
  if (priorities.length === 0) fail("No priorities found for project.");

  const statusIds = statuses
    .map((item) => item.id)
    .filter((id) => Number.isInteger(id));
  const priorityIds = priorities
    .map((item) => item.id)
    .filter((id) => Number.isInteger(id));
  const tagIds = tags
    .map((item) => item.id)
    .filter((id) => Number.isInteger(id));

  let created = 0;
  let failed = 0;

  console.log(`Creating ${count} tickets...`);
  for (let i = 0; i < count; i += 1) {
    const payload = {
      title: buildTitle(),
      description: maybe(0.9) ? buildDescription() : "",
      status_id: pick(statusIds),
      priority_id: pick(priorityIds),
      tag_ids: randomSubset(tagIds, 3),
    };

    const createRes = await request(
      `/projects/${projectId}/tickets`,
      {
        method: "POST",
        body: JSON.stringify(payload),
      },
      cookie,
    );

    if (!createRes.response.ok) {
      failed += 1;
      console.error(
        `Failed ticket ${i + 1}/${count} (${createRes.response.status}): ${JSON.stringify(createRes.body)}`,
      );
      continue;
    }

    created += 1;
    const ticketNumber =
      createRes.body && typeof createRes.body === "object"
        ? createRes.body.ticket_number
        : "?";
    console.log(`Created #${ticketNumber}: ${payload.title}`);
  }

  console.log(`Done. Created: ${created}, Failed: ${failed}`);
}

main().catch((error) => {
  fail(
    `Unexpected error: ${error instanceof Error ? error.message : String(error)}`,
  );
});
