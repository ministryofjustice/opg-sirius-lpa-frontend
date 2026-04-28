describe("Delete a note", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/notes/456", "DELETE", {
      status: 204,
    });

    cy.visit("/delete-note?donorId=123&noteId=456");
  });

  it("deletes a note", () => {
    cy.contains("Delete this event?");
    cy.contains("Any documents attached to this event will NOT be deleted.");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "history");
  });
});
