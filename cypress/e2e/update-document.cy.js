describe("Edit a document", () => {
  const visitWithNotifySettings = (sendingLettersViaNotify, usesNotify) => {
    cy.addMock("/lpa-api/v1/templates/lpa", "GET", {
      status: 200,
      body: {
        "LP-A": {
          label: "LPA Registered Notification",
          inserts: {},
          govukNotify: usesNotify,
        },
      },
    });

    cy.addMock("/config", "GET", {
      status: 200,
      body: {
        featureToggles: {
          poasBlankSections: true,
          poasSendingLettersViaNotify: sendingLettersViaNotify,
        },
      },
    });

    cy.visit("/edit-document?id=800&case=lpa");
  };

  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/800", "GET", {
      status: 200,
      body: {
        id: 800,
        caseType: "LPA",
        donor: { id: 33 },
      },
    });
    cy.addMock("/lpa-api/v1/lpas/800/documents?type[]=Draft", "GET", {
      status: 200,
      body: [
        {
          childCount: 0,
          correspondent: {
            firstname: "Consuela",
            id: 189,
            surname: "Aysien",
          },
          createdDate: "15/12/2022 13:41:04",
          direction: "Outgoing",
          filename: "LP-A.pdf",
          friendlyDescription:
            "Dr Consuela Aysien - LPA perfect + reg due date: applicant",
          id: 1,
          mimeType: "application\\/pdf",
          systemType: "LP-A",
          type: "Draft",
          uuid: "dfef6714-b4fe-44c2-b26e-90dfe3663e95",
        },
      ],
    });

    cy.addMock("/lpa-api/v1/templates/lpa", "GET", {
      status: 200,
      body: {
        "LP-A": {
          label: "LPA Registered Notification",
          inserts: {},
          govukNotify: false,
        },
      },
    });

    cy.addMock("/config", "GET", {
      status: 200,
      body: {
        featureToggles: {
          poasBlankSections: true,
          poasSendingLettersViaNotify: false,
        },
      },
    });

    cy.visit("/update-document?id=800&case=lpa");
  });

  it("displays draft to edit", () => {
    cy.contains("Edit draft document");
    cy.contains("button", "Save draft");
    cy.contains("button", "Preview draft");
    cy.contains("button", "Delete draft");
    cy.contains("button", "Publish draft");
    cy.contains("button", "Publish and post via Notify").should("not.exist");
    cy.contains("button", "Publish and Send").should("not.exist");
    cy.contains(
      "When you publish this letter it will automatically be sent via GOV.UK Notify",
    ).should("not.exist");
    cy.contains("button", "Save and exit");
    cy.contains("span", "Insert blank sections of the LPA");
  });

  it("shows the send warning and button when the template uses Notify but the toggle is off", () => {
    visitWithNotifySettings(false, true);

    cy.contains(
      "When you publish this letter it will automatically be sent via GOV.UK Notify",
    ).should("exist");
    cy.contains("button", "Publish and Send").should("exist");
    cy.contains("button", "Publish draft").should("not.exist");
    cy.contains("button", "Publish and post via Notify").should("not.exist");
  });

  it("shows the publish button without a warning when the toggle is on but the template does not use Notify", () => {
    visitWithNotifySettings(true, false);

    cy.contains("button", "Publish draft").should("exist");
    cy.contains("button", "Publish and post via Notify").should("not.exist");
    cy.contains("button", "Publish and Send").should("not.exist");
    cy.contains("When you publish this letter").should("not.exist");
  });

  it("shows separate publish and send buttons when the toggle is on and the template uses Notify", () => {
    visitWithNotifySettings(true, true);

    cy.contains(
      "When you publish this letter it will automatically be sent via GOV.UK Notify",
    ).should("exist");
    cy.contains("button", "Publish draft").should("exist");
    cy.contains("button", "Publish and post via Notify").should("exist");
    cy.contains("button", "Publish and Send").should("not.exist");
  });

  it("does not display the blank sections option on an EPA case", () => {
    cy.addMock("/lpa-api/v1/cases/800", "GET", {
      status: 200,
      body: {
        id: 800,
        caseType: "EPA",
        donor: { id: 33 },
      },
    });
    cy.visit("/update-document?id=800&case=lpa");

    cy.contains("Edit draft document");
    cy.contains("button", "Save draft");
    cy.contains("button", "Preview draft");
    cy.contains("button", "Delete draft");
    cy.contains("button", "Publish draft");
    cy.contains("button", "Save and exit");
    cy.contains("span", "Insert blank sections of the LPA").should("not.exist");
  });

  it("can select a draft to edit", () => {
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    cy.get("#f-document").select(
      "1: 15/12/2022 13:41:04: Consuela Aysien: LP-A",
    );
    cy.contains("button", "Select").click();
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
  });

  it("saves a draft document", () => {
    cy.get("#documentTextEditor").contains("Test content");
    // can only edit the iframe
    const $iframe = cy
      .get("iframe[id=documentTextEditor_ifr]")
      .its("0.contentDocument.body")
      .should("not.be.empty")
      .then(cy.wrap);

    $iframe.clear().type("Edited test content");

    cy.contains("button", "Save draft").click();
    cy.get("#documentTextEditor").contains("Edited test content");
  });

  it("deletes a draft document", () => {
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    cy.get("#documentTextEditor").contains("Test content");

    cy.contains("button", "Delete draft").click();
  });

  it("previews a draft document", () => {
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    cy.get("#documentTextEditor").contains("Test content");

    cy.addMock(
      "/lpa-api/v1/documents/dfef6714-b4fe-44c2-b26e-90dfe3663e95/download",
      "GET",
      {
        status: 200,
        body: "pdf content",
      },
    );

    cy.contains("button", "Preview draft").click();
  });

  it("publishes a draft document", () => {
    cy.contains("1: 15/12/2022 13:41:04: Consuela Aysien: LP-A");
    cy.get("#documentTextEditor").contains("Test content");

    cy.contains("button", "Publish draft").click();
  });
});
