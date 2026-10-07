describe("Case info panel on the header bar", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/123", "GET", {
      status: 200,
      body: {
        uId: "7000-0000-0123",
        applicationType: "Online",
        caseSubType: "hw",
        assignee: {
          id: 99,
          displayName: "Sarah Jones",
          phoneNumber: "03004560300",
        },
        applicants: [
          {
            id: 1,
            firstname: "Melanie",
            surname: "Vanvolkenburg",
          },
        ],
        receiptDate: "21/06/2026",
        filingDate: "22/06/2026",
        registrationDate: "23/06/2026",
        dispatchDate: "24/06/2026",
        lpaDonorSignatureDate: "17/06/2026",
        caseAttorneySingular: true,
        caseAttorneyJointly: false,
        caseAttorneyJointlyAndSeverally: false,
        caseAttorneyJointlyAndJointlyAndSeverally: false,
        lifeSustainingTreatment: "Option A",
        batchId: "123",
      },
    });

    cy.visit("/sirius-header-case-info?id=123");
  });

  it("displays the case info panel", () => {
    cy.contains("Case owner").should("exist");
    cy.contains("Sarah Jones").should("exist");
    cy.contains("03004560300").should("exist");

    cy.contains("Case ID").should("exist");
    cy.contains("7000-0000-0123").should("exist");

    cy.contains("Receipt date").should("exist");
    cy.contains("21/06/2026").should("exist");

    cy.contains("Filing date").should("exist");
    cy.contains("22/06/2026").should("exist");

    cy.contains("Reg due date / Reg date").should("exist");
    cy.contains("23/06/2026").should("exist");

    cy.contains("Dispatch date").should("exist");
    cy.contains("24/06/2026").should("exist");

    cy.contains("Date donor signed").should("exist");
    cy.contains("17/06/2026").should("exist");

    cy.contains("Applicant").should("exist");
    cy.contains("Melanie Vanvolkenburg").should("exist");

    cy.contains("How attorneys are appointed").should("exist");
    cy.contains("Singular").should("exist");

    cy.contains("LST choice").should("exist");
    cy.contains("Option A").should("exist");

    cy.contains("Batch ID").should("exist");
    cy.contains("123").should("exist");
  });

  it("shows unallocated when there is no assignee", () => {
    cy.addMock("/lpa-api/v1/cases/456", "GET", {
      status: 200,
      body: {
        uId: "7000-0000-0456",
      },
    });

    cy.visit("/sirius-header-case-info?id=456");

    cy.contains("Case owner").should("exist");
    cy.contains("Unallocated").should("exist");
  });
});
