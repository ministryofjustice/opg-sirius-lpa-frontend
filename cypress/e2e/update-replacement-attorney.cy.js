describe("Update Replacement Attorney form on an LPA", () => {
    beforeEach(() => {
        cy.addMock("/lpa-api/v1/replacement-attorneys/3", "PUT", {
            status: 200,
            body: {},
        });

        cy.addMock("/lpa-api/v1/cases/2", "GET", {
            status: 200,
            body: {
                id: 2,
                caseSubtype: "pfa",
                receiptDate: "19/06/2026",
                replacementAttorneys: [
                    {
                        id: 3,
                        firstname: "Rudolph",
                        surname: "Stotesbury",
                    },
                ],
            },
        });

        cy.visit("/update-replacement-attorney?id=1&caseId=2&attorneyId=3");
    });

    it("updates an existing replacement attorney on an LPA", () => {
        cy.contains("Update replacement attorney details");
        cy.get("#f-firstname").should("have.value", "Rudolph");
        cy.get("#f-surname").should("have.value", "Stotesbury");
        cy.get("input[type=submit][name=add-another]").should("not.exist");

        cy.get("#f-firstname").clear().type("Rafael");
        cy.get("button[type=submit]").click();
        cy.url().should("include", "create-lpa");
    });
});
