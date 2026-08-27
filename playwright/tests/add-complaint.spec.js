import { test, expect } from "@playwright/test";

test.describe("Add a complaint", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/add-complaint?id=800&case=lpa");
  });

  test("adds a complaint to the case", async ({ page }) => {
    await expect(page.locator("text=Add Complaint")).toBeVisible();

    await expect(page.locator("text=LPA 7000-0000-0000")).toBeVisible();

    await expect(page.locator(".moj-alert")).not.toBeVisible();

    await page.getByRole('radio', { name: 'Complaint', exact: true }).click();

    await page.locator("#f-investigatingOfficer").fill("Test Officer");

    await page.locator("#f-complainantName").fill("Someones name");

    await page.locator("#f-title").fill("A title");

    await page.locator("#f-description").fill("A description");

    await page.locator("#f-receivedDate").fill("2022-04-05");

    await page.getByRole('radio', { name: 'OPG Decisions', exact: true }).click();

    await page.locator("#f-subCategory-02").selectOption("Fee Decision");

    await page.locator("#f-complainantCategory").selectOption("LPA Donor");

    await page.locator("#f-origin").selectOption("Phone call");

    await page.locator("button[type=submit]").click();

    await expect(page.locator(".moj-alert")).toBeVisible();
  });
});

