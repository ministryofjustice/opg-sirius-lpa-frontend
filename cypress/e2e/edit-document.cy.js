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
        cy.get("#f-document").select("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
        cy.contains("button", "Select").click();
        cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    });

    it("saves a draft document", () => {
        cy.get("#documentTextEditor").contains("Test content");
        // can only edit the iframe
        cy.get("#documentTextEditor_ifr").type("Edited test content");

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

        // Make the link open in a new tab because otherwise Cypress
        // throws cross-domain errors
        cy.contains("a", "Download").then(($input) => {
            $input[0].setAttribute("target", "_blank");
        });

        cy.contains("a", "Download").click();
        cy.contains("a", "Download")
            .invoke("attr", "class")
            .should("contain", "govuk-button--disabled");
        cy.contains("Your download will open in a new window when ready");
    });

    it("publishes a draft document", () => {
        cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
        cy.get("#documentTextEditor").contains("Test content");

        cy.contains("button", "Publish draft").click();
    });
});
