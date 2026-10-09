describe("Create notified person form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons", "POST", {
      status: 201,
      body: {},
    });

    cy.addMock("/lpa-api/v1/cases/3", "GET", {
      status: 200,
      body: {
        id: 3,
        notifiedPersons: [
          {
            id: 11,
            salutation: "Prof",
            firstname: "Melanie",
            middlenames: "Josefina",
            surname: "Vanvolkenburg",
            addressLine1: "29737 Andrew Plaza",
            addressLine2: "Apt. 814",
            addressLine3: "Gislasonside",
            town: "Hirthehaven",
            county: "Saskatchewan",
            postcode: "S7R 9F9",
            country: "Canada",
            noticeGivenDate: "01/01/2023",
          },
          {
            id: 12,
          },
        ],
      },
    });

    cy.visit("/update-notified-person?id=1&caseId=3&notifiedPersonId=11");
  });

  it("edits a notified person when next notified person id is set", () => {
    cy.contains("Update notified person details");
    cy.get("#f-salutation").should("have.value", "Prof");
    cy.get("#f-firstname").should("have.value", "Melanie");
    cy.get("#f-middlenames").should("have.value", "Josefina");
    cy.get("#f-surname").should("have.value", "Vanvolkenburg");
    cy.get("#f-addressLine1").should("have.value", "29737 Andrew Plaza");
    cy.get("#f-addressLine2").should("have.value", "Apt. 814");
    cy.get("#f-addressLine3").should("have.value", "Gislasonside");
    cy.get("#f-town").should("have.value", "Hirthehaven");
    cy.get("#f-county").should("have.value", "Saskatchewan");
    cy.get("#f-postcode").should("have.value", "S7R 9F9");
    cy.get("#f-country").should("have.value", "Canada");
    cy.get("#f-noticeGivenDate").should("have.value", "2023-01-01");
    cy.get("input[type=submit][name=next-notified-person]").click();
    cy.url().should("include", "update-notified-person");
  });
});
