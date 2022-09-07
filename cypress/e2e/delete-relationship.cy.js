describe("Delete a relationship", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.setCookie("OPG-Bypass-Membrane", "1");
    cy.visit("/delete-relationship?id=189");
  });

  it("deletes a relationship", () => {
    cy.contains("Delete Relationship");
    cy.contains("John Doe");
    cy.get(".moj-banner").should("not.exist");
    cy.contains("label", "John Doe").click();
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
