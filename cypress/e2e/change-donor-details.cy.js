describe("Change donor details", () => {
    beforeEach(() => {
        cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0001", "GET", {
            status: 200,
            body: {
                uId: "M-0000-0000-0001",
                "opg.poas.sirius": {
                    id: 666,
                    uId: "M-0000-0000-0001",
                    status: "Draft",
                    caseSubtype: "personal-welfare",
                    createdDate: "31/10/2023",
                    investigationCount: 0,
                    complaintCount: 0,
                    taskCount: 0,
                    warningCount: 0,
                    donor: {
                        id: 33,
                    },
                    application: {
                        donorFirstNames: "Agnes",
                        donorLastName: "Hartley",
                        donorDob: "27/05/1998",
                        donorEmail: "agnes@host.example",
                        donorPhone: "073656249524",
                        donorAddress: {
                            addressLine1: "Apartment 3",
                            addressLine2: "Gherkin Building",
                            addressLine3: "33 London Road",
                            country: "GB",
                            postcode: "B15 3AA",
                            town: "Birmingham",
                        },
                        correspondentFirstNames: "Kendrick",
                        correspondentLastName: "Lamar",
                        correspondentAddress: {
                            addressLine1: "Flat 3",
                            addressLine2: "Digital LPA Lane",
                            addressLine3: "Somewhere",
                            country: "GB",
                            postcode: "SW1 1AA",
                            town: "London",
                        },
                    },
                },
                "opg.poas.lpastore": {
                    donor: {
                        uid: "5ff557dd-1e27-4426-9681-ed6e90c2c08d",
                        firstNames: "James",
                        lastName: "Rubin",
                        otherNamesKnownBy: "Somebody",
                        dateOfBirth: "1990-02-22",
                        address: {
                            postcode: "B29 6BL",
                            country: "GB",
                            town: "Birmingham",
                            line1: "29 Grange Road"
                        },
                        contactLanguagePreference: "en",
                        email: "jrubin@mail.example"
                    },
                    attorneys: [
                        {
                            firstNames: "Esther",
                            lastName: "Greenwood",
                            status: "active",
                        },
                        {
                            firstNames: "Rico",
                            lastName: "Welch",
                            status: "replacement",
                            signedAt: "2022-12-19T09:12:59Z",
                        },
                    ],
                    certificateProvider: {
                        uid: "e4d5e24e-2a8d-434e-b815-9898620acc71",
                        firstNames: "Timothy",
                        lastNames: "Turner",
                        signedAt: "2022-12-18T11:46:24Z",
                    },
                    signedAt: "2022-12-18T11:46:24Z",
                    lpaType: "pw",
                    channel: "online",
                    registrationDate: "2022-12-18",
                    peopleToNotify: [],
                },
            },
        });


        cy.addMock("/lpa-api/v1/cases/666", "GET", {
            status: 200,
            body: {
                id: 666,
                uId: "M-0000-0000-0001",
                caseType: "DIGITAL_LPA",
                donor: {
                    id: 666,
                },
                status: "Draft",
            },
        });

        cy.addMock("/lpa-api/v1/cases/666/warnings", "GET", {
            status: 200,
            body: [],
        });

        cy.addMock(
            "/lpa-api/v1/cases/666/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
            "GET",
            {
                status: 200,
                body: {
                    tasks: [],
                },
            },
        );


        cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0001/change-donor-details", "PUT", {
            status: 204,
        });

        cy.visit("/change-donor-details?uid=M-0000-0000-0001");
    });

    it("Can edit the change donor details form", () => {
        cy.contains("Change donor details");
        cy.contains("Details that apply to all LPAs for this donor");
        cy.get(".moj-banner").should("not.exist");

        cy.get("#f-firstNames").clear()
        cy.get("#f-firstNames").type("Coleen Stephanie");

        cy.get("#f-lastName").clear();
        cy.get("#f-lastName").type("Morneault");

        cy.get("#f-dob-day").clear();
        cy.get("#f-dob-month").clear();
        cy.get("#f-dob-year").clear();

        cy.get("#f-dob-day").type("8");
        cy.get("#f-dob-month").type("4");
        cy.get("#f-dob-year").type("1952");

        cy.get("button[type=submit]").click();
        cy.get(".moj-banner").should("exist");
    });
});
