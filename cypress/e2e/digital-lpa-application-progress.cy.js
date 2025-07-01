import * as digitalLpas from "../mocks/digitalLpas";
import * as cases from "../mocks/cases";

describe("View the application progress for a digital LPA", () => {
  beforeEach(() => {
    const mocks = Promise.allSettled([
      digitalLpas.get("M-1111-1111-1111"),
      cases.warnings.empty("1111"),
      cases.tasks.empty("1111"),
      digitalLpas.objections.empty("M-1111-1111-1111"),
      digitalLpas.progressIndicators.feesInProgress("M-1111-1111-1111"),

      digitalLpas.get("M-2222-2222-2222", {
        "opg.poas.sirius": {
          id: 2222,
        },
      }),
      cases.warnings.empty("2222"),
      cases.tasks.empty("2222"),
      digitalLpas.objections.empty("M-2222-2222-2222"),
      digitalLpas.progressIndicators.defaultCannotStart("M-2222-2222-2222", [
        {
          indicator: "RESTRICTIONS_AND_CONDITIONS",
          status: "COMPLETE",
        },
      ]),

      digitalLpas.get("M-3333-3333-3333", {
        "opg.poas.sirius": {
            id: 3333,
            application: {
                source: "paper",
                donorIdentityCheck: {
                    state: "COUNTER_SERVICE_STARTED",
                    checkedAt: "2024-07-01T16:06:08Z",
                    reference: "712254d5-4cf4-463c-96c1-67744b70043e",
                },
            },
        },
        "opg.poas.lpastore": {
            channel: "paper",
        },
      }),
      cases.warnings.empty("3333"),
      cases.tasks.empty("3333"),
      digitalLpas.objections.empty("M-3333-3333-3333"),
      digitalLpas.progressIndicators.feesInProgress("M-3333-3333-3333"),

      digitalLpas.get("M-4444-4444-4444", {
        "opg.poas.sirius": {
            id: 4444,
            application: {
                source: "paper",
            },
        },
        "opg.poas.lpastore": {
          channel: "paper",
          donor: {
            lastName: "Rubix",
            uid: "5ff557dd-1e27-4426-9681-ed6e90c2c08d",
            address: {
              postcode: "W8A 0IK",
              country: "GB",
              town: "Edinburgh",
              line1: "1 Scotland Street",
            },
            dateOfBirth: "1938-03-18",
            firstNames: "Jack",
            contactLanguagePreference: "en",
            identityCheck: {
              type: "opg-paper-id",
              checkedAt: "2025-06-29T15:06:29Z",
            },
            email: "jrubix@mail.example",
          },
        },
      }),
      cases.warnings.empty("4444"),
      cases.tasks.empty("4444"),
      digitalLpas.objections.empty("M-4444-4444-4444"),
      digitalLpas.progressIndicators.defaultCannotStart("M-4444-4444-4444", [
        {
            indicator: "DONOR_ID",
            status: "COMPLETE",
        },
      ]),
    ]);

    cy.wrap(mocks);

    cy.visit("/lpa/M-1111-1111-1111");
  });

  it("shows not started Restrictions and Conditions progress indicator", () => {
    cy.contains("Restrictions and conditions (Not started)");
  });

  it("shows complete Restrictions and Conditions progress indicator", () => {
    cy.visit("/lpa/M-2222-2222-2222");

    cy.contains("Restrictions and conditions (Complete)");
  });

  it("shows unable to attempt online ID check PI content", () => {
    cy.visit("/lpa/M-3333-3333-3333");

    cy.contains("Donor identity confirmation").click();

    cy.contains(
        "Donor unable to attempt online ID check on 1 July 2024 - Post Office to check identity",
    );
  });

  it("shows unable to attempt paper ID check PI content", () => {
    const mocks = Promise.allSettled([
        digitalLpas.get("M-3333-3333-3333", {
            "opg.poas.sirius": {
                id: 3333,
                application: {
                    source: "PHONE",
                    donorIdentityCheck: {
                        state: "COUNTER_SERVICE_STARTED",
                        checkedAt: "2024-07-01T16:06:08Z",
                        reference: "712254d5-4cf4-463c-96c1-67744b70043e",
                    },
                },
            },
        }),
    ]);

    cy.wrap(mocks);

    cy.visit("/lpa/M-3333-3333-3333");

    cy.contains("Donor identity confirmation").click();

    cy.contains(
        "Donor unable to attempt phone ID check on 1 July 2024 - Post Office to check identity",
    );
  });

  it("shows complete Donor identity confirmation progress indicator", () => {
    cy.visit("/lpa/M-4444-4444-4444");

    cy.contains("Donor identity confirmation").click();

    cy.contains("Passed phone identity check on 29 June 2025");
  });
});
