import * as cases from "../mocks/cases";

describe("Change draft form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-1111-2222-1111", "GET", {
      status: 200,
      body: {
        uId: "M-1111-2222-1111",
        "opg.poas.sirius": {
          id: 565,
          uId: "M-1111-2222-1111",
          status: "Draft",
          caseSubtype: "property-and-affairs",
          createdDate: "31/10/2023",
          investigationCount: 2,
          complaintCount: 1,
          taskCount: 2,
          warningCount: 4,
          donor: {
            id: 44,
            firstname: "Agnes",
            surname: "Hartley",
            dob: "17/06/1982",
            addressLine1: "Apartment 3",
            addressLine2: "Gherkin Building",
            town: "London",
            postcode: "B15 3AA",
            country: "GB",
            phone: "073656249524",
            email: "agnes@host.example",
          },
          application: {
            donorFirstNames: "Agnes",
            donorLastName: "Hartley",
            donorDob: "17/06/1982",
            donorEmail: "agnes@host.example",
            donorPhone: "073656249524",
            donorAddress: {
              addressLine1: "Apartment 3",
              addressLine2: "Gherkin Building",
              country: "GB",
              postcode: "B15 3AA",
              town: "London",
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
          linkedDigitalLpas: [],
        },
        "opg.poas.lpastore": null,
      },
    });

    cy.addMock("/lpa-api/v1/cases/565", "GET", {
      status: 200,
      body: {
        id: 565,
        uId: "M-1111-2222-1111",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 44,
        },
        status: "Draft",
      },
    });

    cases.warnings.empty("565");

    cy.addMock(
      "/lpa-api/v1/cases/565/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.visit("/lpa/M-1111-2222-1111/change-draft");
  });

  it("can be visited from the LPA draft Change link", () => {
    cy.visit("/lpa/M-1111-2222-1111/lpa-details").then(() => {
      cy.get("#f-change-draft").click();
      cy.contains("Change donor details");
      cy.url().should("contain", "/lpa/M-1111-2222-1111/change-draft");
    });
  });

  it("populates draft details", () => {
    cy.get("#f-firstNames").should("have.value", "Agnes");
    cy.get("#f-lastName").should("have.value", "Hartley");

    cy.get("#f-dob-day").should("have.value", "17");
    cy.get("#f-dob-month").should("have.value", "6");
    cy.get("#f-dob-year").should("have.value", "1982");

    cy.get("#f-address\\.Line1").should("have.value", "Apartment 3");
    cy.get("#f-address\\.Line2").should("have.value", "Gherkin Building");
    cy.get("#f-address\\.Line3").should("have.value", "");
    cy.get("#f-address\\.Town").should("have.value", "London");
    cy.get("#f-address\\.Postcode").should("have.value", "B15 3AA");
    cy.get("#f-address\\.Country").should("have.value", "GB");

    cy.get("#f-phoneNumber").should("have.value", "073656249524");
    cy.get("#f-email").should("have.value", "agnes@host.example");
  });

  it("can go Back to LPA details", () => {
    cy.contains("Back to LPA details").click();
    cy.url().should("contain", "/lpa/M-1111-2222-1111/lpa-details");
  });

  it("can be cancelled, returning to the LPA details", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-1111-2222-1111/lpa-details");
  });

  it("can edit all draft details and redirect to lpa details", () => {
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-1111-2222-1111/change-draft",
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
    cy.get("#f-address\\.Line2").clear().type("Halls");
    cy.get("#f-address\\.Line3").clear().type("19 Grange Road");
    cy.get("#f-address\\.Town").clear().type("Birmingham");
    cy.get("#f-address\\.Postcode").clear().type("B29 6BL");

    cy.get("#f-phoneNumber").clear().type("07777777777");
    cy.get("#f-email").clear().type("jimR@mail.example");

    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");

    cy.url().should("contain", "/lpa/M-1111-2222-1111/lpa-details");
  });
});
