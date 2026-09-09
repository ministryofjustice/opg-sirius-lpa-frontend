describe("Create Digital LPA draft", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "private-mlpa"],
      },
    });
    cy.addMock(`/lpa-api/v1/digital-lpas/M-GHIJ-7890-KLMN`, "GET", {
      status: 200,
      body: {
        "opg.poas.sirius": {
          donor: {
            id: 33,
            firstname: "Coleen Stephanie",
            surname: "Morneault",
          },
        },
      },
    });

    cy.addMock(`/lpa-api/v1/digital-lpas/M-ABCD-1234-EF56`, "GET", {
      status: 200,
      body: {
        "opg.poas.sirius": {
          donor: {
            id: 33,
            firstname: "Coleen Stephanie",
            surname: "Morneault",
          },
        },
      },
    });

    cy.visit("/digital-lpa/create");
  });

  it("creates a digital LPA", () => {
    cy.addMock("/lpa-api/v1/donors/130/digital-lpas", "POST", {
      status: 201,
      body: [
        {
          caseSubtype: "personal-welfare",
          uId: "M-GHIJ-7890-KLMN",
        },
        {
          caseSubtype: "property-and-affairs",
          uId: "M-ABCD-1234-EF56",
        },
      ],
    });

    cy.contains("Create a draft LPA");

    cy.contains("Personal welfare").click();
    cy.contains("Property and affairs").click();

    cy.get("#f-donorFirstname").type("Coleen Stephanie");
    cy.get("#f-donorSurname").type("Morneault");

    cy.get("#f-dob-day").type("8");
    cy.get("#f-dob-month").type("4");
    cy.get("#f-dob-year").type("1952");
    cy.get(
      '[data-app-address-finder-label="Donor’s address"] > :nth-child(1) > .govuk-details > .govuk-details__summary',
    ).click();

    // Override address manually
    cy.get("#f-donorAddress\\.Line1").type("Fluke House");
    cy.get("#f-donorAddress\\.Line2").type("Summit");
    cy.get("#f-donorAddress\\.Line3").type("Houston");
    cy.get("#f-donorAddress\\.Town").type("South Bend");
    cy.get("#f-donorAddress\\.Postcode").type("AI1 6VW");

    cy.contains("Another person").click();
    cy.get("#f-correspondentFirstname").type("Leon");
    cy.get("#f-correspondentSurname").type("Selden");

    cy.get("#f-correspondentSurname")
      .closest(".govuk-radios__conditional")
      .within(() => {
        cy.contains("Enter address manually").click();
        cy.get("#f-correspondentAddress\\.Line1").type(
          "Nitzsche, Nader And Schuppe",
        );
        cy.get("#f-correspondentAddress\\.Line2").type("6064 Alessandro Plain");
        cy.get("#f-correspondentAddress\\.Line3").type("Pittsfield");
        cy.get("#f-correspondentAddress\\.Town").type("Concord");
        cy.get("#f-correspondentAddress\\.Postcode").type("JN2 7UO");
      });

    cy.get("#f-donorPhone").type("07893932118");
    cy.get("#f-donorEmail").type("c.morneault@example.com");

    cy.contains("Confirm and create draft LPA").click();
    cy.get(".govuk-notification-banner")
      .should("be.visible")
      .within(() => {
        cy.contains(
          "2 draft LPAs for Coleen Stephanie Morneault have been created.",
        );
        cy.contains("M-GHIJ-7890-KLMN Personal welfare");
        cy.contains("M-ABCD-1234-EF56 Property and affairs");
      });
  });
});
