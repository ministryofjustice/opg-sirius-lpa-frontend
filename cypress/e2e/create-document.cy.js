// describe("Create a document", () => {
//     beforeEach(() => {
//         cy.visit("/create-document?id=800&case=lpa");
//     });
//
//     it("creates a document on the case", () => {
//         cy.contains("Create Draft");
//         cy.contains("7000-0000-0000");
//         cy.get(".moj-banner").should("not.exist");
//         cy.get("#f-title").type("A title");
//         cy.get("#f-information").type("Some information");
//         cy.contains(".govuk-radios__label", "Priority").click();
//         cy.get("#f-dateReceived").type("2022-03-04");
//         cy.get("button[type=submit]").click();
//         cy.get(".moj-banner").should("exist");
//     });
// });
