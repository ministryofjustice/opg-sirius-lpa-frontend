describe("Create a document", () => {
    beforeEach(() => {
        cy.visit("/create-document?id=800&case=lpa");
        cy.contains("Create Draft");
        cy.get(".moj-banner").should("not.exist");
        cy.get("#f-templateId").type("DD");
        cy.get(".autocomplete__menu").contains("DD: Donor deceased: Blank template").click();
        cy.contains("button", "Select template").click();
        cy.contains("Template: DD");
        // cy.get("#f-DD1").click();
        cy.contains("button", "Continue").click();
    });


    it("creates a document on the case by selecting a recipient", () => {
        cy.contains(".govuk-radios__item", "Select").find("input").check();
        cy.get("#f-189").click();
        cy.contains("button", "Select recipient").click();
    });

    it("generates a recipient", () => {
        cy.contains("Select or generate a recipient");
        cy.contains(".govuk-radios__item", "Generate").find("input").check();

        cy.get("#f-salutation").type("Prof");
        cy.get("#f-firstName").type("Melanie");
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

        cy.contains("button", "Create new recipient").click();
        cy.get(".moj-banner").should("exist");
        cy.get(".moj-banner").contains(
            "You have successfully created a recipient."
        );
    });
});
