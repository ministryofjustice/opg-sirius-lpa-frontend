describe("Change draft form", () => {
    beforeEach(() => {
        cy.addMock("/lpa-api/v1/digital-lpas/M-1111-2222-1111", "GET", {
            status: 200,
            body: {
                uId: "M-1111-2222-1111",
                "opg.poas.sirius": {
                    id: 565,
                    uId: "M-1111-2222-1111",
                    status: "Draft",
                    caseSubtype: "personal-welfare",
                    createdDate: "24/04/2024",
                    investigationCount: 0,
                    complaintCount: 0,
                    taskCount: 0,
                    warningCount: 0,
                    donor: {
                        id: 44,
                        firstname: "Peter",
                        surname: "MacPhearson",
                        dob: "1971-11-27",
                        addressLine1: "15 Cameron Approach",
                        addressLine2: "Nether Collier",
                        addressLine3: "",
                        town: "Worcestershire",
                        postcode: "BL2 6DI",
                        country: "GB",
                        phone: "0123456789",
                        email: "peter@test.com",
                    },
                    application: {
                        source: "PHONE",
                        donorFirstNames: "Peter",
                        donorLastName: "MacPhearson",
                        donorDob: "27/05/1968",
                        donorEmail: "peter@test.com",
                        donorPhone: "0123456789",
                        donorAddress: {
                            addressLine1: "15 Cameron Approach",
                            addressLine2: "Nether Collier",
                            addressLine3: "",
                            town: "Worcestershire",
                            country: "GB",
                            postcode: "BL2 6DI",
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
                uId: "M-1111-2222-1111",
                caseType: "DIGITAL_LPA",
                donor: {
                    id: 44,
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

    });

    it("can be visited from the LPA draft Change link", () => {
        cy.visit("/lpa/M-1111-2222-1111/lpa-details").then(() => {
            cy.get("#f-change-draft").click();
            cy.contains("Change donor details");
            cy.url().should(
                "contain",
                "/lpa/M-1111-2222-1111/change-draft",
            );
        });
    });

    it("populates draft details", () => {
        cy.visit("/lpa/M-1111-2222-1111/change-draft");

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
        cy.get("#f-email").should("have.value", "peter@test.com");
    });

    it("can go Back to LPA details", () => {
        cy.visit("/lpa/M-1111-2222-1111/change-draft");
        cy.contains("Back to LPA details").click();
        cy.url().should("contain", "/lpa/M-1111-2222-1111/lpa-details");
    });

    it("can be cancelled, returning to the LPA details", () => {
        cy.visit("/lpa/M-1111-2222-1111/change-draft");
        cy.contains("Cancel").click();
        cy.url().should("contain", "/lpa/M-1111-2222-1111/lpa-details");
    });
});