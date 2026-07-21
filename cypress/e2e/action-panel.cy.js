describe("Action Panel", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/1/cases", "GET", {
      status: 200,
      body: {
        cases: [
          {
            caseType: "LPA",
            caseSubtype: "pfa",
            id: 34,
            uId: "7000-1234-1234",
          },
          { caseType: "LPA", caseSubtype: "hw", id: 78, uId: "7000-5678-5678" },
          {
            caseType: "EPA",
            caseSubtype: "pfa",
            id: 111,
            uId: "7000-9876-5432",
          },
        ],
      },
    });

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0,case:34&limit=999",
      "GET",
      {
        status: 200,
        body: {
          total: 0,
          documents: [],
        },
      },
    );

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0,case:78&limit=999",
      "GET",
      {
        status: 200,
        body: {
          total: 0,
          documents: [],
        },
      },
    );

    // needed for the header bar
    cy.addMock("/lpa-api/v1/persons/1", "GET", {
      status: 200,
      body: {},
    });

    cy.addMock("/lpa-api/v1/permissions", "GET", {
      status: 200,
      body: {
        "v1-persons": {
          permissions: ["GET"],
        },
        "v1-persons-cases": {
          permissions: ["GET"],
        },
        "v1-cases-tasks-post": {
          permissions: ["POST"],
        },
        "v1-donors": {
          permissions: ["POST", "PUT"],
        },
        "v1-donors-epas": {
          permissions: ["POST"],
        },
        "v1-donors-lpas": {
          permissions: ["POST"],
        },
        "v1-lpas": {
          permissions: ["PUT"],
        },
        "v1-lpas-documents-draft": {
          permissions: ["POST"],
        },
        "v1-lpas-investigations": {
          permissions: ["POST"],
        },
        "v1-notes": {
          permissions: ["POST"],
        },
        "v1-payments": {
          permissions: ["GET"],
        },
        "v1-person-links": {
          permissions: ["POST", "PATCH"],
        },
        "v1-person-references": {
          permissions: ["DELETE"],
        },
        "v1-persons-references": {
          permissions: ["POST"],
        },
        "v1-poa-tasks": {
          permissions: ["PUT"],
        },
        "v1-users-updateusercases": {
          permissions: ["PUT"],
        },
        "v1-warnings": {
          permissions: ["POST"],
        },
        reporting: {
          permissions: ["GET"],
        },
      },
    });

    cy.addMock("/lpa-api/v1/lpas/34/draft-count", "GET", {
      status: 200,
      body: {
        draftCount: 1,
      },
    });

    cy.addMock("/lpa-api/v1/lpas/78/draft-count", "GET", {
      status: 200,
      body: {
        draftCount: 1,
      },
    });

    cy.addMock(
      "/lpa-api/v1/cases/34/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [
            { id: 990, name: "Review application", dueDate: "01/07/2026" },
          ],
        },
      },
    );

    cy.addMock("/lpa-api/v1/persons/1/references", "GET", {
      status: 200,
      body: [
        {
          referenceId: 123,
        },
      ],
    });

    cy.visit("/donor/1/documents?uid[]=7000-1234-1234");
  });

  it("can open and close the action panel", () => {
    cy.get("#actions-content").should("be.visible");

    cy.get("#actions-tab").click();
    cy.get("#actions-content").should("not.be.visible");
  });

  it("does not display any buttons on the action panel without permissions", () => {
    cy.addMock("/lpa-api/v1/permissions", "GET", {
      status: 200,
      body: {
        "v1-persons-cases": {
          permissions: ["GET"],
        },
      },
    });

    cy.addMock("/lpa-api/v1/lpas/78/draft-count", "GET", {
      status: 200,
      body: {
        draftCount: 0,
      },
    });

    cy.addMock(
      "/lpa-api/v1/cases/78/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.visit("/donor/1/documents?uid[]=7000-5678-5678");

    cy.contains("Add investigation").should("not.be.visible");
    cy.contains("Add complaint").should("not.be.visible");
    cy.contains("Allocate Case").should("not.be.visible");
    cy.contains("Change status").should("not.be.visible");
    cy.contains("Create document").should("not.be.visible");
    cy.contains("Create donor").should("not.be.visible");
    cy.contains("Create epa case").should("not.be.visible");
    cy.contains("Create lpa case").should("not.be.visible");
    cy.contains("Create event").should("not.be.visible");
    cy.contains("Create relationship").should("not.be.visible");
    cy.contains("Create warning").should("not.be.visible");
    cy.contains("Delete relationship").should("not.be.visible");
    cy.contains("Edit dates").should("not.be.visible");
    cy.contains("Edit donor").should("not.be.visible");
    cy.contains("Edit case").should("not.be.visible");
    cy.contains("Fees").should("not.be.visible");
    cy.contains("Unlink record").should("not.be.visible");
    cy.contains("Link record").should("not.be.visible");
    cy.contains("MI reporting").should("not.be.visible");
    cy.contains("New task").should("not.be.visible");
    cy.contains("Retrieve draft").should("not.be.visible");
  });

  it("displays the warning button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Create warning");

    cy.get("a#action-panel-button-create-warning").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create Warning");
  });

  it("displays the event button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Create event");

    cy.get("a#action-panel-button-create-event").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create Event");
  });

  it("displays the complaint button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Add complaint");

    cy.addMock("/lpa-api/v1/cases/34", "GET", { status: 200, body: {} });

    cy.get("a#action-panel-button-add-complaint").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Add Complaint");
  });

  it("displays the new task button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("New task");

    cy.addMock("/lpa-api/v1/cases/34", "GET", { status: 200, body: {} });

    cy.get("a#action-panel-button-new-task").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create Task");
  });

  it("displays the create document button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Create document");

    cy.addMock("/lpa-api/v1/cases/34", "GET", { status: 200, body: {} });
    cy.addMock("/lpa-api/v1/documents/%s", "GET", { status: 200, body: {} });
    cy.addMock("/lpa-api/v1/templates/lpa", "GET", {
      status: 200,
      body: {
        DD: {
          label: "Donor deceased: Blank template",
          inserts: {
            all: {
              DD1: {
                label: "DD1 - Case complete",
                order: 0,
              },
            },
          },
        },
      },
    });

    cy.get("a#action-panel-button-create-document").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Select a document template");
  });

  it("displays the edit document button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Retrieve draft");

    cy.addMock("/lpa-api/v1/cases/34", "GET", {
      status: 200,
      body: {
        donor: { id: 33 },
      },
    });
    cy.addMock("/lpa-api/v1/lpas/34/documents?type[]=Draft", "GET", {
      status: 200,
      body: [
        {
          id: 789,
          uuid: "6789",
        },
      ],
    });
    cy.addMock("/lpa-api/v1/documents/6789", "GET", { status: 200, body: {} });
    cy.addMock("/lpa-api/v1/templates/lpa", "GET", {
      status: 200,
      body: {
        DD: {
          label: "Donor deceased: Blank template",
          inserts: {
            all: {
              DD1: {
                label: "DD1 - Case complete",
                order: 0,
              },
            },
          },
        },
      },
    });

    cy.get("a#action-panel-button-retrieve-draft").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Edit draft document");
  });

  it("displays the change status button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Change status");

    cy.addMock("/lpa-api/v1/cases/34", "GET", {
      status: 200,
      body: { caseType: "lpa" },
    });
    cy.addMock("/lpa-api/v1/lpas/34/available-statuses", "GET", {
      status: 200,
      body: ["Cancelled", "Withdrawn"],
    });

    cy.get("a#action-panel-button-change-status").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Change Status");
  });

  it("displays fees button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Fees");

    cy.addMock("/lpa-api/v1/cases/34", "GET", {
      status: 200,
      body: {
        caseType: "LPA",
        caseSubtype: "pfa",
        donor: {
          id: 1,
        },
        id: 34,
        uId: "7000-1234-1234",
        expectedPaymentTotal: 8200,
      },
    });
    cy.addMock("/lpa-api/v1/cases/34/payments", "GET", {
      status: 200,
      body: [{ amount: 100 }],
    });
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "Reduced Fees User"],
      },
    });

    cy.get("a#action-panel-button-fees").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Fee details");

    cy.get(".action-panel__form").contains("Add payment").click();
    cy.get(".action-panel__form h1").contains("Add a payment");
    cy.get(".action-panel__form .govuk-link").contains("Cancel").click();

    cy.get(".action-panel__form").contains("Apply fee reduction").click();
    cy.get(".action-panel__form h1").contains("Apply a fee reduction");
    cy.get(".action-panel__form .govuk-link").contains("Cancel").click();
  });

  it("displays the create donor button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Create donor");

    cy.addMock("/lpa-api/v1/cases/34", "GET", { status: 200, body: {} });

    cy.get("a#action-panel-button-create-donor").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create Donor");
  });

  it("displays the edit donor button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Edit donor");

    cy.addMock("/lpa-api/v1/cases/34", "GET", { status: 200, body: {} });

    cy.get("a#action-panel-button-edit-donor").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Edit Donor");
  });

  it("displays the allocate case button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Allocate Case");

    cy.addMock("/lpa-api/v1/cases/34", "GET", { status: 200, body: {} });

    cy.get("a#action-panel-button-allocate-case").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Allocate Case");
  });

  it("displays the edit dates button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Edit dates");

    cy.addMock("/lpa-api/v1/cases/34", "GET", {
      status: 200,
      body: {
        uid: "7000-1234-1234",
        caseType: "LPA",
        donor: { id: 1 },
      },
    });

    cy.get("a#action-panel-button-edit-dates").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Edit Dates");
  });

  it("displays the MI reporting button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("MI reporting");

    cy.addMock("/lpa-api/v1/reporting/config", "GET", {
      status: 200,
      body: {},
    });

    cy.get("a#action-panel-button-mi-reporting").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("MI Reporting");
  });

  it("displays the create relationship button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Create relationship");

    cy.get("a#action-panel-button-create-relationship").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create Relationship");
  });

  it("displays the link record button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Link record");

    cy.get("a#action-panel-button-link-record").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Link record");
  });

  it("displays the delete relationship button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Delete relationship");

    cy.get("a#action-panel-button-delete-relationship").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Delete Relationship");
  });

  it("displays the add investigation button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Add investigation");

    cy.addMock("/lpa-api/v1/cases/34", "GET", { status: 200, body: {} });

    cy.get("a#action-panel-button-add-investigation").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create Investigation");
  });

  it("displays the create epa button on the action panel and can click through to subforms", () => {
    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0&limit=999",
      "GET",
      {
        status: 200,
        body: {
          total: 0,
          documents: [],
        },
      },
    );

    cy.addMock("/lpa-api/v1/donors/1/epas", "POST", {
      status: 200,
      body: { id: 123 },
    });

    cy.addMock("/lpa-api/v1/epas/123", "PUT", {
      status: 200,
      body: { id: 123 },
    });

    cy.addMock("/lpa-api/v1/cases/123", "GET", {
      status: 200,
      body: {
        id: 123,
        receiptDate: "10/06/2026",
        attorneys: [],
      },
    });

    cy.addMock("/lpa-api/v1/epas/123/attorneys", "POST", {
      status: 200,
      body: { id: 123 },
    });

    cy.visit("/donor/1/documents"); // create epa button is only clickable when no cases are selected

    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Create epa case");

    cy.get("a#action-panel-button-create-epa-case").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create an EPA");

    cy.get("#f-receiptDate").type("2026-06-19");

    cy.get(".action-panel__form").contains("Add attorney").click();
    cy.get(".action-panel__form h1").contains("Add an attorney");
    cy.get(".action-panel__form .govuk-button")
      .contains("Save and add another attorney")
      .click();

    // assert we're scrolled to the top of the new form
    cy.get(".action-panel__form h1")
      .contains("Add an attorney")
      .should("be.visible");
    cy.get(".action-panel__form .govuk-link").contains("Cancel").click();

    // assert we are back on the create epa form and scrolled to step 3
    cy.get(
      "#accordion-create-epa-heading-3 span.govuk-accordion__section-heading-text-focus",
    ).should("be.visible");

    cy.get(".action-panel__form").contains("Add correspondent").click();
    cy.get(".action-panel__form h1").contains("Add a correspondent");
    cy.get(".action-panel__form .govuk-link").contains("Cancel").click();
  });

  it("displays the create lpa button on the action panel", () => {
    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0&limit=999",
      "GET",
      {
        status: 200,
        body: {
          total: 0,
          documents: [],
        },
      },
    );

    cy.visit("/donor/1/documents"); // create lpa button is only clickable when no cases are selected

    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Create lpa case");

    cy.get("a#action-panel-button-create-lpa-case").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create an LPA");
  });

  it("displays the edit case button on the action panel", () => {
    cy.addMock("/lpa-api/v1/cases/111", "GET", {
      status: 200,
      body: { id: 111 },
    });

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0,case:111&limit=999",
      "GET",
      {
        status: 200,
        body: {
          total: 0,
          documents: [],
        },
      },
    );

    cy.addMock("/lpa-api/v1/epas/111/draft-count", "GET", {
      status: 200,
      body: {
        draftCount: 1,
      },
    });

    cy.addMock(
      "/lpa-api/v1/cases/111/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.visit("/donor/1/documents?uid[]=7000-9876-5432");

    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Edit case");

    cy.get("a#action-panel-button-edit-case").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Edit EPA");
  });

  it("displays the assign task button on the action panel", () => {
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Assign task");

    cy.addMock("/lpa-api/v1/tasks/990", "GET", {
      status: 200,
      body: {
        id: 990,
        name: "Review application",
        caseItems: [{ caseType: "LPA", uId: "7000-1234-1234" }],
      },
    });

    cy.addMock("/lpa-api/v1/teams", "GET", {
      status: 200,
      body: [{ id: 23, displayName: "Cool Team" }],
    });

    cy.get("a#action-panel-button-assign-task").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Assign Task");
  });

  it("applies large font styling to the action panel when an accessible theme is set", () => {
    cy.setCookie("siriusTheme", "accessible-light");
    cy.visit("/donor/1/documents?uid[]=7000-1234-1234");

    cy.get("html").should("have.class", "app-!-html-class--large-font");
    cy.get(".action-panel__button").should("not.be.visible");
    cy.get("a#action-panel-button-create-warning")
      .should("be.visible")
      .and("contain", "Create warning");
  });
});
