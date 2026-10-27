describe("Create a document", () => {
    beforeEach(() => {
        cy.visit("/notify-error/678/dfef6714-b4fe-44c2-b26e-90dfe3663e95/LPA/8/700000000021?defaultDocumentUUID=ad8bbdc6-163d-49f9-9cac-a8141310e409");
    });

    it("creates a document on the case by selecting a recipient", () => {
        cy.get('[data-cy="notify-warning-message"]').contains("Document has been published named Dr Consuela Aysien - LPA perfect + reg due date: applicant. The published document has more than 10 pages and will not be sent automatically via GOV.UK,Notify. It must be manually uploaded to Xerox for sending.")
        cy.get('[data-cy="notify-download-document-btn"]').should("exist").contains("Download document");
        cy.get('[data-cy="notify-go-to-drafts-btn"]').should("exist").contains("Go to drafts");
    });
});
