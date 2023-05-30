describe("Create Digital LPA draft", () => {
  beforeEach(() => {
    cy.visit("/digital-lpa/create");
  });

  it("creates a digital LPA", () => {
    cy.contains("Create a draft LPA");

    cy.contains("Health and welfare").click();
    cy.contains("Property and finance").click();

    cy.get("#f-donorFirstname").type("Coleen");
    cy.get("#f-donorMiddlename").type("Stephanie");
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
    cy.get("#f-donorAddressLine1").should(
      "have.value",
      "Office of the Public Guardian"
    );
    cy.get("#f-donorAddressLine2").should("have.value", "1 Something Street");
    cy.get("#f-donorAddressLine3").should("have.value", "Someborough");
    cy.get("#f-donorTown").should("have.value", "Someton");
    cy.get("#f-donorPostcode").should("have.value", "SW1A 1AA");
    cy.get("#f-donorCountry").should("have.value", "GB");

    // Override address manually
    cy.get("#f-donorAddressLine1").clear().type("Fluke House");
    cy.get("#f-donorAddressLine2").clear().type("Summit");
    cy.get("#f-donorAddressLine3").clear().type("Houston");
    cy.get("#f-donorTown").clear().type("South Bend");
    cy.get("#f-donorPostcode").clear().type("AI1 6VW");

    cy.contains("Another person").click();
    cy.get("#f-correspondentFirstname").type("Leon");
    cy.get("#f-correspondentMiddlename").type("Marius");
    cy.get("#f-correspondentSurname").type("Selden");

    cy.get("#f-correspondentSurname")
      .closest(".govuk-radios__conditional")
      .within(() => {
        cy.contains("Enter address manually").click();
        cy.get("#f-correspondentAddressLine1").type(
          "Nitzsche, Nader And Schuppe"
        );
        cy.get("#f-correspondentAddressLine2").type("6064 Alessandro Plain");
        cy.get("#f-correspondentAddressLine3").type("Pittsfield");
        cy.get("#f-correspondentTown").type("Concord");
        cy.get("#f-correspondentPostcode").type("JN2 7UO");
      });

    cy.get("#f-donorPhone").type("07893932118");
    cy.get("#f-donorEmail").type("c.morneault@somehost.example");

    cy.contains("Save draft LPA").click();
    cy.get(".govuk-panel").contains(
      "Draft Health and Welfare and Property and Finance LPAs for the donor Coleen Stephanie Morneault have been saved"
    );
    cy.get(".govuk-panel").contains(
      "Health and Welfare case reference number is M-GHIJ-7890-KLMN"
    );
    cy.get(".govuk-panel").contains(
      "Property and Finance case reference number is M-ABCD-1234-EF56"
    );
  });
});
