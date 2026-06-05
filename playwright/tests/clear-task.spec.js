import { test, expect } from "@playwright/test";

test.describe("Clear task on a digital LPA", () => {
  test.beforeEach(async ({ page, context }) => {
    // Mock the API endpoint
    await context.routeFromHAR("pacts/sirius-lpa-frontend-sirius.json", {
      url: /\/lpa-api\/v1\/tasks\/990\/mark-as-completed/,
      update: "recordings",
    });

    // Navigate to the page
    await page.goto("http://localhost:8888/clear-task?id=990");
  });

  test("marks a task as completed", async ({ page }) => {
    // Assert "Clear Task" heading is visible
    await expect(page.locator("text=Clear Task")).toBeVisible();

    // Assert moj-alert does not exist initially
    await expect(page.locator(".moj-alert")).not.toBeVisible();

    // Assert "Task:" text is visible
    await expect(page.locator("text=Task:")).toBeVisible();

    // Click the submit button
    await page.locator("button[type=submit]").click();

    // Assert moj-alert is now visible after clicking
    await expect(page.locator(".moj-alert")).toBeVisible();
  });
});
