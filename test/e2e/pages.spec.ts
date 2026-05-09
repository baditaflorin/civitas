import { Buffer } from "node:buffer";
import { expect, test } from "@playwright/test";

const apiBase = process.env.PLAYWRIGHT_API_BASE_URL ?? "http://127.0.0.1:18089";

test("loads the published investigation workspace", async ({ page }) => {
  await page.addInitScript((value) => {
    localStorage.setItem("civitas_api_base_url", value);
  }, apiBase);
  await page.goto(".");
  await expect(page.getByRole("heading", { name: "Investigation OS" })).toBeVisible();
  await expect(page.getByRole("link", { name: /star on github/i })).toHaveAttribute(
    "href",
    "https://github.com/baditaflorin/civitas"
  );
  await expect(page.getByRole("link", { name: /paypal/i })).toHaveAttribute(
    "href",
    "https://www.paypal.com/paypalme/florinbadita"
  );
  await expect(page.getByLabel("Evidence relationship map")).toBeVisible();
});

test("fresh user creates a case, uploads pasted evidence, and exports", async ({
  page,
}) => {
  await page.addInitScript((value) => {
    localStorage.setItem("civitas_api_base_url", value);
  }, apiBase);
  await page.goto(".");

  await expect(page.getByText(/Backend 0\.2\.0/)).toBeVisible();
  await page.getByPlaceholder("Case title").fill("Stranger workflow");
  await page.getByPlaceholder("Case notes").fill("Smoke test case");
  await page.getByRole("button", { name: /new case/i }).click();
  await expect(page.getByRole("button", { name: /stranger workflow/i })).toBeVisible();

  await page
    .getByPlaceholder("Paste text or HTML evidence")
    .fill("Contract signed on 2026-05-10 by source@example.org for EUR 1200.");
  await page.getByRole("button", { name: /upload paste/i }).click();
  await expect(page.getByText(/Uploaded 1\/1 file/)).toBeVisible();
  await expect(
    page.getByRole("article").getByText("pasted-evidence.txt")
  ).toBeVisible();

  await page.getByRole("button", { name: /^safe export/i }).click();
  await expect(page.getByText(/schema_version: phase2\.export\.v1/)).toBeVisible();
  await expect(page.getByRole("button", { name: /copy export/i })).toBeVisible();
  await expect(page.getByRole("button", { name: /download markdown/i })).toBeVisible();
  await expect(page.getByRole("button", { name: /download state/i })).toBeEnabled();

  const stateSnippet = await page
    .locator("pre")
    .filter({ hasText: "/state > civitas-case-state.json" })
    .textContent();
  const caseId = stateSnippet?.match(/cases\/([^/]+)\/state/)?.[1];
  expect(caseId).toBeTruthy();
  const stateResponse = await page.request.get(
    `${apiBase}/api/v1/cases/${caseId}/state`,
  );
  expect(stateResponse.ok()).toBeTruthy();
  const state = await stateResponse.text();
  const fileChooserPromise = page.waitForEvent("filechooser");
  await page.getByText("Import state").click();
  const fileChooser = await fileChooserPromise;
  await fileChooser.setFiles({
    name: "stranger-workflow-state.json",
    mimeType: "application/json",
    buffer: Buffer.from(state),
  });
  await expect(page.getByText(/Imported Stranger workflow/)).toBeVisible();
  await expect(
    page.getByRole("button", { name: /Stranger workflow \(imported\)/i }),
  ).toBeVisible();
});
