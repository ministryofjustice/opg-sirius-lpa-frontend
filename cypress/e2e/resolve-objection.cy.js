import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Update objection form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0005", "GET", {
      status: 200,
      body: {
        uId: "M-0000-0000-0005",
        "opg.poas.sirius": {
          id: 5,
          uId: "M-0000-0000-0005",
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
          linkedDigitalLpas: [
            {
              uId: "M-0000-0000-0004",
              caseSubtype: "property-and-affairs",
              status: "Draft",
              createdDate: "01/11/2023",
            },
            {
              uId: "M-0000-0000-0003",
              caseSubtype: "personal-welfare",
              status: "Registered",
              createdDate: "02/11/2023",
            },
          ],
        },
      },
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0005/objections", "GET", {
      status: 200,
      body: [
        {
          id: 15,
          notes: "test",
          objectionType: "factual",
          receivedDate: "2025-01-01",
          lpaUids: ["M-0000-0000-0005"],
        },
      ],
    });

    cy.addMock("/lpa-api/v1/objections/15", "GET", {
      status: 200,
      body: {
        id: 15,
        notes: "test",
        objectionType: "factual",
        receivedDate: "2025-02-01",
        lpaUids: ["M-0000-0000-0005", "M-0000-0000-0004"],
      },
    });

    cy.addMock("/lpa-api/v1/cases/5", "GET", {
      status: 200,
      body: {
        id: 5,
        uId: "M-0000-0000-0005",
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
      digitalLpas.progressIndicators.feesInProgress("M-0000-0000-0005"),
    ]);

    cy.wrap(mocks);

    cy.visit("/lpa/M-0000-0000-0005/objection/15/resolve");
  });

  it("can be visited from the case summary dashboard", () => {
    cy.visit("/lpa/M-0000-0000-0005/lpa-details").then(() => {
      cy.contains("Objections (1)").click();
      cy.contains("Record objection outcome").click();

      cy.url().should("include", "/lpa/M-0000-0000-0005/objection/15");
      cy.contains("Record objection outcome");
    });
  });

  it("can go Back to LPA page", () => {
    cy.contains("Back").click();
    cy.url().should("contain", "/lpa/M-0000-0000-0005");
  });

  it("can be cancelled, returning to the LPA page", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-0000-0000-0005");
  });

  it("Can resolve objection", () => {
    cy.addMock("/lpa-api/v1/objections/15/resolution/M-0000-0000-0005", "PUT", {
      status: 204,
      body: {
        resolution: "upheld",
        resolutionNotes: "Test",
      },
    });

    cy.addMock("/lpa-api/v1/objections/15/resolution/M-0000-0000-0004", "PUT", {
      status: 204,
      body: {
        resolution: "notUpheld",
        resolutionNotes: "Test",
      },
    });

    cy.contains("What is the outcome for M-0000-0000-0005");
    cy.get("#f-resolution-upheld-M-0000-0000-0005").click();
    cy.get("#f-resolutionNotes-M-0000-0000-0005").type("Test");
    cy.get("#f-resolution-notUpheld-M-0000-0000-0004").click();
    cy.get("#f-resolutionNotes-M-0000-0000-0004").type("Test");

    cy.get("button[type=submit]").click();

    cy.url().should("contain", "/lpa/M-0000-0000-0005");
  });

  it("Displays upheld resolved objection", () => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0005/objections", "GET", {
      status: 200,
      body: [
        {
          id: 15,
          notes: "test",
          objectionType: "factual",
          receivedDate: "2025-01-01",
          lpaUids: ["M-0000-0000-0005"],
          objectionLpas: [
            {
              uid: "M-0000-0000-0005",
              resolution: "upheld",
              resolutionDate: "2025-02-02",
            },
          ],
        },
      ],
    });

    cy.visit("/lpa/M-0000-0000-0005").then(() => {
      cy.contains("Objection upheld");
      cy.contains("Outcome recorded on 2 February 2025");
    });
  });

  it("Displays not upheld resolved objection", () => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0005/objections", "GET", {
      status: 200,
      body: [
        {
          id: 15,
          notes: "test",
          objectionType: "factual",
          receivedDate: "2025-01-01",
          lpaUids: ["M-0000-0000-0005"],
          objectionLpas: [
            {
              uid: "M-0000-0000-0005",
              resolution: "notUpheld",
              resolutionDate: "2025-02-02",
            },
          ],
        },
      ],
    });

    cy.visit("/lpa/M-0000-0000-0005").then(() => {
      cy.contains("Objection not upheld");
      cy.contains("Outcome recorded on 2 February 2025");
    });
  });
});
