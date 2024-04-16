describe("Create Additional Digital LPA draft", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "private-mlpa"],
      },
    });

    cy.addMock("/lpa-api/v1/persons/130", "GET", {
      status: 200,
      body: {
        id: 130,
        dob: "05/05/1970",
        firstname: "John",
        surname: "Doe",
        uId: "7000-0000-0007",
        addressLine1: "Road",
        town: "City",
        postcode: "A12 3CD",
        country: "GB",
        cases: [
          {
            caseSubtype: "personal-welfare",
            caseType: "DIGITAL_LPA",
            id: 56,
            status: "Draft",
            uId: "M-ABCD-0000-1234",
          },
        ],
      },
    });

    cy.visit("/create-additional-draft-lpa?id=130");
  });

  it("creates an additional digital LPA using the same Donor address", () => {
    cy.addMock("/lpa-api/v1/donors/130/digital-lpas", "POST", {
      status: 201,
      body: [
        {
          caseSubtype: "personal-welfare",
          uId: "M-0101-ABCD-0101",
        },
      ],
    });
    cy.contains("Create a draft LPA for John Doe");

    cy.contains("Personal welfare").click();
    cy.contains("Property and affairs").click();
    cy.contains("The donor, using the address above").click();
    cy.contains("Confirm and create draft LPA").click();

    cy.get(".govuk-notification-banner__content").contains(
      "A draft LPA for John Doe has been created.",
    );
    cy.contains("M-0101-ABCD-0101");
  });

  it("creates 2 additional digital LPAs using a different donor address", () => {
    cy.addMock("/lpa-api/v1/donors/130/digital-lpas", "POST", {
      status: 201,
      body: [
        {
          caseSubtype: "personal-welfare",
          uId: "M-0101-ABCD-0101",
        },
        {
          caseSubtype: "property-and-affairs",
          uId: "M-0202-WXYZ-0202",
        },
      ],
    });
    cy.contains("Create a draft LPA for John Doe");

    cy.contains("Personal welfare").click();
    cy.contains("Property and affairs").click();
    cy.contains("The donor, using a different address").click();

    cy.contains("Enter a UK postcode, or enter the address manually")
      .nextUntil(".govuk-input")
      .next()
      .type("SW1A 1AA");
    cy.contains("Find address").click();

    cy.contains("Confirm and create draft LPA").click();

    cy.get(".govuk-notification-banner__content").contains(
      "2 draft LPAs for John Doe have been created.",
    );

    cy.contains("M-0101-ABCD-0101");
    cy.contains("M-0202-WXYZ-0202");
  });

  it("creates an additional digital LPA with a correspondent", () => {
    cy.addMock("/lpa-api/v1/donors/130/digital-lpas", "POST", {
      status: 201,
      body: [
        {
          caseSubtype: "property-and-affairs",
          uId: "M-0202-WXYZ-0202",
        },
      ],
    });
    cy.contains("Create a draft LPA for John Doe");

    cy.contains("Personal welfare").click();
    cy.contains("Property and affairs").click();
    cy.contains("Another person").click();
    cy.get("#f-correspondentFirstname").type("Simon");
    cy.get("#f-correspondentSurname").type("Sheldon");
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

    cy.contains("Confirm and create draft LPA").click();

    cy.get(".govuk-notification-banner__content").contains(
      "A draft LPA for John Doe has been created.",
    );
    cy.contains("M-0202-WXYZ-0202");
  });
});
