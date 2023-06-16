describe("Edit a document", () => {
  beforeEach(() => {
    cy.visit("/edit-document?id=800&case=lpa");
  });

  it("displays draft to edit", () => {
    cy.contains("Edit draft document");
    cy.contains("button", "Save draft");
    cy.contains("button", "Preview draft");
    cy.contains("button", "Delete draft");
    cy.contains("button", "Publish draft");
    cy.contains("button", "Save and exit");
  });

  it("can select a draft to edit", () => {
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    cy.get("#f-document").select(
      "1: 15/12/2022 13:41:04: Consuela Aysien: LP-A"
    );
    cy.contains("button", "Select").click();
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
  });

  it("saves a draft document", () => {
    cy.get("#documentTextEditor").contains("Test content");
    // can only edit the iframe
    const $iframe = cy
      .get("iframe[id=documentTextEditor_ifr]")
      .its("0.contentDocument.body")
      .should("not.be.empty")
      .then(cy.wrap);

    $iframe.clear().type("Edited test content");

    cy.contains("button", "Save draft").click();
    cy.get("#documentTextEditor").contains("Edited test content");
  });

  it("deletes a draft document", () => {
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    cy.get("#documentTextEditor").contains("Test content");

    cy.contains("button", "Delete draft").click();
  });

  it("previews a draft document", () => {
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    cy.get("#documentTextEditor").contains("Test content");

    cy.contains("button", "Preview draft").click();
  });

  it("publishes a draft document", () => {
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    cy.get("#documentTextEditor").contains("Test content");

    cy.contains("button", "Publish draft").click();
  });
});
