describe("Create a document", () => {
    beforeEach(() => {
        cy.visit("/create-document?id=800&case=lpa");
        cy.contains("7000-0000-0000");
        cy.contains("Select a document template");
        cy.get("#f-templateId").type("DD");
        cy.get(".autocomplete__menu").contains("DD: Donor deceased: Blank template").click();
        cy.contains("button", "Continue").click();

        cy.contains("Select document inserts");
        cy.contains("DD1: DD1 - Case complete");
        cy.get("#f-DD1").click();
        cy.contains("button", "Continue").click();
    });


    it("creates a document on the case by selecting a recipient", () => {
        cy.contains("Select a recipient");
        cy.get("#f-189").click();
        cy.contains("button", "Continue").click();
    });

    it("creates a new recipient via new recipient form", () => {
        cy.contains("Select a recipient");
        cy.contains("button", "Add new recipient").click();

        cy.contains("Add a new recipient");
        cy.get("#f-salutation").type("Prof");
        cy.get("#f-firstname").type("Melanie");
        cy.get("#f-middlenames").type("Josefina");
        cy.get("#f-surname").type("Vanvolkenburg");
        cy.get("#f-addressLine1").type("29737 Andrew Plaza");
        cy.get("#f-addressLine2").type("Apt. 814");
        cy.get("#f-addressLine3").type("Gislasonside");
        cy.get("#f-town").type("Hirthehaven");
        cy.get("#f-county").type("Saskatchewan");
        cy.get("#f-postcode").type("S7R 9F9");
        cy.get("#f-isAirmailRequired").click();
        cy.get("#f-phoneNumber").type("072345678");
        cy.get("#f-email").type("m.vancolkenburg@ca.test");
        cy.get("#f-correspondenceBy-email").click();
        cy.get("#f-correspondenceBy-phone").click();

        cy.contains("button", "Continue").click();
        cy.get(".moj-banner").should("exist");
        cy.get(".moj-banner").contains("New recipient added");
    });
});
