import { expect, test } from "@playwright/test";

test("loads the published investigation workspace", async ({ page }) => {
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
