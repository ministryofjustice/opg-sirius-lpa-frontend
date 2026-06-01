describe("Action Panel", () => {
  beforeEach(() => {
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
          { caseType: "LPA", caseSubtype: "hw", id: 78, uId: "7000-5678-5678" },
        ],
      },
    });

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0&limit=999",
      "GET",
      {
        status: 200,
        body: {
          total: 0,
          documents: [],
        },
      },
    );
    cy.visit("/donor/1/documents");
  });

  it("can open and close the action panel", () => {
    cy.get("#actions-tab").click();
    cy.get("#actions-content").should("be.visible");

    cy.get("#actions-tab").click();
    cy.get("#actions-content").should("not.be.visible");
  });

  it("displays the warning button on the action panel", () => {
    cy.get("#actions-tab").click();
    cy.get("#actions-content").should("be.visible");
    cy.get("#actions-content").contains("Create warning");

    cy.get("a#action-panel-button-create-warning").click();
    cy.get(".action-panel__form").should("exist");
    cy.get(".action-panel__form").contains("Create Warning");
  });
});
