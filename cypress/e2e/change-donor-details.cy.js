describe("Change donor details form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0001", "GET", {
      status: 200,
      body: {
        uId: "M-0000-0000-0001",
        "opg.poas.sirius": {
          id: 666,
          uId: "M-0000-0000-0001",
          status: "Draft",
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
            donorEmail: "jrubin@mail.example",
            donorPhone: "073656249524",
            donorAddress: {
              addressLine1: "Apartment 3",
              country: "GB",
              postcode: "B15 3AA",
              town: "Birmingham",
            },
            correspondentFirstNames: "Kendrick",
            correspondentLastName: "Lamar",
            correspondentAddress: {
              addressLine1: "Flat 3",
              country: "GB",
              postcode: "SW1 1AA",
              town: "London",
            },
          },
        },
        "opg.poas.lpastore": {
          donor: {
            uid: "5ff557dd-1e27-4426-9681-ed6e90c2c08d",
            firstNames: "James",
            lastName: "Rubin",
            otherNamesKnownBy: "Somebody",
            dateOfBirth: "1990-02-22",
            address: {
              line1: "Apartment 3",
              town: "Birmingham",
              country: "GB",
              postcode: "B15 3AA",
            },
            contactLanguagePreference: "en",
            email: "jrubin@mail.example",
          },
          attorneys: [
            {
              firstNames: "Esther",
              lastName: "Greenwood",
              status: "active",
            },
          ],
          certificateProvider: {
            uid: "e4d5e24e-2a8d-434e-b815-9898620acc71",
            firstNames: "Timothy",
            lastNames: "Turner",
            signedAt: "2022-12-18T11:46:24Z",
          },
          signedAt: "2024-10-18T11:46:24Z",
          lpaType: "pw",
          channel: "online",
          registrationDate: "2024-11-11",
          peopleToNotify: [],
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/666", "GET", {
      status: 200,
      body: {
        id: 666,
        uId: "M-0000-0000-0001",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 666,
        },
        status: "Draft",
      },
    });

    cy.addMock("/lpa-api/v1/cases/666/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      "/lpa-api/v1/cases/666/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.visit("/change-donor-details?uid=M-0000-0000-0001");
  });

  it("can be visited from the LPA details Change link", () => {
    cy.visit("/lpa/M-0000-0000-0001/lpa-details").then(() => {
      Cypress.$("span:contains('Donor')").closest("button")[0].click();
      cy.contains("Change").click();
      cy.contains("Details that apply to all LPAs for this donor");
    });
  });

  it("populates donor details", () => {
    cy.get("#f-firstNames").should("have.value", "James");
    cy.get("#f-lastName").should("have.value", "Rubin");
    cy.get("#f-otherNamesKnownBy").should("have.value", "Somebody");

    cy.get("#f-dob-day").should("have.value", "22");
    cy.get("#f-dob-month").should("have.value", "2");
    cy.get("#f-dob-year").should("have.value", "1990");

    cy.get("#f-address\\.Line1").should("have.value", "Apartment 3");
    cy.get("#f-address\\.Line2").should("have.value", "");
    cy.get("#f-address\\.Line3").should("have.value", "");
    cy.get("#f-address\\.Town").should("have.value", "Birmingham");
    cy.get("#f-address\\.Postcode").should("have.value", "B15 3AA");
    cy.get("#f-address\\.Country").should("have.value", "GB");

    cy.get("#f-phoneNumber").should("have.value", "073656249524");
    cy.get("#f-email").should("have.value", "jrubin@mail.example");

    cy.get("#f-lpaSignedOn-day").should("have.value", "18");
    cy.get("#f-lpaSignedOn-month").should("have.value", "10");
    cy.get("#f-lpaSignedOn-year").should("have.value", "2024");
  });

  it("can go Back to LPA details", () => {
    cy.contains("Back to LPA details").click();
    cy.url().should("contain", "/lpa/M-0000-0000-0001/lpa-details");
  });

  it("can be cancelled, returning to the LPA details", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-0000-0000-0001/lpa-details");
  });

  it("Can edit all donor details", () => {
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-0000-0000-0001/change-donor-details",
      "PUT",
      {
        status: 204,
      },
    );

    cy.get("#f-firstNames").clear();
    cy.get("#f-firstNames").type("Jonathan");

    cy.get("#f-lastName").clear();
    cy.get("#f-lastName").type("Ruby");

    cy.get("#f-otherNamesKnownBy").clear();
    cy.get("#f-otherNamesKnownBy").type("Jim");

    cy.get("#f-dob-day").clear();
    cy.get("#f-dob-month").clear();
    cy.get("#f-dob-year").clear();

    cy.get("#f-dob-day").type("31");
    cy.get("#f-dob-month").type("1");
    cy.get("#f-dob-year").type("2000");

    cy.get("#f-address\\.Line2").clear();
    cy.get("#f-address\\.Line3").clear();
    cy.get("#f-address\\.Town").clear();
    cy.get("#f-address\\.Postcode").clear();

    cy.get("#f-address\\.Line1").type("4");
    cy.get("#f-address\\.Line2").type("Gherkin Building");
    cy.get("#f-address\\.Line3").type("33 London Road");
    cy.get("#f-address\\.Town").type("London");
    cy.get("#f-address\\.Postcode").type("B29 6BL");

    cy.get("#f-phoneNumber").clear();
    cy.get("#f-phoneNumber").type("07777777777");

    cy.get("#f-email").clear();
    cy.get("#f-email").type("jimR@mail.example");

    cy.get("#f-lpaSignedOn-day").clear();
    cy.get("#f-lpaSignedOn-month").clear();
    cy.get("#f-lpaSignedOn-year").clear();

    cy.get("#f-lpaSignedOn-day").type("11");
    cy.get("#f-lpaSignedOn-month").type("11");
    cy.get("#f-lpaSignedOn-year").type("2023");

    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");

    cy.url().should("contain", "/lpa/M-0000-0000-0001/lpa-details");

    cy.contains("Donor").click();
    //todo: Add check for new details after lpa-store changes added
  });
});
