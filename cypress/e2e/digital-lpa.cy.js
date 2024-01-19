describe("View a digital LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-LPA3-3333",
        "opg.poas.sirius": {
          id: 333,
          uId: "M-DIGI-LPA3-3333",
          status: "Draft",
          caseSubtype: "pfa",
          createdDate: "31/10/2023",
          investigationCount: 2,
          complaintCount: 1,
          taskCount: 2,
          warningCount: 4,
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
          lpaType: "pf",
          registrationDate: "2022-12-18",
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/333/warnings", "GET", {
      status: 200,
      body: [
        {
          id: 44,
          warningType: "Court application in progress",
          warningText: "Court notified",
          dateAdded: "24/08/2022 13:13:13",
          caseItems: [{ uId: "M-DIGI-LPA3-3333", caseSubtype: "pw" }],
        },
        {
          id: 22,
          warningType: "Complaint Received",
          warningText: "Complaint from donor",
          dateAdded: "12/12/2023 12:12:12",
          caseItems: [
            { uId: "M-DIGI-LPA3-3333", caseSubtype: "pw" },
            { uId: "M-DIGI-LPA3-5555", caseSubtype: "hw" },
          ],
        },
        {
          id: 24,
          warningType: "Donor Deceased",
          warningText: "Advised of donor death",
          dateAdded: "05/01/2022 10:10:00",
          caseItems: [
            { uId: "M-DIGI-LPA3-3333", caseSubtype: "pw" },
            { uId: "M-DIGI-LPA3-5555", caseSubtype: "hw" },
            { uId: "M-DIGI-LPA3-6666", caseSubtype: "pw" },
          ],
        },
      ],
    });

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
    cy.get(".app-caseworker-summary > div:nth-child(2) li").should((elts) => {
      expect(elts).to.have.length(3);

      // check donor deceased is at the top, date is properly-formatted,
      // and applies to text for 3+ cases is correct
      expect(elts[0]).to.contain("Donor Deceased");
      expect(elts[0]).to.contain(
        "this case, HW M-DIGI-LPA3-5555 and PW M-DIGI-LPA3-6666",
      );

      // check sorting has worked properly and case applies text is correct for 2 cases
      expect(elts[1]).to.contain("Complaint Received");
      expect(elts[1]).to.contain("this case and HW M-DIGI-LPA3-5555");

      // check case applies text is correct for 1 case
      expect(elts[2]).to.contain("Court application in progress");
      expect(elts[2]).not.to.contain("this case");
    });
  });

  it("creates a task via case actions", () => {
    cy.get(".govuk-button").contains("Case actions").click();
    cy.contains("Create a task").click();
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
    cy.get(".govuk-button").contains("Case actions").click();
    cy.contains("Create a warning").click();
    cy.url().should("include", "/create-warning?id=33");
    cy.get("#f-warningType").select("Complaint Received");
    cy.get("#f-warningText").type("Be warned!");
    cy.get("button[type=submit]").click();

    cy.get(".moj-banner").should("exist");
    cy.get(".moj-banner").contains("Warning created");
    cy.get("h1").contains("Agnes Hartley");
    cy.location("pathname").should("eq", "/lpa/M-DIGI-LPA3-3333");
  });

  it("shows lpa details from lpa store", () => {
    cy.contains("LPA details").click();
    cy.contains("lpaType:pf");
    cy.contains("registrationDate:2022-12-18");
  });
});
