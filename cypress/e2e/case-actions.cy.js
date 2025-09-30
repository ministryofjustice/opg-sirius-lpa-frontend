import * as digitalLpas from "../mocks/digitalLpas";
import * as cases from "../mocks/cases";

describe("Case actions drop down", () => {
  beforeEach(() => {
    const mocks = Promise.allSettled([
      cases.tasks.empty("1111"),
      cases.warnings.empty("1111"),
      digitalLpas.get("M-1111-1111-1111"),
      digitalLpas.objections.empty("M-1111-1111-1111"),
      digitalLpas.progressIndicators.feesInProgress("M-1111-1111-1111"),
    ]);

    cy.wrap(mocks);

    cy.addMock("/lpa-api/v1/persons/1111/cases", "GET", {
      status: 200,
      body: {
        cases: [
          {
            caseSubtype: "property-and-affairs",
            id: 1111,
            uId: "M-1111-1111-1111",
            status: "Draft",
            caseType: "DIGITAL_LPA",
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/cases/1111", "GET", {
      status: 200,
      body: {
        id: 1111,
        uId: "M-1111-1111-1111",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 1111,
        },
      },
    });

    cy.visit("/lpa/M-1111-1111-1111");
  });

  it("can create a task", () => {
    cy.addMock("/lpa-api/v1/cases/1111/tasks", "POST", {
      status: 201,
      body: { tasks: [] },
    });
    cy.contains(".govuk-button", "Case actions").click();
    cy.contains("Create a task").click();
    cy.url().should("include", "/create-task?id=1111");
    cy.contains("M-1111-1111-1111");
    cy.get("#f-taskType").select("Check Application");
    cy.get("#f-name").type("Do this task");
    cy.get("#f-description").type("This task, do");
    cy.contains("label", "Team").click();
    cy.get("#f-assigneeTeam").select("Cool Team");
    cy.get("#f-dueDate").type("2035-01-01");
    cy.get("button[type=submit]").click();

    cy.get(".moj-alert").should("exist");
    cy.get(".moj-alert").contains("Task created");
    cy.get("h1").contains("Steven Munnell");
    cy.location("pathname").should("eq", "/lpa/M-1111-1111-1111");
  });

  it("can cancel changing the status", () => {
    cy.addMock("/lpa-api/v1/reference-data/caseChangeReason", "GET", {
      status: 200,
      body: [
        {
          handle: "LPA_DOES_NOT_WORK",
          label: "The LPA does not work and cannot be changed",
          parentSources: ["cannot-register"],
        },
      ],
    });
    cy.contains("Case actions").click();
    cy.contains("Change case status").click();

    cy.url().should("include", "/change-case-status?uid=M-1111-1111-1111");
    cy.contains("Change case status");
    cy.get(".govuk-button-group").contains("Cancel").click();

    cy.url().should("include", "/lpa/M-1111-1111-1111");
    cy.contains("Case summary");
  });

  it("can cancel creating a warning", () => {
    cy.contains("Case actions").click();
    cy.contains("Create a warning").click();

    cy.url().should("include", "/create-warning?id=1111");
    cy.contains("Create Warning");
    cy.contains("Cancel").click();

    cy.url().should("include", "/lpa/M-1111-1111-1111");
    cy.contains("Case summary");
  });

  it("creates a warning via case actions", () => {
    cy.contains(".govuk-button", "Case actions").click();
    cy.contains("Create a warning").click();
    cy.url().should("include", "/create-warning?id=1111");
    cy.get("#f-warningType").select("Complaint Received");
    cy.get("#f-warningText").type("Be warned!");
    cy.get("button[type=submit]").click();

    cy.get(".moj-alert").should("exist");
    cy.get(".moj-alert").contains("Warning created");
    cy.get("h1").contains("Steven Munnell");
    cy.location("pathname").should("eq", "/lpa/M-1111-1111-1111");
  });
});
