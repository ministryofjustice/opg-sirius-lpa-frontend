describe("Compare documents", () => {
  beforeEach(() => {
    let $documentWithCase = {
      createdDate: "15/12/2022 13:41:04",
      friendlyDescription:
        "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
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
    };

    let $documentWithoutCase = {
      createdDate: "15/12/2022 13:41:04",
      friendlyDescription: "Dr Consuela Aysien - A document not linked to case",
      uuid: "e5b5acd1-c11c-41fe-a921-7fdd07e8f670",
      persons: [
        {
          id: 33,
        },
      ],
    };

    cy.addMock(
      "/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95",
      "GET",
      {
        status: 200,
        body: $documentWithCase,
      },
    );

    cy.addMock(
      "/lpa-api/v1/documents/e5b5acd1-c11c-41fe-a921-7fdd07e8f670",
      "GET",
      {
        status: 200,
        body: $documentWithoutCase,
      },
    );

    cy.addMock(
      "/lpa-api/v1/persons/33/documents?filter=draft:0,preview:0,case:34&limit=999",
      "GET",
      {
        status: 200,
        body: {
          total: 2,
          documents: [$documentWithCase, $documentWithoutCase],
        },
      },
    );

    cy.addMock(
      "/lpa-api/v1/persons/33/documents?filter=draft:0,preview:0&limit=999",
      "GET",
      {
        status: 200,
        body: {
          total: 2,
          documents: [$documentWithCase, $documentWithoutCase],
        },
      },
    );
  });

  it("shows document alongside document list when first selecting compare", () => {
    cy.visit("/compare/33/34?pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95");
    cy.contains("7001-0000-5678");
    cy.get("#main-content > :nth-child(1) > .govuk-button")
      .contains("Back to list")
      .should("have.attr", "href")
      .and("include", "/compare/33/34");
    cy.get(".govuk-table__head .govuk-table__row").within(() => {
      cy.get("th").eq(0).should("contain", "Select");
      cy.get("th").eq(1).should("contain", "Name");
      cy.get("th").eq(2).should("contain", "Date created");
      cy.get("th").eq(3).should("contain", "Document Type");
    });
  });

  it("shows document not linked to a case alongside document list when first selecting compare", () => {
    //need to revisit this
    cy.visit("/compare/33/34?pane1=e5b5acd1-c11c-41fe-a921-7fdd07e8f670");
    cy.get("#main-content > :nth-child(1) > .govuk-button")
      .contains("Back to list")
      .should("have.attr", "href")
      .and("include", "/compare/33/34");
    cy.get(".govuk-table__head .govuk-table__row").within(() => {
      cy.get("th").eq(0).should("contain", "Select");
      cy.get("th").eq(1).should("contain", "Name");
      cy.get("th").eq(2).should("contain", "Date created");
      cy.get("th").eq(3).should("contain", "Document Type");
    });
  });

  it("shows two documents alongside each other", () => {
    cy.visit(
      "compare/33/34?pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95&pane2=e5b5acd1-c11c-41fe-a921-7fdd07e8f670",
    );
    cy.contains("7001-0000-5678");
    cy.contains("A document not linked to case");
    cy.get("#main-content > :nth-child(1) > .govuk-button")
      .contains("Back to list")
      .should("have.attr", "href")
      .and(
        "include",
        "/compare/33/34?pane2=e5b5acd1-c11c-41fe-a921-7fdd07e8f670",
      );
    cy.get("#main-content > :nth-child(2) > .govuk-button")
      .contains("Back to list")
      .should("have.attr", "href")
      .and(
        "include",
        "/compare/33/34?pane1=dfef6714-b4fe-44c2-b26e-90dfe3663e95",
      );
  });

  it("shows one document in view on the right", () => {
    cy.visit("compare/33/34?pane2=e5b5acd1-c11c-41fe-a921-7fdd07e8f670");
    cy.contains("7001-0000-5678");
    cy.contains("A document not linked to case");
    cy.get("#main-content > :nth-child(2) > .govuk-button")
      .contains("Back to list")
      .should("have.attr", "href")
      .and("include", "/compare/33/34");
  });

  it("shows two lists", () => {
    cy.visit("compare/33/34");
    cy.contains("7001-0000-5678");
    cy.get(".govuk-grid-column-one-half")
      .eq(0)
      .find("button")
      .contains("Download")
      .should("exist");

    cy.get(".govuk-grid-column-one-half")
      .eq(0)
      .find("table.govuk-table")
      .should("exist");

    cy.get(".govuk-grid-column-one-half")
      .eq(1)
      .find("button")
      .contains("Download")
      .should("exist");

    cy.get(".govuk-grid-column-one-half")
      .eq(1)
      .find("table.govuk-table")
      .should("exist");
  });
});
