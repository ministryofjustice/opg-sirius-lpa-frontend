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
            },
          ],
        },
      },
    );
  });

  it("on a person", () => {
    cy.visit("/view-document/dfef6714-b4fe-44c2-b26e-90dfe3663e95");
    cy.contains("7001-0000-5678");
  });
});
