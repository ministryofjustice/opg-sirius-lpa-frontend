describe("Edit a fee reduction", () => {
    beforeEach(() => {
        cy.visit("/edit-fee-reduction?id=124");
    });

    it("edits an existing fee reduction", () => {
        cy.contains("Edit fee reduction");
        cy.contains("7000-0000-0001");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#f-feeReductionType" ).select("Remission");
        cy.get("#f-paymentEvidence").type("Edited test evidence");
        cy.get("#f-paymentDate").type("2022-03-26");
        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
