import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Update objection form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0008", "GET", {
      status: 200,
      body: {
        uId: "M-0000-0000-0008",
        "opg.poas.sirius": {
          id: 8,
          uId: "M-0000-0000-0008",
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
              uId: "M-0000-0000-0009",
              caseSubtype: "property-and-affairs",
              status: "Draft",
              createdDate: "01/11/2023",
            },
            {
              uId: "M-0000-0000-0007",
              caseSubtype: "personal-welfare",
              status: "Registered",
              createdDate: "02/11/2023",
            },
          ],
        },
      },
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0008/objections", "GET", {
      status: 200,
      body: [
        {
          id: 18,
          notes: "test",
          objectionType: "factual",
          receivedDate: "2025-01-01",
          lpaUids: ["M-0000-0000-0008"],
        },
      ],
    });

    cy.addMock("/lpa-api/v1/objections/18", "GET", {
      status: 200,
      body: {
        id: 18,
        notes: "test",
        objectionType: "factual",
        receivedDate: "2025-01-01",
        lpaUids: ["M-0000-0000-0008"],
      },
    });

    cy.addMock("/lpa-api/v1/cases/8", "GET", {
      status: 200,
      body: {
        id: 8,
        uId: "M-0000-0000-0008",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 8,
        },
        status: "Draft",
      },
    });

    const mocks = Promise.allSettled([
      cases.warnings.empty("8"),
      cases.tasks.empty("8"),
      digitalLpas.progressIndicators.feesInProgress("M-0000-0000-0008"),
    ]);

    cy.wrap(mocks);

    cy.visit("/lpa/M-0000-0000-0008/objection/18");
  });

  it("can be visited from the case summary dashboard", () => {
    cy.visit("/lpa/M-0000-0000-0008/lpa-details").then(() => {
      cy.contains("Objection received").click();

      cy.url().should("include", "/lpa/M-0000-0000-0008/objection/18");
      cy.contains("Update Objection");
    });
  });

  it("can go Back to LPA page", () => {
    cy.contains("Back").click();
    cy.url().should("contain", "/lpa/M-0000-0000-0008");
  });

  it("can be cancelled, returning to the LPA page", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-0000-0000-0008");
  });

  it("Can update objection", () => {
    cy.addMock("/lpa-api/v1/objections/18", "PUT", {
      status: 204,
      body: {
        lpaUids: ["M-0000-0000-0008", "M-0000-0000-0009"],
        receivedDate: "13/12/2024",
        objectionType: "prescribed",
        notes: "Test",
      },
    });

    cy.contains("PA M-0000-0000-0009").click();
    cy.get("#f-receivedDate-day").clear();
    cy.get("#f-receivedDate-day").type("13");
    cy.get("#f-receivedDate-month").clear();
    cy.get("#f-receivedDate-month").type("12");
    cy.get("#f-receivedDate-year").clear();
    cy.get("#f-receivedDate-year").type("2024");
    cy.contains("Prescribed").click();

    cy.get("button[type=submit]").click();

    cy.contains("Confirm screen");

    cy.get("input[name=step]").invoke("val", "confirm");
    cy.get("button[type=submit]").click();

    cy.url().should("contain", "/lpa/M-0000-0000-0008");
  });
});
