describe("Person info panel on the header bar", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/123", "GET", {
      status: 200,
      body: {
        donor: {
          id: 1,
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
        attorneys: [
          {
            id: 2,
            salutation: "Miss",
            firstname: "Amanda",
            middlenames: "Percy",
            surname: "Jenkins",
            systemStatus: true,
          },
          {
            id: 3,
            salutation: "Dr",
            firstname: "Will",
            middlenames: "Oswald",
            surname: "Niesborella",
            systemStatus: true,
          },
        ],
        correspondent: {
          id: 4,
          salutation: "Mr",
          firstname: "Test",
          middlenames: "J",
          surname: "Name",
        },
      },
    });

    cy.visit("/sirius-header-people-info?id=123");
  });

  it("displays the person info panel", () => {
    cy.contains("Donor:");
    cy.contains("Prof Melanie Josefina Vanvolkenburg");
    cy.contains("Attorney 1:");
    cy.contains("Miss Amanda Percy Jenkins");
    cy.contains("Attorney 2:");
    cy.contains("Dr Will Oswald Niesborella");
    cy.contains("Correspondent:");
    cy.contains("Mr Test J Name");
    cy.contains("Company Name:");
    cy.contains("ACME");
    cy.contains("DOB:");
    cy.contains("19/04/1978");
    cy.contains("Address:");
    cy.contains(
      "AIRMAIL 29737 Andrew Plaza, Apt. 814, Gislasonside, Hirthehaven, Saskatchewan, S7R 9F9, Canada",
    );
    cy.contains("Tel:");
    cy.contains("072345678");
    cy.contains("Email:");
    cy.contains("m.vancolkenburg@ca.test");
  });
});
