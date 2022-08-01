describe("Unlink records", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.setCookie("OPG-Bypass-Membrane", "1");
    cy.visit("/unlink-person?id=189");
  });

  it("unlinks the persons records", () => {
    cy.contains("Unlink Record");
    cy.contains("John Doe");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#child-id-0").click();
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
