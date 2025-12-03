describe("View documents", () => {
    beforeEach(() => {
        cy.addMock("/lpa-api/v1/persons/1/cases", "GET", {
            status: 200,
            body: {
                cases: [
                    {
                        caseType: "LPA",
                        caseSubtype: "pfa",
                        id: 34,
                        uId: "7000-1234-1234",
                    },
                    {
                        caseType: "LPA",
                        caseSubtype: "hw",
                        id: 78,
                        uId: "7000-5678-5678",
                    },
                    {
                        caseType: "EPA",
                        caseSubtype: "hw",
                        id: 990,
                        uId: "7001-0000-5678",
                    },
                ],
            },
        });

        cy.visit("/donor/1/documents");
    });

    it("has title", () => {
        cy.contains("Documents");
    });

    it("looks at documents with a single uid", () => {
        cy.visit("/donor/1/documents?uid[]=7000-1234-1234");

        cy.contains("Documents");
    });

    it("looks at documents with multiple uid", () => {
        cy.visit("/donor/1/documents?uid[]=7000-1234-1234&uid[]=7000-5678-5678");

        cy.contains("Documents");
    });


    it("looks at EPA documents", () => {
        cy.visit("/donor/1/documents?uid[]=7001-0000-5678");

        cy.contains("Documents");
    });
});
