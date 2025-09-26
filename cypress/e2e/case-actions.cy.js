import * as digitalLpas from "../mocks/digitalLpas";

describe("Case actions drop down", () => {
    beforeEach(() => {

        const mocks = Promise.allSettled([
            digitalLpas.get("M-1111-1111-1111"),
            digitalLpas.objections.empty("M-1111-1111-1111"),
            digitalLpas.progressIndicators.feesInProgress("M-1111-1111-1111"),
        ]);

        cy.wrap(mocks);

        cy.addMock(
            "/lpa-api/v1/cases/1111/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
            "GET",
            {
                status: 200,
                body: {
                    tasks: [
                        {
                            id: 1,
                            name: "Review restrictions and conditions",
                            duedate: "10/12/2023",
                            status: "OPEN",
                            assignee: { displayName: "Super Team" },
                        },
                    ],
                },
            },
        );

        cy.addMock("/lpa-api/v1/persons/1111/cases", "GET", {
            status: 200,
            body: {
                cases: [
                    {
                        caseSubtype: "property-and-affairs",
                        id: 1111,
                        uId: "M-1111-1111-1111",
                        status: "Draft",
                    },
                ],
            },
        });

        cy.addMock("/lpa-api/v1/tasks/1", "GET", {
            status: 200,
            body: {
                caseItems: [
                    {
                        caseType: "DIGITAL_LPA",
                        uId: "M-1111-1111-1111",
                    },
                ],
                dueDate: "10/01/2022",
                id: 1,
                name: "Create physical case file",
                status: "Not Started",
            },
        });


        cy.addMock("/lpa-api/v1/cases/1111/tasks", "POST", {
            status: 201,
            body: { tasks: [] },
        });

        cy.addMock("/lpa-api/v1/cases/1111", "GET", {
            status: 200,
            body: {
                id: 1111,
                uId: "M-1111-1111-1111",
                caseType: "DIGITAL_LPA",
                donor: {
                    id: 1111,
                },
                status: "Processing",
                expectedPaymentTotal: 8200,
            },
        });

        cy.addMock("/lpa-api/v1/cases/1111/warnings", "GET", {
            status: 200,
            body: [
                {
                    id: 44,
                    warningType: "Court application in progress",
                    warningText: "Court notified",
                    dateAdded: "24/08/2022 13:13:13",
                    caseItems: [
                        { uId: "M-1111-1111-1111", caseSubtype: "property-and-affairs" },
                    ],
                },
            ],
        });

        cy.addMock("/lpa-api/v1/warnings", "POST", {
            status: 201,
            body: {
                personId: 1111,
                warningText: "Be warned!",
                warningType: "Complaint Received",
                caseIds: [1111],
            },
        });

        cy.addMock("/lpa-api/v1/warnings", "POST", {
            status: 201,
        });

        cy.visit("/lpa/M-1111-1111-1111");
    });

    it("can create a task", () => {
        cy.contains(".govuk-button", "Case actions").click();
        cy.contains("Create a task").click();
        cy.url().should("include", "/create-task?id=1111");
        cy.contains("M-1111-1111-1111");
        cy.get("#f-taskType").select("Check Application");
        cy.get("#f-name").type("Do this task");
        cy.get("#f-description").type("This task, do");
        cy.contains("label", "Team").click();
        cy.get("#f-assigneeTeam").select("Cool Team");
        cy.get("#f-dueDate").type("2035-01-01");
        cy.get("button[type=submit]").click();

        cy.get(".moj-alert").should("exist");
        cy.get(".moj-alert").contains("Task created");
        cy.get("h1").contains("Steven Munnell");
        cy.location("pathname").should("eq", "/lpa/M-1111-1111-1111");
    });

    it("can cancel reassign task", () => {
        cy.contains("Reassign task").click();
        cy.url().should("include", "/assign-task?id=1");
        cy.contains("Assign Task");

        cy.contains("Cancel").click();
        cy.url().should("include", "/lpa/M-1111-1111-1111");
        cy.contains("Case summary");
    });

    it("can cancel clear task", () => {
        cy.contains("Clear task").click();
        cy.url().should("include", "/clear-task?id=1");
        cy.contains("Save and clear task");

        cy.contains("Cancel").click();
        cy.url().should("include", "/lpa/M-1111-1111-1111");
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

        cy.url().should("contain", "/lpa/M-1111-1111-1111");
        cy.contains("Case summary");
    });

    it("can cancel changing the status", () => {
        cy.addMock("/lpa-api/v1/reference-data/caseChangeReason", "GET", {
            status: 200,
            body: [
                {
                    handle: "LPA_DOES_NOT_WORK",
                    label: "The LPA does not work and cannot be changed",
                    parentSources: ["cannot-register"],
                },
            ],
        });
        cy.contains("Case actions").click();
        cy.contains("Change case status").click();

        cy.url().should("include", "/change-case-status?uid=M-1111-1111-1111");
        cy.contains("Change case status");
        cy.get(".govuk-button-group").contains("Cancel").click();

        cy.url().should("include", "/lpa/M-1111-1111-1111");
        cy.contains("Case summary");
    });

    it("can cancel creating a warning", () => {
        cy.contains("Case actions").click();
        cy.contains("Create a warning").click();

        cy.url().should("include", "/create-warning?id=1111");
        cy.contains("Create Warning");
        cy.contains("Cancel").click();

        cy.url().should("include", "/lpa/M-1111-1111-1111");
        cy.contains("Case summary");
    });

    it("creates a warning via case actions", () => {
        cy.contains(".govuk-button", "Case actions").click();
        cy.contains("Create a warning").click();
        cy.url().should("include", "/create-warning?id=1111");
        cy.get("#f-warningType").select("Complaint Received");
        cy.get("#f-warningText").type("Be warned!");
        cy.get("button[type=submit]").click();

        cy.get(".moj-alert").should("exist");
        cy.get(".moj-alert").contains("Warning created");
        cy.get("h1").contains("Steven Munnell");
        cy.location("pathname").should("eq", "/lpa/M-1111-1111-1111");
    });
});