import * as digitalLpas from "../mocks/digitalLpas";
import * as cases from "../mocks/cases";

describe("Case actions drop down", () => {
    beforeEach(() => {

        const mocks = Promise.allSettled([
            digitalLpas.get("M-1111-1111-1111"),
            cases.warnings.empty("1111"),
            cases.tasks.empty("1111"),
            digitalLpas.objections.empty("M-1111-1111-1111"),
            digitalLpas.progressIndicators.feesInProgress("M-1111-1111-1111"),
        ]);

        cy.wrap(mocks);

        cy.addMock("/lpa-api/v1/digital-lpas/M-1111-1111-1111/anomalies", "GET", {
            status: 200,
            body: [],
        });

        cy.addMock("/lpa-api/v1/cases/1111/tasks", "GET", {
            status: 200,
            body: { tasks: [] },
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

        cy.addMock("/lpa-api/v1/cases/1111/payments", "GET", {
            status: 200,
            body: [
                {
                    amount: 4100,
                    case: {
                        id: 1111,
                    },
                    id: 1234,
                    paymentDate: "12/02/2025",
                    source: "ONLINE",
                    locked: true,
                },
            ],
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

    it("creates a warning via case actions", () => {
        cy.addMock("/lpa-api/v1/warnings", "POST", {
            status: 201,
            body: {
                personId: 1111,
                warningText: "Be warned!",
                warningType: "Complaint Received",
            },
        });
        cy.contains(".govuk-button", "Case actions").click();
        cy.contains("Create a warning").click();
        cy.url().should("include", "/create-warning?id=1111");
        cy.get("#f-warningType").select("Complaint Received");
        cy.get("#f-warningText").type("Be warned!");
        cy.get("button[type=submit]").click();

        cy.get(".moj-alert").should("exist");
        cy.get(".moj-alert").contains("Warning created");
        cy.get("h1").contains("Agnes Hartley");
        cy.location("pathname").should("eq", "/lpa/M-1111-1111-1111");
    });
});