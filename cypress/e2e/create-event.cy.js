describe("Create an event", () => {
  beforeEach(() => {
    cy.setCookie("Other", "other");
    cy.setCookie("XSRF-TOKEN", "abcde");
    cy.setCookie("OPG-Bypass-Membrane", "1");
    cy.visit("/create-event?id=800&entity=lpa");
  });

  it("creates an event", () => {
    cy.contains("Create Event");
    cy.contains("LPA 7000-0000-0000");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-type").select("Application processing");
    cy.get("#f-name").type("A title");
    cy.get("#f-description").type("More words");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });
});
