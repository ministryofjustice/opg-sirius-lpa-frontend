describe("Compare documents", () => {
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
              id: 34,
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

    cy.addMock(
      "/lpa-api/v1/persons/33/documents?filter=draft:0,preview:0,case:34&limit=999",
      "GET",
      {
        status: 200,
        body: {
          limit: 999,
          metadata: {
            doctype: {
              correspondence: 1,
              order: 0,
              report: 0,
              visit: 0,
              finance: 0,
              other: 0,
            },
            direction: {
              Incoming: 1,
              Outgoing: 0,
              Internal: 0,
            },
          },
          pages: {
            current: 1,
            total: 1,
          },
          total: 1,
          documents: [
            {
              id: 1,
              uuid: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
              type: "Save",
              friendlyDescription: "LP-A",
              createdDate: "15/12/2022 13:41:04",
              direction: "Outgoing",
              filename: "LP-A.pdf",
              mimeType: "application/pdf",
              systemType: "LP-A",
              caseItems: [
                {
                  id: 34,
                  uId: "7001-0000-5678",
                  caseSubtype: "pfa",
                  caseType: "EPA",
                },
              ],
              persons: [],
            },
          ],
        },
      },
    );
  });

  it("compares documents for a donor", () => {
    cy.visit(
      "/compare/33/documents?uid[]=dfef6714-b4fe-44c2-b26e-90dfe3663e95",
    );
    cy.contains("7001-0000-5678");
    cy.get("#main-content > :nth-child(1) > .govuk-button").contains(
      "Back to list",
    );
    cy.get(".govuk-table__head .govuk-table__row").within(() => {
      cy.get("th").eq(0).should("contain", "Select");
      cy.get("th").eq(1).should("contain", "Name");
      cy.get("th").eq(2).should("contain", "Date created");
      cy.get("th").eq(3).should("contain", "Document Type");
    });
  });
});
