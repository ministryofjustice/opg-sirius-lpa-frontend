describe("Case info panel on the header bar", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/123", "GET", {
      status: 200,
      body: {
        uId: "7000-0000-0123",
        applicationType: "Online",
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
        lpaDonorSignatureDate: "17/06/2026",
        caseAttorneySingular: true,
        caseAttorneyJointly: false,
        caseAttorneyJointlyAndSeverally: false,
        caseAttorneyJointlyAndJointlyAndSeverally: false,
      },
    });

    cy.visit("/sirius-header-case-info?id=123");
  });

  it("displays the case info panel", () => {
    cy.contains("Case owner:").should("exist");
    cy.contains("Sarah Jones").should("exist");
    cy.contains("03004560300").should("exist");

    cy.contains("Case ID:").should("exist");
    cy.contains("7000-0000-0123").should("exist");

    cy.contains("Who applied to register:").should("exist");
    cy.contains("Melanie Vanvolkenburg").should("exist");

    cy.contains("Online or Classic application:").should("exist");
    cy.contains("Online").should("exist");

    cy.contains("Receipt date:").should("exist");
    cy.contains("21/06/2026").should("exist");

    cy.contains("Date Donor signed Instrument:").should("exist");
    cy.contains("17/06/2026").should("exist");

    cy.contains("Attorneys appointed:").should("exist");
    cy.contains("Singular").should("exist");
  });

  it("does not display fields with no data", () => {
    cy.contains("Notification date:").should("not.exist");
    cy.contains("Registration due date:").should("not.exist");
    cy.contains("Registration date:").should("not.exist");
    cy.contains("Dispatch date:").should("not.exist");
    cy.contains("Closed date:").should("not.exist");
    cy.contains("Attorney declaration signature date:").should("not.exist");
    cy.contains("Notice given date:").should("not.exist");
    cy.contains("Life sustaining treatment:").should("not.exist");
    cy.contains("Batch ID:").should("not.exist");
    cy.contains("CaseRec number:").should("not.exist");
  });

  it("shows unallocated when there is no assignee", () => {
    cy.addMock("/lpa-api/v1/cases/456", "GET", {
      status: 200,
      body: {
        uId: "7000-0000-0456",
      },
    });

    cy.visit("/sirius-header-case-info?id=456");

    cy.contains("Case owner:").should("exist");
    cy.contains("Unallocated").should("exist");
  });
});
