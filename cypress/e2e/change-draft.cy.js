describe("Change draft form", () => {
    beforeEach(() => {
        cy.addMock("/lpa-api/v1/digital-lpas/M-1111-2222-1110", "GET", {
            status: 200,
            body: {
                uId: "M-1111-2222-1110",
                "opg.poas.sirius": {
                    id: 565,
                    uId: "M-1111-2222-1110",
                    status: "Draft",
                    caseSubtype: "personal-welfare",
                    createdDate: "24/04/2024",
                    investigationCount: 0,
                    complaintCount: 0,
                    taskCount: 0,
                    warningCount: 0,
                    donor: {
                        id: 33,
                        firstname: "Peter",
                        surname: "MacPhearson",
                    },
                    application: {
                        source: "PHONE",
                        donorFirstNames: "Peter",
                        donorLastName: "MacPhearson",
                        donorDob: "27/05/1968",
                        donorEmail: "peter@test",
                        donorPhone: "073656249524",
                        donorAddress: {
                            addressLine1: "Flat 9999",
                            addressLine2: "Flaim House",
                            addressLine3: "33 Marb Road",
                            town: "Birmingham",
                            country: "GB",
                            postcode: "X15 3XX",
                        },
                        correspondentFirstNames: "Salty",
                        correspondentLastName: "McNab",
                        correspondentAddress: {
                            addressLine1: "Flat 3",
                            addressLine2: "Digital LPA Avenue",
                            addressLine3: "Noplace",
                            country: "GB",
                            postcode: "SW1 1AA",
                        },
                    },
                    linkedDigitalLpas: [],
                },
            },
        });

        cy.addMock("/lpa-api/v1/cases/565", "GET", {
            status: 200,
            body: {
                id: 565,
                uId: "M-1111-2222-1110",
                caseType: "DIGITAL_LPA",
                donor: {
                    id: 33,
                },
            },
        });

        cy.addMock("/lpa-api/v1/cases/565/warnings", "GET", {
            status: 200,
            body: [],
        });

        cy.addMock(
            "/lpa-api/v1/cases/565/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
            "GET",
            {
                status: 200,
                body: {
                    tasks: [],
                },
            },
        );
        cy.visit("/lpa/M-1111-2222-1110/change-draft");
    });

    it("can be visited from the LPA draft Change link", () => {
        cy.visit("/lpa/M-1111-2222-1110/lpa-details").then(() => {
            cy.get("#f-change-draft").click();
            cy.contains("Change donor details");
            cy.url().should(
                "contain",
                "/lpa/M-1111-2222-1110/change-draft",
            );
        });
    });

    it("populates draft details", () => {
        cy.get("#f-firstNames").should("have.value", "Peter");
        cy.get("#f-lastName").should("have.value", "MacPhearson");

        cy.get("#f-dob-day").should("have.value", "27");
        cy.get("#f-dob-month").should("have.value", "11");
        cy.get("#f-dob-year").should("have.value", "1971");

        cy.get("#f-address\\.Line1").should("have.value", "15 Cameron Approach");
        cy.get("#f-address\\.Line2").should("have.value", "Nether Collier");
        cy.get("#f-address\\.Line3").should("have.value", "");
        cy.get("#f-address\\.Town").should("have.value", "Worcestershire");
        cy.get("#f-address\\.Postcode").should("have.value", "BL2 6DI");
        cy.get("#f-address\\.Country").should("have.value", "GB");

        cy.get("#f-phoneNumber").should("have.value", "0123456789");
        cy.get("#f-email").should("have.value", "j@example.com");
    });

    it("can go Back to LPA details", () => {
        cy.contains("Back to LPA details").click();
        cy.url().should("contain", "/lpa/M-1111-2222-1110/lpa-details");
    });

    it("can be cancelled, returning to the LPA details", () => {
        cy.contains("Cancel").click();
        cy.url().should("contain", "/lpa/M-1111-2222-1110/lpa-details");
    });
});