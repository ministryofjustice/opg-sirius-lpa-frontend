describe("Unlink records", () => {
  beforeEach(() => {
    cy.visit("/unlink-person?id=189");
  });

  it("unlinks the persons records", () => {
    cy.contains("Unlink Record");
    cy.contains("John Doe");
    cy.get(".moj-banner").should("not.exist");
    cy.get("label[for=child-id-0]").click();
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
