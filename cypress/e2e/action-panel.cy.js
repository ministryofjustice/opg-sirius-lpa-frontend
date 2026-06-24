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

    // needed for the header bar
    cy.addMock("/lpa-api/v1/persons/1", "GET", {
      status: 200,
      body: {},
    });

    cy.addMock("/lpa-api/v1/permissions", "GET", {
      status: 200,
      body: {},
    });

    cy.addMock("/lpa-api/v1/lpas/34/draft-count", "GET", {
      status: 200,
      body: {
        draftCount: 1,
      },
    });

    cy.visit("/donor/1/documents?uid[]=7000-1234-1234");
  });

  it("can open and close the action panel", () => {
    cy.get("#actions-content").should("be.visible");

    cy.get("#actions-tab").click();
    cy.get("#actions-content").should("not.be.visible");
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
});
