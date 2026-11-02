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
        replacementAttorneys: [
          {
            id: 5,
            salutation: "Mr",
            firstname: "Rupert",
            surname: "Replacement",
            systemStatus: false,
          },
        ],
        certificateProviders: [
          {
            id: 6,
            salutation: "Ms",
            firstname: "Cora",
            surname: "Provider",
            isAirmailRequired: true,
          },
        ],
        correspondent: {
          id: 4,
          salutation: "Mr",
          firstname: "Test",
          middlenames: "J",
          surname: "Name",
        },
        applicants: [{ id: 2 }],
      },
    });
  });

  it("lists the people involved in the case", () => {
    cy.visit("/sirius-header-people-info?id=123");

    cy.contains("h2", "Details on the people involved in this case");

    cy.contains("tr", "Donor").should(
      "contain.text",
      "Prof Melanie Josefina Vanvolkenburg",
    );
    cy.contains("tr", "Attorney 1").should(
      "contain.text",
      "Miss Amanda Percy Jenkins",
    );
    cy.contains("tr", "Attorney 2").should(
      "contain.text",
      "Dr Will Oswald Niesborella",
    );
    cy.contains("tr", "Replacement attorney 1").should(
      "contain.text",
      "Mr Rupert Replacement",
    );
    cy.contains("tr", "Certificate provider").should(
      "contain.text",
      "Ms Cora Provider",
    );
    cy.contains("tr", "Correspondent").should("contain.text", "Mr Test J Name");
  });

  it("shows status, applicant and airmail indicators", () => {
    cy.visit("/sirius-header-people-info?id=123");

    cy.contains("tr", "Donor").within(() => {
      cy.get('svg[aria-label="Airmail"]').should("exist");
      cy.contains("Active").should("not.exist");
      cy.contains("Inactive").should("not.exist");
    });

    cy.contains("tr", "Attorney 1").within(() => {
      cy.get('[title="Applicant"]').should("exist");
      cy.contains("Active");
    });

    cy.contains("tr", "Attorney 2").within(() => {
      cy.get('[title="Applicant"]').should("not.exist");
      cy.contains("Active");
    });

    cy.contains("tr", "Replacement attorney 1").within(() => {
      cy.contains("Inactive");
    });

    cy.contains("tr", "Certificate provider").within(() => {
      cy.get('svg[aria-label="Airmail"]').should("not.exist");
      cy.contains("Active").should("not.exist");
      cy.contains("Inactive").should("not.exist");
    });
  });

  it("shows the donor's details by default", () => {
    cy.visit("/sirius-header-people-info?id=123");

    cy.get("tr.selected").should("contain.text", "Donor");
    cy.get(".app-person-info__details-heading").should(
      "contain.text",
      "Donor: Prof Melanie Josefina Vanvolkenburg",
    );

    cy.get(".app-person-info__details").within(() => {
      cy.contains("dt", "Company name").next().should("have.text", "ACME");
      cy.contains("dt", "Date of birth")
        .next()
        .should("have.text", "19/04/1978");
      cy.contains("dt", "Address")
        .next()
        .should(
          "have.text",
          "29737 Andrew Plaza, Apt. 814, Gislasonside, Hirthehaven, Saskatchewan, Canada",
        );
      cy.contains("dt", "Postcode").next().should("have.text", "S7R 9F9");
      cy.contains("dt", "Telephone number")
        .next()
        .should("have.text", "072345678");
      cy.contains("dt", "Email address")
        .next()
        .find("a")
        .should("have.attr", "href", "mailto:m.vancolkenburg@ca.test");
      cy.contains("AIRMAIL").should("not.exist");
    });
  });

  it("shows the selected person's details with empty rows still labelled", () => {
    cy.visit("/sirius-header-people-info?id=123&selected=2");

    cy.get("tr.selected").should("contain.text", "Attorney 1");
    cy.get(".app-person-info__details-heading").should(
      "contain.text",
      "Attorney 1: Miss Amanda Percy Jenkins",
    );

    cy.get(".app-person-info__details").within(() => {
      [
        "Date of birth",
        "Address",
        "Postcode",
        "Telephone number",
        "Email address",
      ].forEach((label) => {
        cy.contains("dt", label).next().invoke("text").should("match", /^\s*$/);
      });
    });
  });
});
