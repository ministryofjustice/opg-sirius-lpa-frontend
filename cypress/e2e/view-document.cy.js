describe("View documents", () => {
  beforeEach(() => {
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
  });

  it("views a document as a user with system admin role", () => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User", "System Admin"],
      },
    });
    cy.visit("/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95");
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

  it("views a document as a user without system admin role", () => {
    cy.addMock("/lpa-api/v1/users/current", "GET", {
      status: 200,
      body: {
        roles: ["OPG User"],
      },
    });
    cy.visit("/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95");
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
    cy.visit("/view-document/e5b5acd1-c11c-41fe-a921-7fdd07e8f670");
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
