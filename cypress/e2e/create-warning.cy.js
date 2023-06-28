describe("Create a warning", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/189", "GET", {
      status: 200,
      body: {
        dob: "05/05/1970",
        firstname: "John",
        id: 189,
        surname: "Doe",
        uId: "7000-0000-0007",
      },
    });
    cy.visit("/create-warning?id=189");
  });

  it("creates a warning", () => {
    cy.addMock("/lpa-api/v1/warnings", "POST", {
      status: 201,
      body: {
        personId: 189,
        warningText: "Some warning notes",
        warningType: "Complaint Received",
      },
    });
    cy.contains("Create Warning");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#case-id-0").should("not.exist");
    cy.get("select").select("Complaint Received");
    cy.get("textarea").type("Some warning notes");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});

describe("Create a warning on multiple cases", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/400", "GET", {
      status: 200,
      body: {
        dob: "05/05/1970",
        firstname: "John",
        id: 400,
        surname: "Doe",
        uId: "7000-0000-0007",
        cases: [
          {
            caseSubtype: "pfa",
            caseType: "LPA",
            id: 405,
            status: "Perfect",
            uId: "7000-5382-4438",
          },
          {
            caseSubtype: "hw",
            caseType: "LPA",
            id: 406,
            status: "Pending",
            uId: "7000-5382-8764",
          },
        ],
      },
    });

    cy.visit("/create-warning?id=400");
  });

  it("creates a warning on multiple cases", () => {
    cy.addMock("/lpa-api/v1/warnings", "POST", {
      status: 201,
      body: {
        personId: 400,
        warningText: "Some warning notes for multiple cases",
        warningType: "Complaint Received",
        caseIds: [405, 406],
      },
    });

    cy.contains("Create Warning");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#case-id-0").click();
    cy.get("#case-id-1").click();
    cy.get("select").select("Complaint Received");
    cy.get("textarea").type("Some warning notes for multiple cases");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
