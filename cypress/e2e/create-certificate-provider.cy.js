describe("Create certificate provider form", () => {
    beforeEach(() => {
        cy.addMock("/lpa-api/v1/cases/2", "GET", {
            status: 200,
            body: {
                id: 2,
            },
        });

        cy.addMock("/lpa-api/v1/persons", "POST", {
            status: 201,
            body: {},
        });

        cy.visit("/create-certificate-provider?id=1&caseId=2");
    });

    it("redirects to lpa form on submit", () => {
        cy.contains("Add a certificate provider");
        cy.get("#f-salutation").type("Prof");
        cy.get("#f-firstname").type("Melanie");
        cy.get("#f-middlenames").type("Josefina");
        cy.get("#f-surname").type("Vanvolkenburg");
        cy.get(".govuk-details__summary").click();
        cy.get("#f-addressLine1").type("29737 Andrew Plaza");
        cy.get("#f-addressLine2").type("Apt. 814");
        cy.get("#f-addressLine3").type("Gislasonside");
        cy.get("#f-town").type("Hirthehaven");
        cy.get("#f-county").type("Saskatchewan");
        cy.get("#f-postcode").type("S7R 9F9");
        cy.get("#f-country").type("Canada");
        cy.get("button[type=submit]").click();
        cy.url().should("include", "create-lpa");
    });

    it("redirects to certificate provider form when adding another", () => {
        cy.get("#f-salutation").type("Prof");
        cy.get("#f-firstname").type("Melanie");
        cy.get("#f-middlenames").type("Josefina");
        cy.get("#f-surname").type("Vanvolkenburg");
        cy.get(".govuk-details__summary").click();
        cy.get("#f-addressLine1").type("29737 Andrew Plaza");
        cy.get("#f-addressLine2").type("Apt. 814");
        cy.get("#f-addressLine3").type("Gislasonside");
        cy.get("#f-town").type("Hirthehaven");
        cy.get("#f-county").type("Saskatchewan");
        cy.get("#f-postcode").type("S7R 9F9");
        cy.get("#f-country").type("Canada");
        cy.get("input[type=submit][name=add-another]").click();
        cy.url().should("include", "create-certificate-provider");
        cy.contains("Add a certificate provider");
    });

    it("has a back link to the LPA form", () => {
        cy.get(".govuk-back-link")
            .should("exist")
            .and("have.attr", "href")
            .and(
                "include",
                "/create-lpa?id=1&caseId=2#accordion-create-lpa-heading-3",
            );
    });
});
