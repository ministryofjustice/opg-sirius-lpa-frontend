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
    cy.contains("Case owner:");
    cy.contains("Sarah Jones");
    cy.contains("03004560300");

    cy.contains("Case ID:");
    cy.contains("7000-0000-0123");

    cy.contains("Who applied to register:");
    cy.contains("Melanie Vanvolkenburg");

    cy.contains("Online or Classic application:");
    cy.contains("Online");

    cy.contains("Receipt date:");
    cy.contains("21/06/2026");

    cy.contains("Date Donor signed Instrument:");
    cy.contains("17/06/2026");

    cy.contains("Attorneys appointed:");
    cy.contains("Singular");
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

    cy.contains("Case owner:");
    cy.contains("Unallocated");
  });
});
