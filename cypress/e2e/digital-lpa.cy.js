describe("View a digital LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-LPA3-3333",
        "opg.poas.sirius": {
          id: 333,
          uId: "M-DIGI-LPA3-3333",
          status: "Processing",
          caseSubtype: "property-and-affairs",
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
        "opg.poas.lpastore": {
          attorneys: [
            {
              firstNames: "Esther",
              lastName: "Greenwood",
              status: "active",
            },
            {
              firstNames: "Volo",
              lastName: "McSpolo",
              status: "active",
            },
            {
              firstNames: "Susanna",
              lastName: "Kaysen",
              status: "removed",
            },
            {
              firstNames: "Philomena",
              lastName: "Guinea",
              status: "replacement",
            },
          ],
          lpaType: "pf",
          registrationDate: "2022-12-18",
          peopleToNotify: [],
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/333", "GET", {
      status: 200,
      body: {
        id: 333,
        uId: "M-DIGI-LPA3-3333",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 33,
        },
        status: "Processing",
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

    cy.addMock(
      "/lpa-api/v1/lpas/333/documents?type[-][]=Draft&type[-][]=Preview",
      "GET",
      {
        status: 200,
        body: [
          {
            id: 1,
            uuid: "7327f57d-e3d5-4300-95a8-67b3337c7231",
            friendlyDescription: "Mr Test Person - Blank Template",
            direction: "Outgoing",
            createdDate: "24/08/2023 15:27:16",
            systemType: "EP-BB",
            correspondent: {
              uId: "7000-0000-0013",
              firstname: "Test",
              surname: "Person",
              personType: "Donor",
            },
          },
          {
            id: 2,
            uuid: "40fa2847-27ae-4976-a93a-9f45ec0a4e98",
            friendlyDescription: "Mr John Doe - Reduced fee evidence",
            direction: "Incoming",
            createdDate: "15/05/2023 11:09:28",
            receivedDateTime: "15/05/2023 11:09:28",
            type: "Application Related",
            subType: "Reduced fee request evidence",
            correspondent: {
              uId: "7000-0000-0013",
              firstname: "John",
              surname: "Doe",
              personType: "Correspondent",
            },
          },
        ],
      },
    );

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3334", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-LPA3-3334",
        "opg.poas.sirius": {
          id: 334,
          uId: "M-DIGI-LPA3-3334",
          status: "Draft",
          caseSubtype: "property-and-affairs",
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
          linkedDigitalLpas: [
            {
              uId: "M-DIGI-LPA3-3333",
              caseSubtype: "property-and-affairs",
              status: "Processing",
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
        "opg.poas.lpastore": null,
      },
    });

    cy.addMock("/lpa-api/v1/cases/334", "GET", {
      status: 200,
      body: {
        id: 334,
        uId: "M-DIGI-LPA3-3334",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 33,
        },
        status: "Draft",
      },
    });

    cy.addMock("/lpa-api/v1/cases/334/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      "/lpa-api/v1/cases/334/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.addMock(
      `/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333/progress-indicators`,
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-DIGI-LPA3-3333",
          progressIndicators: [
            { indicator: "FEES", status: "NOT_STARTED" },
            { indicator: "FEES", status: "COMPLETE" },
            { indicator: "FEES", status: "IN_PROGRESS" },
          ],
        },
      },
    );
  });

  it("shows case information", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.get("h1").contains("Agnes Hartley");

    cy.contains("M-DIGI-LPA3-3333");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3333'] .govuk-tag").contains("Processing");

    cy.contains("PW M-DIGI-LPA3-3334");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3334'] .govuk-tag").contains("Draft");

    cy.contains("PW M-DIGI-LPA3-3335");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3335'] .govuk-tag").contains("Registered");
  });

  it("shows payment information", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("M-DIGI-LPA3-3333");
    cy.get("h1").contains("Agnes Hartley");

    cy.contains("Fees").click();
    cy.contains("£41.00 expected");
  });

  it("shows document information", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("M-DIGI-LPA3-3333");
    cy.get("h1").contains("Agnes Hartley");
    cy.contains("Documents").click();

    cy.contains("Mr Test Person - Blank Template");
    cy.contains("[OUT]");
    cy.contains("24 August 2023");
    cy.contains("EP-BB");

    cy.contains("Mr John Doe - Reduced fee evidence");
    cy.contains("[IN]");
    cy.contains("15 May 2023");
    cy.contains("Application Related");
    cy.contains("Reduced fee request evidence");
  });

  it("shows task table", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

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
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.get(".app-caseworker-summary > div:nth-child(2) li").should((elts) => {
      expect(elts).to.have.length(3);

      // check donor deceased is at the top, date is properly-formatted,
      // and applies to text for 3+ cases is correct
      expect(elts[0]).to.contain("Donor Deceased");
      expect(elts[0]).to.contain(
        "this case, PA M-DIGI-LPA3-5555 and PW M-DIGI-LPA3-6666",
      );

      // check sorting has worked properly and case applies text is correct for 2 cases
      expect(elts[1]).to.contain("Complaint Received");
      expect(elts[1]).to.contain("this case and PA M-DIGI-LPA3-5555");

      // check case applies text is correct for 1 case
      expect(elts[2]).to.contain("Court application in progress");
      expect(elts[2]).not.to.contain("this case");
    });
  });

  it("creates a task via case actions", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.get(".govuk-button").contains("Case actions").click();
    cy.contains("Create a task").click();
    cy.url().should("include", "/create-task?id=333");
    cy.contains("M-DIGI-LPA3-3333");
    cy.get("#f-taskType").select("Check Application");
    cy.get("#f-name").type("Do this task");
    cy.get("#f-description").type("This task, do");
    cy.contains("label", "Team").click();
    cy.get("#f-assigneeTeam").select("Cool Team");
    cy.get("#f-dueDate").type("2035-01-01");
    cy.get("button[type=submit]").click();

    cy.get(".moj-banner").should("exist");
    cy.get(".moj-banner").contains("Task created");
    cy.get("h1").contains("Agnes Hartley");
    cy.location("pathname").should("eq", "/lpa/M-DIGI-LPA3-3333");
  });

  it("creates a warning via case actions", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

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

  it("shows lpa details from store when status is Processing", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("LPA details").click();
    cy.contains("Attorneys (2)");
    cy.contains("Replacement attorneys (1)");
    cy.contains("Notified people (0)");
    cy.contains("Correspondent");
  });

  it("shows application details when status is Draft", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3334");

    cy.contains("LPA details").click();
    cy.contains("Application format");
    cy.contains("Paper");
  });
});
