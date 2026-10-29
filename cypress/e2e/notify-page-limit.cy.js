describe("Notify page limit", () => {
  it("displays go to drafts btn", () => {
    cy.visit("/notify-error/678/dfef6714-b4fe-44c2-b26e-90dfe3663e95/LPA/8/700000000021?defaultDocumentUUID=ad8bbdc6-163d-49f9-9cac-a8141310e409");

    cy.get('.moj-alert__content').should('have.text', 'Document has been published but has not been posted via GOV.UK,Notify.\nThe published document has more than 10 pages and will not be posted automatically via GOV.UK,Notify. It must be manually uploaded to Xerox for posting.');

    cy.get('[data-cy="document-list-btn"]').should('have.text', 'Document list');
    cy.get('[data-cy="go-to-drafts-btn"]').should('have.text', 'Go to drafts');
    cy.get('[data-cy="exit-btn"]').should('have.text', 'Exit');
  });

  it("does not display go to drafts btn", () => {
    cy.visit("/notify-error/678/dfef6714-b4fe-44c2-b26e-90dfe3663e95/LPA/8/700000000021?defaultDocumentUUID=");

    cy.get('.moj-alert__content').should('have.text', 'Document has been published but has not been posted via GOV.UK,Notify.\nThe published document has more than 10 pages and will not be posted automatically via GOV.UK,Notify. It must be manually uploaded to Xerox for posting.');

    cy.get('[data-cy="document-list-btn"]').should('have.text', 'Document list');
    cy.get('[data-cy="go-to-drafts-btn"]').should('not.exist');
    cy.get('[data-cy="exit-btn"]').should('have.text', 'Exit');
  });
});
