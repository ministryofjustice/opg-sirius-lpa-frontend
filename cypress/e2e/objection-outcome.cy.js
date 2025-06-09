import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Resolved objections summary", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-1005", "GET", {
      status: 200,
      body: {
        uId: "M-0000-0000-1005",
        "opg.poas.sirius": {
          id: 5,
          uId: "M-0000-0000-1005",
          status: "Draft",
          caseSubtype: "personal-welfare",
          application: {
            donorFirstNames: "James",
            donorLastName: "Rubin",
            donorDob: "22/02/1990",
            donorAddress: {
              addressLine1: "Apartment 3",
              country: "GB",
              postcode: "B15 3AA",
              town: "Birmingham",
            },
          },
        },
      },
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-1005/objections", "GET", {
      status: 200,
      body: [
        {
          id: 1005,
          notes: "test",
          objectionType: "factual",
          receivedDate: "2025-01-01",
          lpaUids: ["M-0000-0000-1005"],
          objectionLpas: [
            {
              uid: "M-0000-0000-1005",
              resolution: "upheld",
              resolutionDate: "2025-02-02",
            },
          ],
        },
      ],
    });

    cy.addMock("/lpa-api/v1/objections/1005", "GET", {
      status: 200,
      body: {
        id: 1005,
        notes: "test",
        objectionType: "factual",
        receivedDate: "2025-02-01",
        lpaUids: ["M-0000-0000-1005"],
        objectionLpas: [
          {
            uid: "M-0000-0000-1005",
            resolution: "upheld",
            resolutionDate: "2025-02-02",
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/cases/5", "GET", {
      status: 200,
      body: {
        id: 5,
        uId: "M-0000-0000-1005",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 5,
        },
        status: "Draft",
      },
    });

    const mocks = Promise.allSettled([
      cases.warnings.empty("5"),
      cases.tasks.empty("5"),
      digitalLpas.progressIndicators.feesInProgress("M-0000-0000-1005"),
    ]);

    cy.wrap(mocks);

    cy.visit("/lpa/M-0000-0000-1005/objection/1005/outcome");
  });

  it("can be visited from the case summary dashboard", () => {
    cy.visit("/lpa/M-0000-0000-1005/lpa-details").then(() => {
      cy.contains("Objections (1)").click();
      cy.contains("Objection upheld").click();

      cy.url().should(
        "include",
        "/lpa/M-0000-0000-1005/objection/1005/outcome",
      );
      cy.contains("Objection - upheld");
    });
  });

  it("can go Back to LPA page", () => {
    cy.contains("Back").click();
    cy.url().should("contain", "/lpa/M-0000-0000-1005");
  });

  it("can be cancelled, returning to the LPA page", () => {
    cy.contains("Exit").click();
    cy.url().should("contain", "/lpa/M-0000-0000-1005");
  });

  it("Displays upheld objection", () => {
    cy.contains("Objection - upheld");
    cy.contains("Will this stop the LPA from being registered?");
    cy.contains("Will an attorney need to be removed?");
    cy.contains("No notes provided");
  });

  it("Dispalys not upheld objection", () => {
    cy.addMock("/lpa-api/v1/objections/1005", "GET", {
      status: 200,
      body: {
        id: 1005,
        notes: "test",
        objectionType: "factual",
        receivedDate: "2025-02-01",
        lpaUids: ["M-0000-0000-1005"],
        objectionLpas: [
          {
            uid: "M-0000-0000-1005",
            resolution: "notUpheld",
            resolutionDate: "2025-02-02",
            resolutionNotes: "some notes for testing",
          },
        ],
      },
    });

    cy.visit("/lpa/M-0000-0000-1005/objection/1005/outcome").then(() => {
      cy.contains("Objection - not upheld");
      cy.contains("some notes for testing");
      cy.contains("Will this stop the LPA from being registered?").should(
        "not.exist",
      );
      cy.contains("Will an attorney need to be removed?").should("not.exist");
    });
  });
});
