import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Change case status", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-LPA3-3333",
        "opg.poas.sirius": {
          id: 333,
          uId: "M-DIGI-LPA3-3333",
          status: "Draft",
          caseSubtype: "property-and-affairs",
          createdDate: "31/10/2023",
          dueDate: "01/12/2023",
          donor: {
            id: 33,
          },
          application: {
            donorFirstNames: "Agnes",
            donorLastName: "Hartley",
            donorDob: "27/05/1998",
            donorEmail: "agnes@host.example",
            donorPhone: "073656249524",
            donorAddress: {
              addressLine1: "Apartment 3",
              addressLine2: "Gherkin Building",
              addressLine3: "33 London Road",
              country: "GB",
              postcode: "B15 3AA",
              town: "Birmingham",
            },
            correspondentFirstNames: "Kendrick",
            correspondentLastName: "Lamar",
            correspondentAddress: {
              addressLine1: "Flat 3",
              addressLine2: "Digital LPA Lane",
              addressLine3: "Somewhere",
              country: "GB",
              postcode: "SW1 1AA",
              town: "London",
            },
          },
        },
        "opg.poas.lpastore": {
          attorneys: [
            {
              firstNames: "Esther",
              lastName: "Greenwood",
              status: "active",
              appointmentType: "original",
            },
            {
              firstNames: "Volo",
              lastName: "McSpolo",
              status: "active",
              appointmentType: "original",
            },
            {
              firstNames: "Susanna",
              lastName: "Kaysen",
              status: "removed",
              appointmentType: "original",
            },
            {
              firstNames: "Philomena",
              lastName: "Guinea",
              status: "inactive",
              appointmentType: "replacement",
            },
          ],
          lpaType: "pf",
          channel: "online",
          status: "draft",
          peopleToNotify: [],
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/333", "GET", {
      status: 200,
      body: {
        id: 333,
        uId: "M-DIGI-LPA3-3333",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 33,
        },
        status: "Draft",
      },
    });

    cy.addMock("/lpa-api/v1/reference-data/caseChangeReason", "GET", {
      status: 200,
      body: [{
        handle: "LPA_DOES_NOT_WORK",
        label: "The LPA does not work and cannot be changed",
        parentSources: ["cannot-register"],
      }],
    });

    const mocks = Promise.allSettled([
      cases.warnings.empty("333"),
      cases.tasks.empty("333"),
      digitalLpas.objections.empty("M-DIGI-LPA3-3333"),
    ]);

    cy.wrap(mocks);

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333/update-case-status",
      "PUT",
      {
        status: 204,
        body: [],
      },
    );
  });

  it("changes the digital lpa case status", () => {
    cy.visit("/change-case-status?uid=M-DIGI-LPA3-3333");
    cy.contains("Change case status");
    cy.contains("M-DIGI-LPA3-3333");
    cy.get(".moj-alert").should("not.exist");
    cy.contains(".govuk-radios__label", "Draft")
      .parent()
      .get("input")
      .should("be.checked");
    cy.contains(".govuk-radios__label", "Cannot register").click();
    cy.contains(
      ".govuk-radios__label",
      "The LPA does not work and cannot be changed",
    ).click();
    cy.get("button[type=submit]").click();
    cy.url().should("contain", "/lpa/M-DIGI-LPA3-3333");
  });
});
