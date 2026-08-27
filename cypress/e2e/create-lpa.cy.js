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
    cy.get("#f-caseAttorney").click();

    // section 2
    cy.get("#f-attorneyActDecisions").click();
    cy.get("#f-preferencesAndInstructions-3").click();

    // section 3 tested separately in create-attorney and create-correspondent
    cy.get("#f-lpaDonorSignatureDate").type("2026-06-19");

    // section 4
    cy.get("#f-applicantType").click();
    cy.get("#f-applicantSignatureDate").type("2026-06-19");
    cy.get("#f-applicationFee").click();
    cy.get("#f-cardPaymentContact").type("07700 900000");
    cy.get("#f-anyOtherInfo").click();
    cy.get("#f-additionalInfo").type("None");

    cy.contains("button", "Save and exit").click();
  });

  it("displays details of an existing LPA", () => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        id: 2,
        caseSubtype: "pfa",
        applicationType: "Online",
        onlineLpaId: "A12345678901",
        receiptDate: "19/06/2026",
        caseAttorneySingular: true,
        attorneyActDecisions: "When Registered",
        applicationHasGuidance: false,
        applicationHasRestrictions: false,
        lpaDonorSignatureDate: "19/06/2026",
        applicantType: "donor",
        applicantSignatureDate: "19/06/2026",
        paymentByDebitCreditCard: true,
        cardPaymentContact: "07700 900000",
        anyOtherInfo: true,
        additionalInfo: "None",
        attorneys: [
          {
            id: 21,
            salutation: "Mr",
            firstname: "Active",
            surname: "Attorney",
            systemStatus: true,
            personType: "Attorney",
          },
        ],
        replacementAttorneys: [
          {
            id: 22,
            salutation: "Mrs",
            firstname: "Replacement",
            surname: "Attorney",
            personType: "Replacement Attorney",
          },
        ],
        certificateProviders: [
          {
            id: 23,
            salutation: "Dr",
            firstname: "Certificate",
            surname: "Provider",
            personType: "Certificate Provider",
          },
        ],
        notifiedPersons: [
          {
            id: 24,
            salutation: "Ms",
            firstname: "Notified",
            surname: "Person",
            personType: "Notified Person",
          },
        ],
        correspondent: {
          id: 25,
          salutation: "Mr",
          firstname: "Correspondent",
          surname: "Correspondent",
          personType: "Correspondent",
        },
        applicants: [
          {
            id: 1,
          },
        ],
      },
    });

    cy.visit("/create-lpa?id=1&caseId=2");
    cy.get(".govuk-accordion__show-all").click();

    // section 1
    cy.get("#f-caseSubtype").should("be.checked");
    cy.get("#f-applicationType").should("be.checked");
    cy.get("#f-onlineLpaId").should("have.value", "A12345678901");
    cy.get("#f-receiptDate-readonly").should("have.text", "19/06/2026");
    cy.get("#f-caseAttorney").should("be.checked");

    cy.contains(".govuk-details__summary-text", "Mr Active Attorney");
    cy.contains(
      "#f-update-attorney-21 .govuk-visually-hidden",
      "attorney Mr Active Attorney",
    );

    cy.contains(".govuk-details__summary-text", "Mrs Replacement Attorney");
    cy.contains(
      "#f-update-replacementAttorney-22 .govuk-visually-hidden",
      "replacement attorney Mrs Replacement Attorney",
    );

    // section 2
    cy.get("#f-attorneyActDecisions").should("be.checked");
    cy.get("#f-preferencesAndInstructions-3").should("be.checked");

    cy.contains(".govuk-details__summary-text", "Ms Notified Person");
    cy.contains(
      "#f-update-notifiedPerson-24 .govuk-visually-hidden",
      "notified person Ms Notified Person",
    );

    // section 3 tested separately in create-attorney and create-correspondent
    cy.get("#f-lpaDonorSignatureDate").should("have.value", "2026-06-19");

    cy.contains(".govuk-details__summary-text", "Dr Certificate Provider");
    cy.contains(
      "#f-update-certificate-provider-23 .govuk-visually-hidden",
      "certificate provider Dr Certificate Provider",
    );

    // section 4
    cy.get("#f-applicantType").should("be.checked");
    cy.get("#f-applicantSignatureDate").should("have.value", "2026-06-19");
    cy.get("#f-applicationFee").should("be.checked");
    cy.get("#f-cardPaymentContact").should("have.value", "07700 900000");
    cy.get("#f-anyOtherInfo").should("be.checked");
    cy.get("#f-additionalInfo").should("have.value", "None");

    cy.contains(
      ".govuk-details__summary-text",
      "Mr Correspondent Correspondent",
    );
    cy.contains(
      "#f-update-correspondent-25 .govuk-visually-hidden",
      "correspondent Mr Correspondent Correspondent",
    );

    cy.get("#f-additionalInfo").clear().type("Updated info");

    cy.contains("button", "Save and exit").click();
  });
});
