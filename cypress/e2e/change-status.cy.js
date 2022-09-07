describe("Change status", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.setCookie("OPG-Bypass-Membrane", "1");
    cy.visit("/change-status?id=800&case=lpa");
  });

  it("chnges the case status", () => {
    cy.contains("Change status");
    cy.contains("LPA 7000-0000-0000");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-status").select("Perfect");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
