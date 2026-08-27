import { test, expect } from "@playwright/test";

test.describe("Add a complaint", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the page
    await page.goto("/add-complaint?id=800&case=lpa");
  });

  test("adds a complaint to the case", async ({ page }) => {
    // Assert "Add Complaint" heading is visible
    await expect(page.locator("text=Add Complaint")).toBeVisible();

    // Assert "LPA 7000-0000-0000" is visible
    await expect(page.locator("text=LPA 7000-0000-0000")).toBeVisible();

    // Assert moj-alert does not exist initially
    await expect(page.locator(".moj-alert")).not.toBeVisible();

    // Click the "Complaint" radio button
    await page.getByRole('radio', { name: 'Complaint', exact: true }).click();

    // Fill in the investigating officer
    await page.locator("#f-investigatingOfficer").fill("Test Officer");

    // Fill in the complainant name
    await page.locator("#f-complainantName").fill("Someones name");

    // Fill in the title
    await page.locator("#f-title").fill("A title");

    // Fill in the description
    await page.locator("#f-description").fill("A description");

    // Fill in the received date
    await page.locator("#f-receivedDate").fill("2022-04-05");

    // Click the "OPG Decisions" radio button
    await page.getByRole('radio', { name: 'OPG Decisions', exact: true }).click();

    // Select "Fee Decision" from sub-category dropdown
    await page.locator("#f-subCategory-02").selectOption("Fee Decision");

    // Select "LPA Donor" from complainant category dropdown
    await page.locator("#f-complainantCategory").selectOption("LPA Donor");

    // Select "Phone call" from origin dropdown
    await page.locator("#f-origin").selectOption("Phone call");

    // Click the submit button
    await page.locator("button[type=submit]").click();

    // Assert moj-alert is now visible after clicking
    await expect(page.locator(".moj-alert")).toBeVisible();
  });
});

