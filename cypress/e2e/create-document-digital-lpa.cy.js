describe("Create a document for a digital LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-GDJ7-QK9R-4XVF", "GET", {
      status: 200,
      body: {
        "opg.poas.sirius": {
          id: 483,
          application: {
            donorFirstNames: "Steven",
            donorLastName: "Munnell",
            donorDob: "17/06/1982",
            donorAddress: {
              addressLine1: "1 Scotland Street",
              postcode: "EH6 18J",
            },
          },
        },
      },
    });

    cy.addMock("/lpa-api/v1/templates/digitallpa", "GET", {
      status: 200,
      body: {
        DD: {
          label: "DLPA Example Form",
          inserts: {
            all: {
              DLPA_INSERT_01: {
                label: "DLPA Insert 1",
                order: 0,
              },
            },
          },
        },
      },
    });

    cy.addMock("/lpa-api/v1/lpas/483/documents/draft", "POST", {
      status: 201,
      body: {},
    });

    cy.visit("/lpa/M-GDJ7-QK9R-4XVF/documents/new");
  });

  it("creates a document on the case", () => {
    cy.contains("Select a document template");
    cy.get("#f-templateId").type("DLPA");
    cy.get(".autocomplete__menu").contains("DLPA Example Form").click();

    cy.contains("Select document inserts");
    cy.contains("DLPA Insert 1").click();

    cy.contains("1 Scotland Street, EH6 18J");
    cy.contains("Steven Munnell (Donor)").click();

    cy.contains("button", "Continue").click();
  });
});
