import * as cases from "../mocks/cases";

describe("Change attorney details form", () => {
  beforeEach(() => {
    const lpaResponse = {
      status: 200,
      body: {
        uId: "M-1111-1111-1110",
        "opg.poas.sirius": {
          id: 555,
          uId: "M-1111-1111-1110",
          status: "in-progress",
          caseSubtype: "personal-welfare",
          createdDate: "31/10/2023",
          investigationCount: 0,
          complaintCount: 0,
          taskCount: 0,
          warningCount: 0,
          donor: {
            id: 33,
          },
          application: {
            donorFirstNames: "James",
            donorLastName: "Rubin",
            donorDob: "22/02/1990",
          },
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
          attorneys: [
            {
              uid: "active-attorney-1",
              firstNames: "Julie",
              lastName: "Rutherford",
              address: {
                line1: "15 Cameron Approach",
                line2: "Nether Collier",
                town: "Worcestershire",
                postcode: "BL2 6DI",
                country: "GB",
              },
              status: "active",
              appointmentType: "original",
              signedAt: "2022-12-19T09:12:59Z",
              dateOfBirth: "1971-11-27",
              mobile: "0123456789",
              email: "j@example.com",
            },
            {
              uid: "replacement-attorney-2",
              firstNames: "Rico",
              lastName: "Welch",
              status: "inactive",
              appointmentType: "replacement",
              signedAt: "2022-12-19T09:12:59Z",
              dateOfBirth: "1998-05-27",
            },
          ],
          status: "in-progress",
          signedAt: "2024-10-18T11:46:24Z",
          lpaType: "pw",
          channel: "online",
          registrationDate: "2024-11-11",
          peopleToNotify: [],
        },
      },
    };

    cy.addMock("/lpa-api/v1/digital-lpas/M-1111-1111-1110", "GET", lpaResponse);
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-1111-1111-1110?presignImages",
      "GET",
      lpaResponse,
    );

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

    cases.warnings.empty("555");

    cy.addMock(
      "/lpa-api/v1/cases/555/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.visit("/lpa/M-1111-1111-1110/attorney/active-attorney-1/change-details");
  });

  it("can be visited from the LPA details attorney Change link", () => {
    cy.visit("/lpa/M-1111-1111-1110/lpa-details").then(() => {
      cy.get(".govuk-accordion__section-button").contains("Attorneys").click();
      cy.get("#f-change-attorney-details").click();
      cy.contains("Change attorney details");
      cy.url().should(
        "contain",
        "/lpa/M-1111-1111-1110/attorney/active-attorney-1/change-details",
      );
      cy.contains("Attorney's first names");
      cy.contains("Attorney's last name");
      cy.contains("Attorney's date of birth");
      cy.contains("Attorney address");
      cy.contains("Attorney's phone number (optional)");
      cy.contains("Attorney's email address (optional)");
    });
  });

  it("can be visited from the LPA details replacement attorney Change link", () => {
    cy.visit("/lpa/M-1111-1111-1110/lpa-details").then(() => {
      cy.get(".govuk-accordion__section-button")
        .contains("Replacement attorneys")
        .click();
      cy.get("#f-change-replacement-attorney-details").click();
      cy.contains("Change replacement attorney details");
      cy.url().should(
        "contain",
        "/lpa/M-1111-1111-1110/attorney/replacement-attorney-2/change-details",
      );
      cy.contains("Replacement attorney's first names");
      cy.contains("Replacement attorney's last name");
      cy.contains("Replacement attorney's date of birth");
      cy.contains("Replacement attorney address");
      cy.contains("Replacement attorney's phone number (optional)");
      cy.contains("Replacement attorney's email address (optional)");
    });
  });

  it("populates attorney details", () => {
    cy.get("#f-firstNames").should("have.value", "Julie");
    cy.get("#f-lastName").should("have.value", "Rutherford");

    cy.get("#f-dob-day").should("have.value", "27");
    cy.get("#f-dob-month").should("have.value", "11");
    cy.get("#f-dob-year").should("have.value", "1971");

    cy.get("#f-address\\.Line1").should("have.value", "15 Cameron Approach");
    cy.get("#f-address\\.Line2").should("have.value", "Nether Collier");
    cy.get("#f-address\\.Line3").should("have.value", "");
    cy.get("#f-address\\.Town").should("have.value", "Worcestershire");
    cy.get("#f-address\\.Postcode").should("have.value", "BL2 6DI");
    cy.get("#f-address\\.Country").should("have.value", "GB");

    cy.get("#f-phoneNumber").should("have.value", "0123456789");
    cy.get("#f-email").should("have.value", "j@example.com");

    cy.get("#f-signedAt-day").should("have.value", "19");
    cy.get("#f-signedAt-month").should("have.value", "12");
    cy.get("#f-signedAt-year").should("have.value", "2022");
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
      "/lpa-api/v1/digital-lpas/M-1111-1111-1110/attorney/active-attorney-1/change-details",
      "PUT",
      {
        status: 204,
      },
    );

    cy.get("#f-firstNames").clear().type("Jonathan");
    cy.get("#f-lastName").clear().type("Ruby");

    cy.get("#f-dob-day").clear().type("31");
    cy.get("#f-dob-month").clear().type("1");
    cy.get("#f-dob-year").clear().type("2000");

    cy.get("#f-address\\.Line1").clear().type("4");
    cy.get("#f-address\\.Line2").clear().type("Gherkin Building");
    cy.get("#f-address\\.Line3").clear().type("33 London Road");
    cy.get("#f-address\\.Town").clear().type("London");
    cy.get("#f-address\\.Postcode").clear().type("B29 6BL");

    cy.get("#f-phoneNumber").clear().type("07777777777");
    cy.get("#f-email").clear().type("jimR@mail.example");

    cy.get("#f-signedAt-day").clear().type("11");
    cy.get("#f-signedAt-month").clear().type("11");
    cy.get("#f-signedAt-year").clear().type("2023");

    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");

    cy.url().should("contain", "/lpa/M-1111-1111-1110/lpa-details");
  });
});
