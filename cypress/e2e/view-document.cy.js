describe("View documents", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/persons/1", "GET", {
      status: 200,
      body: {},
    });

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
            caseSubtype: "pfa",
            id: 990,
            uId: "7001-0000-5678",
          },
        ],
      },
    });

    cy.addMock(
      "/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95",
      "GET",
      {
        status: 200,
        body: {
          childCount: 0,
          content: "Test content",
          correspondent: {
            id: 189,
          },
          createdDate: "15/12/2022 13:41:04",
          direction: "Outgoing",
          filename: "LP-A.pdf",
          friendlyDescription:
            "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
          id: 1,
          mimeType: "application/pdf",
          systemType: "LP-A",
          type: "Save",
          uuid: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
          caseItems: [
            {
              uId: "7001-0000-5678",
              caseSubtype: "pfa",
              caseType: "EPA",
              donor: {
                id: 33,
              },
            },
          ],
        },
      },
    );

    cy.addMock("/lpa-api/v1/cases/1", "GET", {
      status: 200,
      body: {
        id: 1,
        caseType: "EPA",
        caseSubtype: "pfa",
        uId: "7001-0000-5678",
      },
    });

    cy.addMock("/lpa-api/v1/permissions", "GET", {
      status: 200,
      body: {
        "v1-persons": {
          permissions: ["GET"],
        },
        "v1-persons-cases": {
          permissions: ["GET"],
        },
      },
    });

    cy.addMock("/lpa-api/v1/epas/1/draft-count", "GET", {
      status: 200,
      body: {
        draftCount: 0,
      },
    });

    cy.addMock(
      "/lpa-api/v1/cases/1/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.addMock("/lpa-api/v1/persons/1/references", "GET", {
      status: 200,
      body: [
        {
          referenceId: 123,
        },
      ],
    });

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0&limit=999",
      "GET",
      {
        status: 200,
        body: {
          documents: [
            {
              uuid: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
              filename: "LP-A.pdf",
              direction: "Outgoing",
              createdDate: "15/12/2022 13:41:04",
            },
          ],
        },
      },
    );
  });

  it("views a document as a user with system admin role", () => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "System Admin"],
      },
    });
    cy.visit("/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95/1?case=1");
    cy.contains("7001-0000-5678");
    cy.get('a:contains("Back to list")')
      .should("exist")
      .should("have.attr", "href")
      .and("include", "/donor/33/documents?uid[]=7001-0000-5678");
    cy.get('a:contains("Download")')
      .should("exist")
      .should("have.attr", "href")
      .and(
        "include",
        "/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95/download",
      );
    cy.get(".govuk-button--warning").contains("Delete");
  });

  it("views a document and all the pdf actions work", () => {
    cy.mockDocumentFile("dfef6714-b4fe-44c2-b26e-90dfe3663e95");

    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "System Admin"],
      },
    });

    cy.visit("/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95/1?case=1");
    cy.contains("7001-0000-5678");
    cy.get('a:contains("Back to list")')
      .should("exist")
      .should("have.attr", "href")
      .and("include", "/donor/33/documents?uid[]=7001-0000-5678");
    cy.get('a:contains("Download")')
      .should("exist")
      .should("have.attr", "href")
      .and(
        "include",
        "/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95/download",
      );

    cy.get('[data-action="toggle-thumbnails"]').should("exist");
    cy.get('[data-action="prev"]').should("exist");
    cy.get('[data-action="next"]').should("exist");
    cy.get(".pdf-viewer-page-input").should("exist");
    cy.get(".pdf-viewer-zoom-input").should("exist");
  });

  it("views a document as a user without system admin role", () => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User"],
      },
    });
    cy.visit("/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95/1?case=1");
    cy.contains("7001-0000-5678");
    cy.get(".govuk-button--warning").should("not.exist");
    cy.get('a:contains("Back to list")')
      .should("exist")
      .should("have.attr", "href")
      .and("include", "/donor/33/documents?uid[]=7001-0000-5678");
    cy.get('a:contains("Download")')
      .should("exist")
      .should("have.attr", "href")
      .and(
        "include",
        "/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95/download",
      );
  });

  it("views a document linked to a person not a case", () => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User"],
      },
    });

    cy.addMock(
      "/lpa-api/v1/documents/e5b5acd1-c11c-41fe-a921-7fdd07e8f670",
      "GET",
      {
        status: 200,
        body: {
          createdDate: "15/12/2022 13:41:04",
          friendlyDescription:
            "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
          uuid: "e5b5acd1-c11c-41fe-a921-7fdd07e8f670",
          persons: [
            {
              id: 33,
            },
          ],
        },
      },
    );

    cy.visit("/view-document/e5b5acd1-c11c-41fe-a921-7fdd07e8f670/1?case=1");
    cy.get(".govuk-button--warning").should("not.exist");
    cy.get('a:contains("Back to list")')
      .should("exist")
      .should("have.attr", "href")
      .and("include", "/donor/33/documents");
    cy.get('a:contains("Download")')
      .should("exist")
      .should("have.attr", "href")
      .and(
        "include",
        "/lpa-api/v1/documents/e5b5acd1-c11c-41fe-a921-7fdd07e8f670/download",
      );
  });
});
