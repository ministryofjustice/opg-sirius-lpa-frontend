describe("Create an event", () => {
  beforeEach(() => {
    cy.visit("/create-event?id=800&entity=lpa");
  });

  it("creates an event", () => {
    cy.contains("Create Event");
    cy.contains("LPA 7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");
    cy.get("#f-type").select("Application processing");
    cy.get("#f-name").type("Something");
    cy.get("#f-description").type("More words");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });
});
