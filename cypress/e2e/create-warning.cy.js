describe("Create a warning", () => {
  beforeEach(() => {
    cy.visit("/create-warning?id=189");
  });

  it("creates a warning", () => {
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
    cy.visit("/create-warning?id=400");
  });

  it("creates a warning on multiple cases", () => {
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
