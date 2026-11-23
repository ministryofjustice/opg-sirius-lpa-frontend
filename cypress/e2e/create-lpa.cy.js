describe("create an LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/1", "GET", {
      status: 200,
      body: {
        id: 1,
        firstname: "Wilma",
        surname: "Bird",
      },
    });

    cy.addMock("/lpa-api/v1/permissions", "GET", {
      status: 200,
      body: {
        "v1-lpas-edit-dates": {
          permissions: ["PUT"],
        },
      },
    });
  });

  it("creates an LPA", () => {
    cy.addMock("/lpa-api/v1/donors/1/lpas", "POST", {
      status: 201,
      body: {},
    });

    cy.visit("/create-lpa?id=1");
    cy.get(".govuk-accordion__show-all").click();

    // section 1
    cy.get("#f-caseSubtype").click();
    cy.get("#f-applicationType").click();
    cy.get("#f-onlineLpaId").type("A12345678901");
    cy.get("#f-receiptDate").type("2026-06-19");
    cy.get("#f-applicationFee").click();
    cy.get("#f-cardPaymentContact").type("07700 900000");
    cy.get("#f-caseAttorney").click();

    // section 2
    cy.get("#f-attorneyActDecisions").click();
    cy.get("#f-preferencesAndInstructions-3").click();

    // section 3 tested separately in create-attorney and create-correspondent
    cy.get("#f-lpaDonorSignatureDate").type("2026-06-19");

    // section 4
    cy.get("#f-applicantType").click();
    cy.get("#f-applicantSignatureDate").type("2026-06-19");

    cy.contains("button", "Save and exit").click();
  });
});
