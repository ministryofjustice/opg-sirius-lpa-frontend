describe("Create a document for a digital LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-GDJ7-QK9R-4XVF", "GET", {
      status: 200,
      body: {
        "opg.poas.sirius": {
          id: 483,
          donor: {
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
        "opg.poas.lpastore": {
          donor: {
            uid: "572fe550-e465-40b3-a643-ca9564fabab8",
            firstNames: "Steven",
            lastName: "Munnell",
            email: "Steven.Munnell@example.com",
            dateOfBirth: "17/06/1982",
            otherNamesKnownBy: "",
            contactLanguagePreference: "",
            address: {
              line1: "1 Scotland Street",
              line2: "Netherton",
              line3: "Glasgow",
              town: "Edinburgh",
              postcode: "EH6 18J",
              country: "GB",
            },
          }
        },
      },
    });

    cy.addMock("/lpa-api/v1/templates/digitallpa", "GET", {
      status: 200,
      body: {
        DD: {
          label: "DLPA Example Form",
          inserts: {
            all: {
              DLPA_INSERT_01: {
                label: "DLPA Insert 1",
                order: 0,
              },
            },
          },
        },
      },
    });

    cy.addMock("/lpa-api/v1/lpas/483/documents/draft", "POST", {
      status: 201,
      body: {},
    });

    cy.visit("/lpa/M-GDJ7-QK9R-4XVF/documents/new");
  });

  it("creates a document on the case", () => {
    cy.contains("Select a document template");
    cy.get("#f-templateId").type("DLPA");
    cy.get(".autocomplete__menu").contains("DLPA Example Form").click();

    cy.contains("Select document inserts");
    cy.contains("DLPA Insert 1").click();

    cy.contains(
      "1 Scotland Street, Netherton, Glasgow, Edinburgh, EH6 18J, GB",
    );
    cy.contains("Steven Munnell (Donor)").click();

    cy.contains("button", "Continue").click();
  });
});
