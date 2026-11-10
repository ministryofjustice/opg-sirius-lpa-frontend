describe("create an LPA", () => {
  const assertTrimmedText = (selector, expected) => {
    cy.get(selector)
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal(expected));
  };

  const assertActorDisplay = (updateButtonId, name, hiddenText, title) => {
    cy.get(`#${updateButtonId}`)
      .closest(".app-create-lpa__details-row")
      .within(() => {
        assertTrimmedText(".govuk-details__summary-text", name);
        cy.get(".govuk-details__summary").click();
        cy.contains("th", "Title").siblings("td").should("have.text", title);
      });

    assertTrimmedText(`#${updateButtonId} .govuk-visually-hidden`, hiddenText);
  };

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

    assertActorDisplay(
      "f-update-attorney-21",
      "Active Attorney",
      "attorney Active Attorney",
      "Mr",
    );
    assertActorDisplay(
      "f-update-replacementAttorney-22",
      "Replacement Attorney",
      "replacement attorney Replacement Attorney",
      "Mrs",
    );

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

    cy.get(".govuk-details__summary-text")
      .contains("Active Attorney")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Active Attorney"));
    cy.contains(
      "#f-update-attorney-21 .govuk-visually-hidden",
      "attorney Mr Active Attorney",
    );

    cy.get(".govuk-details__summary-text")
      .contains("Replacement Attorney")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Replacement Attorney"));
    cy.contains(
      "#f-update-replacementAttorney-22 .govuk-visually-hidden",
      "replacement attorney Mrs Replacement Attorney",
    );

    // section 2
    cy.get("#f-attorneyActDecisions").should("be.checked");
    cy.get("#f-preferencesAndInstructions-3").should("be.checked");

    assertActorDisplay(
      "f-update-notifiedPerson-24",
      "Notified Person",
      "notified person Notified Person",
      "Ms",
    );

    // section 3 tested separately in create-attorney and create-correspondent
    cy.get("#f-lpaDonorSignatureDate").type("2026-06-19");

    // section 4
    cy.get("#f-applicantType").click();
    cy.get("#f-applicantSignatureDate").type("2026-06-19");

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

    assertActorDisplay(
      "f-update-attorney-21",
      "Active Attorney",
      "attorney Active Attorney",
      "Mr",
    );
    assertActorDisplay(
      "f-update-replacementAttorney-22",
      "Replacement Attorney",
      "replacement attorney Replacement Attorney",
      "Mrs",
    );

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

    cy.get(".govuk-details__summary-text")
      .contains("Active Attorney")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Active Attorney"));
    cy.contains(
      "#f-update-attorney-21 .govuk-visually-hidden",
      "attorney Mr Active Attorney",
    );

    cy.get(".govuk-details__summary-text")
      .contains("Replacement Attorney")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Replacement Attorney"));
    cy.contains(
      "#f-update-replacementAttorney-22 .govuk-visually-hidden",
      "replacement attorney Mrs Replacement Attorney",
    );

    // section 2
    cy.get("#f-attorneyActDecisions").should("be.checked");
    cy.get("#f-preferencesAndInstructions-3").should("be.checked");

    cy.get(".govuk-details__summary-text")
      .contains("Notified Person")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Notified Person"));
    cy.contains(
      "#f-update-notifiedPerson-24 .govuk-visually-hidden",
      "notified person Ms Notified Person",
    );

    // section 3 tested separately in create-attorney and create-correspondent
    cy.get("#f-lpaDonorSignatureDate").should("have.value", "2026-06-19");

    assertActorDisplay(
      "f-update-certificate-provider-23",
      "Certificate Provider",
      "certificate provider Certificate Provider",
      "Dr",
    );

    // section 4
    cy.get("#f-applicantType").click();
    cy.get("#f-applicantSignatureDate").type("2026-06-19");

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

    assertActorDisplay(
      "f-update-attorney-21",
      "Active Attorney",
      "attorney Active Attorney",
      "Mr",
    );
    assertActorDisplay(
      "f-update-replacementAttorney-22",
      "Replacement Attorney",
      "replacement attorney Replacement Attorney",
      "Mrs",
    );

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

    cy.get(".govuk-details__summary-text")
      .contains("Active Attorney")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Active Attorney"));
    cy.contains(
      "#f-update-attorney-21 .govuk-visually-hidden",
      "attorney Mr Active Attorney",
    );

    cy.get(".govuk-details__summary-text")
      .contains("Replacement Attorney")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Replacement Attorney"));
    cy.contains(
      "#f-update-replacementAttorney-22 .govuk-visually-hidden",
      "replacement attorney Mrs Replacement Attorney",
    );

    // section 2
    cy.get("#f-attorneyActDecisions").should("be.checked");
    cy.get("#f-preferencesAndInstructions-3").should("be.checked");

    assertActorDisplay(
      "f-update-notifiedPerson-24",
      "Notified Person",
      "notified person Notified Person",
      "Ms",
    );

    // section 3 tested separately in create-attorney and create-correspondent
    cy.get("#f-lpaDonorSignatureDate").type("2026-06-19");

    // section 4
    cy.get("#f-applicantType").click();
    cy.get("#f-applicantSignatureDate").type("2026-06-19");

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

    assertActorDisplay(
      "f-update-attorney-21",
      "Active Attorney",
      "attorney Active Attorney",
      "Mr",
    );
    assertActorDisplay(
      "f-update-replacementAttorney-22",
      "Replacement Attorney",
      "replacement attorney Replacement Attorney",
      "Mrs",
    );

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

    cy.get(".govuk-details__summary-text")
      .contains("Active Attorney")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Active Attorney"));
    cy.contains(
      "#f-update-attorney-21 .govuk-visually-hidden",
      "attorney Mr Active Attorney",
    );

    cy.get(".govuk-details__summary-text")
      .contains("Replacement Attorney")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Replacement Attorney"));
    cy.contains(
      "#f-update-replacementAttorney-22 .govuk-visually-hidden",
      "replacement attorney Mrs Replacement Attorney",
    );

    // section 2
    cy.get("#f-attorneyActDecisions").should("be.checked");
    cy.get("#f-preferencesAndInstructions-3").should("be.checked");

    cy.get(".govuk-details__summary-text")
      .contains("Notified Person")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Notified Person"));
    cy.contains(
      "#f-update-notifiedPerson-24 .govuk-visually-hidden",
      "notified person Ms Notified Person",
    );

    // section 3 tested separately in create-attorney and create-correspondent
    cy.get("#f-lpaDonorSignatureDate").should("have.value", "2026-06-19");

    cy.get(".govuk-details__summary-text")
      .contains("Certificate Provider")
      .invoke("text")
      .then((text) => expect(text.trim()).to.equal("Certificate Provider"));
    cy.contains(
      "#f-update-certificate-provider-23 .govuk-visually-hidden",
      "certificate provider Dr Certificate Provider",
    );

    // section 4
    cy.get("#f-applicantType").should("be.checked");
    cy.get("#f-applicantSignatureDate").should("have.value", "2026-06-19");
    cy.get("#f-applicationFee").should("be.checked");
    cy.get("#f-cardPaymentContact").should("have.value", "07700 900000");

    assertTrimmedText('label[for="f-applicantIds-21"]', "Active Attorney");

    assertActorDisplay(
      "f-update-correspondent-25",
      "Correspondent Correspondent",
      "correspondent Correspondent Correspondent",
      "Mr",
    );

    cy.contains("button", "Save and exit").click();
  });
});
