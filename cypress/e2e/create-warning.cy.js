describe("Create a warning", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.setCookie("OPG-Bypass-Membrane", "1");
    cy.visit("/create-warning?id=189");
  });

  it("creates a warning", () => {
    cy.contains("Create Warning");
    cy.get(".moj-banner").should("not.exist");
    cy.get("select").select("Complaint Received");
    cy.get("textarea").type("Some warning notes");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
