describe("View a digital LPA", () => {
  beforeEach(() => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");
  });

  it("shows case information", () => {
    cy.contains("M-DIGI-LPA3-3333");
    cy.get("h1").contains("Agnes Hartley");
    cy.get(".govuk-tag.app-tag--draft").contains("Draft");
  });

  it("shows payment information", () => {
    cy.contains("M-DIGI-LPA3-3333");
    cy.get("h1").contains("Agnes Hartley");

    cy.contains("Fees").click();
    cy.contains("£41.00 expected");
  });

  it("shows document information", () => {
    cy.contains("M-DIGI-LPA3-3333");
    cy.get("h1").contains("Agnes Hartley");
    cy.contains("Documents").click();

    cy.contains("Mr Test Person - Blank Template");
    cy.contains("[OUT]");
    cy.contains("24 August 2023");
    cy.contains("EP-BB");

    cy.contains("John Doe - Donor deceased: Case Withdrawn");
    cy.contains("[OUT]");
    cy.contains("15 May 2023");
    cy.contains("DD-4");
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
    cy.get("[data-role=warnings-list] [data-role=warning-detail]").should(
      (elts) => {
        expect(elts).to.have.length(3);

        // check donor deceased is at the top, date is properly-formatted,
        // and applies to text for 3+ cases is correct
        expect(elts[0]).to.contain("05 Jan 2022");
        expect(elts[0]).to.contain("Donor Deceased");
        expect(elts[0]).to.contain(
          "this case, M-DIGI-LPA3-5555 and M-DIGI-LPA3-6666",
        );

        // check sorting has worked properly and case applies text is correct for 2 cases
        expect(elts[1]).to.contain("12 Dec 2023");
        expect(elts[1]).to.contain("Complaint Received");
        expect(elts[1]).to.contain("this case and M-DIGI-LPA3-5555");

        // check case applies text is correct for 1 case
        expect(elts[2]).to.contain("24 Aug 2022");
        expect(elts[2]).to.contain("Court application in progress");
        expect(elts[2]).to.contain("this case");
      },
    );
  });

  it("creates a task via case actions", () => {
    cy.get("select#case-actions").select("Create a task");
    cy.url().should("include", "/create-task?id=333");
    cy.contains("M-DIGI-LPA3-3333");
    cy.get("#f-taskType").select("Check Application");
    cy.get("#f-name").type("Do this task");
    cy.get("#f-description").type("This task, do");
    cy.contains(".govuk-radios__item", "Team").find("input").check();
    cy.get("#f-assigneeTeam").select("Cool Team");
    cy.get("#f-dueDate").type("2035-01-01");
    cy.get("button[type=submit]").click();

    cy.get(".moj-banner").should("exist");
    cy.get(".moj-banner").contains("Task created");
    cy.get("h1").contains("Agnes Hartley");
    cy.location("pathname").should("eq", "/lpa/M-DIGI-LPA3-3333");
  });

  it("creates a warning via case actions", () => {
    cy.get("select#case-actions").select("Create a warning");
    cy.url().should("include", "/create-warning?id=33");
    cy.get("#f-warningType").select("Complaint Received");
    cy.get("#f-warningText").type("Be warned!");
    cy.get("button[type=submit]").click();

    cy.get(".moj-banner").should("exist");
    cy.get(".moj-banner").contains("Warning created");
    cy.get("h1").contains("Agnes Hartley");
    cy.location("pathname").should("eq", "/lpa/M-DIGI-LPA3-3333");
  });
});
