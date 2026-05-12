describe("Delete a note", () => {
  beforeEach(() => {
    cy.addMock(
      "/lpa-api/v1/persons/123/events?filter=eventId:456&sort=id:desc&limit=999",
      "GET",
      {
        status: 200,
        body: {
          events: [
            {
              sourceType: "Note",
              type: "Application processing",
              entity: {
                type: "Application processing",
                name: "Test note",
                description: "This is a test note",
                document: {
                  UUID: "123e4567-e89b-12d3-a456-426614174000",
                  friendlyDescription: "Test document",
                },
              },
              sourceNote: {
                id: 456,
              },
            },
          ],
        },
      },
    );
    cy.addMock("/lpa-api/v1/notes/456", "DELETE", {
      status: 204,
    });

    cy.visit("/delete-note?donorId=123&eventId=456");
  });

  it("deletes a note", () => {
    cy.contains("Delete this event?");
    cy.contains("Any documents attached to this event will NOT be deleted.");
    cy.get("button[type=submit]").click();
    cy.url().should("include", "history");
  });
});
