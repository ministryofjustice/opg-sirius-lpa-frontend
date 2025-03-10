import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Change certificate provider details form", () => {
  beforeEach(() => {
    digitalLpas.get("M-1111-1111-1111", {
      "opg.poas.lpastore": {
        certificateProvider: {
          uid: "c362e307-71b9-4070-bdde-c19b4cdf5c1a",
          channel: "online",
          firstNames: "Rhea",
          lastName: "Vandervort",
          address: {
            line1: "290 Vivien Road",
            line2: "Lower Court",
            line3: "Tillman",
            town: "Oxfordshire",
            postcode: "JJ80 7QL",
            country: "GB",
          },
          email: "Rhea.Vandervort@example.com",
          phone: "0151 087 7256",
          signedAt: "2025-01-19T09:12:59Z",
        },
      }
    });

    cases.warnings.empty("1111");
    cases.tasks.empty("1111");

    digitalLpas.progressIndicators.feesInProgress("M-1111-1111-1111");

    cy.visit("/lpa/M-1111-1111-1111/certificate-provider/change-details");
  });

  it("can be visited from the LPA details certificate provider Change link", () => {
    cy.visit("/lpa/M-1111-1111-1111/lpa-details").then(() => {
      cy.get(".govuk-accordion__section-button")
        .contains("Certificate provider")
        .click();
      cy.get("#f-change-certificate-provider-details").click();
      cy.contains("Change certificate provider details");
      cy.url().should(
        "contain",
        "/lpa/M-1111-1111-1111/certificate-provider/change-details",
      );
    });
  });

  it("can submit the change details form", () => {
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-1111-1111-1111/change-certificate-provider-details",
      "PUT",
      {
        status: 204,
      },
    );

    cy.get("#f-firstNames").should("have.value", "Rhea");
    cy.get("#f-lastName").should("have.value", "Vandervort");

    cy.get("#f-address\\.Line1").should("have.value", "290 Vivien Road");
    cy.get("#f-address\\.Line2").should("have.value", "Lower Court");
    cy.get("#f-address\\.Line3").should("have.value", "Tillman");
    cy.get("#f-address\\.Town").should("have.value", "Oxfordshire");
    cy.get("#f-address\\.Postcode").should("have.value", "JJ80 7QL");
    cy.get("#f-address\\.Country").should("have.value", "GB");

    cy.get("#f-phone").should("have.value", "0151 087 7256");
    cy.get("#f-email").should("have.value", "Rhea.Vandervort@example.com");

    cy.get("#f-signedAt-day").should("have.value", "19");
    cy.get("#f-signedAt-month").should("have.value", "1");
    cy.get("#f-signedAt-year").should("have.value", "2025");

    cy.get("#f-firstNames").clear().type("Wilfredo");
    cy.get("#f-lastName").clear().type("Morissette");

    cy.get("#f-address\\.Line1").clear().type("8 Christine Ridge");
    cy.get("#f-address\\.Line2").clear().type("Schiller Gardens");
    cy.get("#f-address\\.Line3").clear().type("Stoltenberg");
    cy.get("#f-address\\.Town").clear().type("Dyfed");
    cy.get("#f-address\\.Postcode").clear().type("YH7 4SO");

    cy.get("#f-phone").clear().type("0953 339 6087");
    cy.get("#f-email").clear().type("Wilfredo.Morissette@example.com");

    cy.get("#f-signedAt-day").clear().type("25");
    cy.get("#f-signedAt-month").clear().type("6");
    cy.get("#f-signedAt-year").clear().type("2025");

    cy.contains("Save and continue").click();
    cy.url().should("contain", "/lpa/M-1111-1111-1111/lpa-details");
  });

  it("can go Back to LPA details", () => {
    cy.contains("Back to LPA details").click();
    cy.url().should("contain", "/lpa/M-1111-1111-1111/lpa-details");
  });

  it("can be cancelled, returning to the LPA details", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-1111-1111-1111/lpa-details");
  });
});
