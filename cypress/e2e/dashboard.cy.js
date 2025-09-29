import * as digitalLpas from "../mocks/digitalLpas";
import * as cases from "../mocks/cases";

describe("The Digital LPA Dashboard", () => {
    beforeEach(() => {
        cy.addMock("/lpa-api/v1/tasks/1", "GET", {
            status: 200,
            body: {
                caseItems: [
                    {
                        caseType: "DIGITAL_LPA",
                        uId: "M-DIGI-LPA3-3333",
                    },
                ],
                dueDate: "10/01/2022",
                id: 1,
                name: "Create physical case file",
                status: "Not Started",
            },
        });

        const mocks = Promise.allSettled([
            digitalLpas.get("M-DIGI-LPA3-3333", {
                "opg.poas.sirius": {
                    id: 333,
                    donor: {
                        id: 33,
                    },
                    linkedDigitalLpas: [
                        {
                            uId: "M-DIGI-LPA3-3334",
                            caseSubtype: "personal-welfare",
                            status: "Draft",
                            createdDate: "01/11/2023",
                        },
                        {
                            uId: "M-DIGI-LPA3-3335",
                            caseSubtype: "personal-welfare",
                            status: "Registered",
                            createdDate: "02/11/2023",
                        },
                    ],
                },
            }),
            digitalLpas.progressIndicators.feesInProgress("M-DIGI-LPA3-3333"),
            cases.warnings.empty("333"),
            digitalLpas.objections.empty("M-DIGI-LPA3-3333"),
        ]);

        cy.wrap(mocks);

        cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");
    });

    it("shows case information", () => {
        cy.get("h1").contains("Steven Munnell");

        cy.contains("M-DIGI-LPA3-3333");
        cy.get("a[href='/lpa/M-DIGI-LPA3-3333'] .govuk-tag").contains("Draft");

        cy.contains("PW M-DIGI-LPA3-3334");
        cy.get("a[href='/lpa/M-DIGI-LPA3-3334'] .govuk-tag").contains("Draft");

        cy.contains("PW M-DIGI-LPA3-3335");
        cy.get("a[href='/lpa/M-DIGI-LPA3-3335'] .govuk-tag").contains("Registered");
    });

    it("shows task table", () => {
        cy.get(
            "table[data-role=tasks-table] [data-role=tasks-table-header] tr th",
        ).should((elts) => {
            expect(elts).to.contain("Tasks");
            expect(elts).to.contain("Due date");
            expect(elts).to.contain("Actions");
        });
        cy.get(
            "table[data-role=tasks-table] tr[data-role=tasks-table-task-row]",
        ).should((elts) => {
            expect(elts).to.have.length(3);
            expect(elts).to.contain("Review reduced fee eligibility");
            expect(elts).to.contain("Review application correspondence");
            expect(elts).to.contain("Another task");
            expect(elts).to.contain("Reassign task");
        });
    });

    it("shows warnings list", () => {
        cy.addMock("/lpa-api/v1/cases/333/warnings", "GET", {
            status: 200,
            body: [
                {
                    id: 44,
                    warningType: "Court application in progress",
                    warningText: "Court notified",
                    dateAdded: "24/08/2022 13:13:13",
                    caseItems: [
                        { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
                    ],
                },
                {
                    id: 22,
                    warningType: "Complaint Received",
                    warningText: "Complaint from donor",
                    dateAdded: "12/12/2023 12:12:12",
                    caseItems: [
                        { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
                        { uId: "M-DIGI-LPA3-5555", caseSubtype: "property-and-affairs" },
                    ],
                },
                {
                    id: 24,
                    warningType: "Donor Deceased",
                    warningText: "Advised of donor death",
                    dateAdded: "05/01/2022 10:10:00",
                    caseItems: [
                        { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
                        { uId: "M-DIGI-LPA3-5555", caseSubtype: "property-and-affairs" },
                        { uId: "M-DIGI-LPA3-6666", caseSubtype: "personal-welfare" },
                    ],
                },
            ],
        });

        cy.visit("/lpa/M-DIGI-LPA3-3333");

        cy.get(".app-caseworker-summary > div:nth-child(2) li").should((elts) => {
            expect(elts).to.have.length(3);

            expect(elts[0]).to.contain("Donor Deceased");
            expect(elts[0]).to.contain(
                "this case, PA M-DIGI-LPA3-5555 and PW M-DIGI-LPA3-6666",
            );

            expect(elts[1]).to.contain("Complaint Received");
            expect(elts[1]).to.contain("this case and PA M-DIGI-LPA3-5555");

            expect(elts[2]).to.contain("Court application in progress");
            expect(elts[2]).not.to.contain("this case");
        });
    });

    it("can cancel reassign task", () => {
        cy.contains("Reassign task").click();
        cy.url().should("include", "/assign-task?id=1");
        cy.contains("Assign Task");

        cy.contains("Cancel").click();
        cy.url().should("include", "/lpa/M-DIGI-LPA3-3333");
        cy.contains("Case summary");
    });

    it("can cancel clear task", () => {
        cy.contains("Clear task").click();
        cy.url().should("include", "/clear-task?id=1");
        cy.contains("Save and clear task");

        cy.contains("Cancel").click();
        cy.url().should("include", "/lpa/M-DIGI-LPA3-3333");
        cy.contains("Case summary");
    });

    it("can clear a task", () => {
        cy.addMock("/lpa-api/v1/tasks/1/mark-as-completed", "PUT", {
            status: 200,
        });

        cy.contains("Clear task").click();
        cy.url().should("include", "/clear-task?id=1");
        cy.get("button[type=submit]").click();

        cy.get(".moj-alert").should("exist");
        cy.get(".moj-alert").contains("Task completed");

        cy.url().should("contain", "/lpa/M-DIGI-LPA3-3333");
        cy.contains("Case summary");
    });
});
