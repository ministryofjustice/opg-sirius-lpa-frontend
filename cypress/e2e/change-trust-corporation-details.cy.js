import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Change trust corporation details form", () => {
  beforeEach(() => {
    const mocks = Promise.allSettled([
      digitalLpas.get("M-1111-1111-1110", {
        "opg.poas.sirius": {
          id: 555,
          uId: "M-1111-1111-1110",
          status: "in-progress",
        },
        "opg.poas.lpastore": {
          donor: {
            uid: "5ff557dd-1e27-4426-9681-ed6e90c2c08d",
            firstNames: "James",
            lastName: "Rubin",
            otherNamesKnownBy: "Somebody",
            dateOfBirth: "1990-02-22",
            contactLanguagePreference: "en",
            email: "jrubin@mail.example",
          },
          trustCorporations: [
            {
              uid: "active-trust-corp-1",
              Name: "Trust Me Ltd.",
              CompanyNumber: "123456789",
              status: "active",
              appointmentType: "original",
              mobile: "077577575757",
              email: "trust.me.once@does.not.exist",
              address: {
                line1: "9 Mount",
                line2: "Pleasant Drive",
                town: "East Harling",
                postcode: "NR16 2GB",
                country: "GB",
              },
            },
            {
              uid: "replacement-trust-corp-2",
              Name: "Trust Me Again Ltd.",
              CompanyNumber: "987654321",
              status: "inactive",
              appointmentType: "replacement",
            },
          ],
          status: "in-progress",
          signedAt: "2024-10-18T11:46:24Z",
          lpaType: "pw",
          channel: "online",
          registrationDate: "2024-11-11",
          peopleToNotify: [],
        },
      }),
      cases.warnings.empty("555"),
      cases.tasks.empty("555"),
      digitalLpas.objections.empty("M-1111-1111-1110"),
    ]);

    cy.wrap(mocks);

    cy.addMock("/lpa-api/v1/cases/555", "GET", {
      status: 200,
      body: {
        id: 555,
        uId: "M-1111-1111-1110",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 33,
        },
      },
    });

    cy.visit(
      "/lpa/M-1111-1111-1110/trust-corporation/active-trust-corp-1/change-details",
    );
  });

  it("can be visited from the LPA details attorney update link", () => {
    cy.visit("/lpa/M-1111-1111-1110/lpa-details");
    cy.get(".govuk-accordion__section-button").contains("Attorneys").click();
    cy.get("#f-change-trust-corporation-details").click();
    cy.contains("Change attorney details");
    cy.url().should(
      "contain",
      "/lpa/M-1111-1111-1110/trust-corporation/active-trust-corp-1/change-details",
    );
    cy.contains("Trust corporation name");
    cy.contains("Company address");
    cy.contains("Company email address (optional)");
    cy.contains("Company phone number (optional)");
    cy.contains("Company registration number");
  });

  it("can be visited from the LPA details replacement attorney update link", () => {
    cy.visit("/lpa/M-1111-1111-1110/lpa-details");
    cy.get(".govuk-accordion__section-button")
      .contains("Replacement attorneys")
      .click();
    cy.get("#f-change-replacement-trust-corporation-details").click();
    cy.contains("Change replacement attorney details");
    cy.contains("Company registration number");
    cy.url().should(
      "contain",
      "/lpa/M-1111-1111-1110/trust-corporation/replacement-trust-corp-2/change-details",
    );
    cy.contains("Trust corporation name");
    cy.contains("Company address");
    cy.contains("Company email address (optional)");
    cy.contains("Company phone number (optional)");
    cy.contains("Company registration number");
  });

  it("populates trust corporation details", () => {
    cy.get("#f-name").should("have.value", "Trust Me Ltd.");

    cy.get(String.raw`#f-address\.Line1`).should("have.value", "9 Mount");
    cy.get(String.raw`#f-address\.Line2`).should(
      "have.value",
      "Pleasant Drive",
    );
    cy.get(String.raw`#f-address\.Line3`).should("have.value", "");
    cy.get(String.raw`#f-address\.Town`).should("have.value", "East Harling");
    cy.get(String.raw`#f-address\.Postcode`).should("have.value", "NR16 2GB");
    cy.get(String.raw`#f-address\.Country`).should("have.value", "GB");

    cy.get("#f-phoneNumber").should("have.value", "077577575757");
    cy.get("#f-email").should("have.value", "trust.me.once@does.not.exist");
    cy.get("#f-companyNumber").should("have.value", "123456789");
  });

  it("can go Back to LPA details", () => {
    cy.contains("Back to LPA details").click();
    cy.url().should("contain", "/lpa/M-1111-1111-1110/lpa-details");
  });

  it("can be cancelled, returning to the LPA details", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-1111-1111-1110/lpa-details");
  });

  it("can edit all attorney details and redirect to lpa details", () => {
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-1111-1111-1110/trust-corporation/active-trust-corp-1/change-details",
      "PUT",
      {
        status: 204,
      },
    );

    cy.get("#f-name").clear().type("Trust Ltd");

    cy.get("#f-address\\.Line1").clear().type("12");
    cy.get("#f-address\\.Line2").clear().type("Building");
    cy.get("#f-address\\.Line3").clear().type("Road");
    cy.get("#f-address\\.Town").clear().type("London");
    cy.get("#f-address\\.Postcode").clear().type("E14 2SH");

    cy.get("#f-phoneNumber").clear().type("07777777777");
    cy.get("#f-email").clear().type("test@test.com");
    cy.get("#f-companyNumber").clear().type("112233");

    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");

    cy.url().should("contain", "/lpa/M-1111-1111-1110/lpa-details");
  });
});
