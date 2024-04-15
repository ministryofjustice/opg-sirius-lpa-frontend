describe("Create Additional Digital LPA draft", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "private-mlpa"],
      },
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-GDJ7-QK9R-4XVF", "GET", {
      status: 200,
      body: {
        "opg.poas.sirius": {
          id: 483,
          donor: {
            id: 12,
            firstname: "Steven",
            surname: "Munnell",
            dob: "17/06/1982",
            addressLine1: "1 Scotland Street",
            addressLine2: "Netherton",
            addressLine3: "Glasgow",
            town: "Edinburgh",
            postcode: "EH6 18J",
            country: "GB",
            personType: "Donor",
          },
          application: {
            donorFirstNames: "Steven",
            donorLastName: "Munnell",
            donorDob: "17/06/1982",
            donorAddress: {
              addressLine1: "1 Scotland Street",
              postcode: "EH6 18J",
            },
          },
        },
      },
    });

    cy.visit("/create-additional-draft-lpa?id=12");
  });

  it("creates an additional digital LPA", () => {
    cy.contains("Create a draft LPA for Steven Munnell");

    cy.contains("Personal welfare").click();
    cy.contains("Property and affairs").click();
    cy.contains("The donor, using the address above").click();

    cy.contains("Confirm and create draft LPA").click();

    cy.get(".govuk-panel").contains(
      "2 draft LPAs for Steven Munnell have been created.",
    );
  });
});
