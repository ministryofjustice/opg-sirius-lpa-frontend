describe("Select or create correspondent", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/2", "GET", {
      status: 200,
      body: {
        attorneys: [
          {
            id: 3,
            salutation: "Prof",
            firstname: "Melanie",
            middlenames: "Josefina",
            surname: "Vanvolkenburg",
            dob: "19/04/1978",
            addressLine1: "29737 Andrew Plaza",
            addressLine2: "Apt. 814",
            addressLine3: "Gislasonside",
            town: "Hirthehaven",
            county: "Saskatchewan",
            postcode: "S7R 9F9",
            country: "Canada",
            isAirmailRequired: true,
            phoneNumber: "072345678",
            email: "m.vancolkenburg@ca.test",
            companyName: "ACME",
          },
          {
            id: 4,
            salutation: "Dr",
            firstname: "Will",
            middlenames: "Oswald",
            surname: "Niesborella",
            dob: "01/07/1995",
            addressLine1: "47209 Stacey Plain",
            addressLine2: "Suite 113",
            addressLine3: "Devonburgh",
            town: "Marquardtville",
            county: "North Carolina",
            postcode: "40936",
            country: "United States",
            isAirmailRequired: true,
            phoneNumber: "0841781784",
            email: "docniesborella@mail.test",
            companyName: "ACME",
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/persons", "POST", {
      status: 201,
      body: {},
    });

    cy.visit("/select-or-create-correspondent?id=1&caseId=2");
  });

  it("can select an existing attorney to create a correspondent", () => {
    cy.contains("Add a correspondent");
    cy.get("label[for=f-attorney-1]").should(
      "contain.text",
      "Melanie Vanvolkenburg",
    );
    cy.get("label[for=f-attorney-2]")
      .should("contain.text", "Will Niesborella")
      .click();
    cy.get("button[type=submit]").click();
    cy.url().should("include", "/create-epa");
  });

  it("select create a new correspondent by default", () => {
    cy.contains("Add a correspondent");
    cy.get("input#f-attorney-new").should("be.checked");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "/create-correspondent");
  });
});
