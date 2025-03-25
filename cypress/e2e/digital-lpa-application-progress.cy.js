import * as digitalLpas from "../mocks/digitalLpas";
import * as cases from "../mocks/cases";

describe("View the application progress for a digital LPA", () => {
  beforeEach(() => {
    const mocks = Promise.allSettled([
      digitalLpas.get("M-1111-1111-1111"),
      cases.warnings.empty("1111"),
      cases.tasks.empty("1111"),
      digitalLpas.progressIndicators.feesInProgress("M-1111-1111-1111"),
      digitalLpas.get("M-2222-2222-2222", {
        "opg.poas.sirius": {
          id: 2222,
        },
      }),
      cases.warnings.empty("2222"),
      cases.tasks.empty("2222"),
      digitalLpas.progressIndicators.defaultCannotStart("M-2222-2222-2222", [{
        indicator: "RESTRICTIONS_AND_CONDITIONS",
        status: "COMPLETE",
      }]),
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
});
