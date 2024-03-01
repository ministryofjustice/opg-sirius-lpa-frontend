describe("Create Digital LPA draft", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "private-mlpa"],
      },
    });

    cy.visit("/digital-lpa/create");
  });

  it("creates a digital LPA", () => {
    cy.contains("Create a draft LPA");

    cy.contains("Personal welfare").click();
    cy.contains("Property and affairs").click();

    cy.get("#f-donorFirstname").type("Coleen Stephanie");
    cy.get("#f-donorSurname").type("Morneault");

    cy.get("#f-dob-day").type("8");
    cy.get("#f-dob-month").type("4");
    cy.get("#f-dob-year").type("1952");

    // Check postcode finder auto-populates
    cy.contains("Enter a UK postcode, or enter the address manually")
      .nextUntil(".govuk-input")
      .next()
      .type("SW1A 1AA");
    cy.contains("Find address").click();
    cy.contains("Enter address manually").click();
    cy.get("#f-donorAddress\\.Line1").should(
      "have.value",
      "Office of the Public Guardian",
    );
    cy.get("#f-donorAddress\\.Line2").should(
      "have.value",
      "1 Something Street",
    );
    cy.get("#f-donorAddress\\.Line3").should("have.value", "Someborough");
    cy.get("#f-donorAddress\\.Town").should("have.value", "Someton");
    cy.get("#f-donorAddress\\.Postcode").should("have.value", "SW1A 1AA");
    cy.get("#f-donorAddress\\.Country").should("have.value", "GB");

    // Override address manually
    cy.get("#f-donorAddress\\.Line1").clear().type("Fluke House");
    cy.get("#f-donorAddress\\.Line2").clear().type("Summit");
    cy.get("#f-donorAddress\\.Line3").clear().type("Houston");
    cy.get("#f-donorAddress\\.Town").clear().type("South Bend");
    cy.get("#f-donorAddress\\.Postcode").clear().type("AI1 6VW");

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
    cy.get("#f-donorEmail").type("c.morneault@somehost.example");

    cy.contains("Confirm and create draft LPA").click();
    cy.get(".govuk-panel").contains(
      "Draft Personal welfare and Property and affairs LPAs for the donor Coleen Stephanie Morneault have been saved",
    );
    cy.get(".govuk-panel").contains(
      "Personal welfare case reference number is M-GHIJ-7890-KLMN",
    );
    cy.get(".govuk-panel").contains(
      "Property and affairs case reference number is M-ABCD-1234-EF56",
    );
  });
});
