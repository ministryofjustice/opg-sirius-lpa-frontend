describe("View documents", () => {
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
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0,case:34&limit=999",
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
              id: 332,
              uuid: "5b4f0ad3-1e4a-4a55-b4a7-3f8e3d2bc3b9",
              type: "LPA",
              friendlyDescription: "LP1F - Finance Instrument",
              createdDate: "29/05/2022 10:07:38",
              direction: "Incoming",
              filename: "LP1F.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-1234-1234",
                  caseSubtype: "pfa",
                  caseType: "LPA",
                },
              ],
              persons: [],
              subtype: "pfa",
            },
          ],
        },
      },
    );

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0,case:34,case:78&limit=999",
      "GET",
      {
        status: 200,
        body: {
          limit: 999,
          metadata: {
            doctype: {
              correspondence: 3,
              order: 0,
              report: 0,
              visit: 0,
              finance: 0,
              other: 0,
            },
            direction: {
              Incoming: 2,
              Outgoing: 1,
              Internal: 0,
            },
          },
          pages: {
            current: 1,
            total: 1,
          },
          total: 4,
          documents: [
            {
              id: 332,
              uuid: "5b4f0ad3-1e4a-4a55-b4a7-3f8e3d2bc3b9",
              type: "LPA",
              friendlyDescription: "LP1F - Finance Instrument",
              createdDate: "29/05/2022 10:07:38",
              direction: "Incoming",
              filename: "LP1F.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-1234-1234",
                  caseSubtype: "pfa",
                  caseType: "LPA",
                },
              ],
              persons: [],
              subtype: "pfa",
            },
            {
              id: 443,
              uuid: "c8e3a1df-7b9b-4d45-94d9-2b8fc0d9e0fd",
              type: "LPA",
              friendlyDescription: "LP1H - Health Instrument",
              createdDate: "01/06/2022 15:39:01",
              direction: "Incoming",
              filename: "LP1H.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-5678-5678",
                  caseSubtype: "hw",
                  caseType: "LPA",
                },
              ],
              subtype: "hw",
            },
            {
              id: 639,
              uuid: "31e6f4c2-5f8b-47c3-bc98-64b47c938e52",
              type: "Save",
              friendlyDescription: "Letter",
              createdDate: "25/07/2022 14:17:13",
              direction: "Outgoing",
              filename: "LP-NA-3A.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-5678-5678",
                  caseSubtype: "hw",
                  caseType: "LPA",
                },
              ],
              persons: [],
              systemType: "LP-NA-3A",
            },
            {
              id: 640,
              uuid: "42e6f4c2-5f8b-47c3-bc98-64b47c938e52",
              type: "Save",
              friendlyDescription: "Letter",
              createdDate: "26/08/2022 08:11:27",
              direction: "Outgoing",
              notifyStatus: "posted",
              filename: "LP-WHAT.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-5678-5678",
                  caseSubtype: "hw",
                  caseType: "LPA",
                },
              ],
              persons: [],
              systemType: "LP-WHAT",
            },
          ],
        },
      },
    );

    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0&limit=999",
      "GET",
      {
        status: 200,
        body: {
          limit: 999,
          metadata: {
            doctype: {
              correspondence: 5,
              order: 0,
              report: 0,
              visit: 0,
              finance: 0,
              other: 0,
            },
            direction: {
              Incoming: 4,
              Outgoing: 2,
              Internal: 0,
            },
          },
          pages: {
            current: 1,
            total: 1,
          },
          total: 6,
          documents: [
            {
              id: 332,
              uuid: "5b4f0ad3-1e4a-4a55-b4a7-3f8e3d2bc3b9",
              type: "LPA",
              friendlyDescription: "LP1F - Finance Instrument",
              createdDate: "29/05/2022 10:07:38",
              direction: "Incoming",
              filename: "LP1F.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-1234-1234",
                  caseSubtype: "pfa",
                  caseType: "LPA",
                },
              ],
              persons: [],
              subtype: "pfa",
            },
            {
              id: 443,
              uuid: "c8e3a1df-7b9b-4d45-94d9-2b8fc0d9e0fd",
              type: "LPA",
              friendlyDescription: "LP1H - Health Instrument",
              createdDate: "01/06/2022 15:39:01",
              direction: "Incoming",
              filename: "LP1H.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-5678-5678",
                  caseSubtype: "hw",
                  caseType: "LPA",
                },
              ],
              subtype: "hw",
            },
            {
              id: 639,
              uuid: "31e6f4c2-5f8b-47c3-bc98-64b47c938e52",
              type: "Save",
              friendlyDescription: "Letter",
              createdDate: "25/07/2022 14:17:13",
              direction: "Outgoing",
              filename: "LP-NA-3A.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-5678-5678",
                  caseSubtype: "hw",
                  caseType: "LPA",
                },
              ],
              persons: [],
              systemType: "LP-NA-3A",
            },
            {
              id: 640,
              uuid: "42e6f4c2-5f8b-47c3-bc98-64b47c938e52",
              type: "Save",
              friendlyDescription: "Letter",
              createdDate: "26/08/2022 08:11:27",
              direction: "Outgoing",
              notifyStatus: "posted",
              filename: "LP-WHAT.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-5678-5678",
                  caseSubtype: "hw",
                  caseType: "LPA",
                },
              ],
              persons: [],
              systemType: "LP-WHAT",
            },
            {
              id: 928,
              uuid: "d9e12f73-3ab2-4d24-9a63-6b0b3e49b1c5",
              type: "Application Related",
              friendlyDescription: "EPA.pdf",
              createdDate: "08/01/2025 10:36:41",
              direction: "Incoming",
              filename: "EPA.pdf",
              mimeType: "application/pdf",
              note: {
                description: "Manual Upload",
              },
              caseItems: [
                {
                  uId: "7001-0000-5678",
                  caseSubtype: "pfa",
                  caseType: "EPA",
                },
              ],
              persons: [],
              subtype: "pfa",
            },
            {
              id: 11,
              uuid: "b829b617-8831-4b6b-864b-327a5d84b925",
              type: "Application Related",
              friendlyDescription: "email.msg",
              createdDate: "04/12/2025 14:56:38",
              direction: "Incoming",
              filename: "6931a1268fac6_receiptdateForm.png",
              mimeType: "email",
              note: {
                description: "test",
              },
              caseItems: [],
              persons: [
                {
                  uId: "7000-0000-0000",
                },
              ],
            },
          ],
        },
      },
    );
  });

  it("on a person", () => {
    cy.visit("/donor/1/documents");
    cy.contains("Documents (6)");
    cy.contains("Viewing 3 POAs:");

    cy.contains("Name");
    cy.contains("Case number");
    cy.contains("Date created");
    cy.contains("Document Type");
  });

  it("on a single case", () => {
    cy.visit("/donor/1/documents?uid[]=7000-1234-1234");

    cy.contains("Documents (1)");
    cy.contains("Viewing 1 POA:");

    cy.contains("Name");
    cy.contains("Case number").should("not.exist");
    cy.contains("Date created");
    cy.contains("Document Type");
  });

  it("on multiple cases", () => {
    cy.visit("/donor/1/documents?uid[]=7000-1234-1234&uid[]=7000-5678-5678");

    cy.contains("Documents (4)");
    cy.contains("Viewing 2 POAs:");

    cy.contains("Name");
    cy.contains("Case number");
    cy.contains("Date created");
    cy.contains("Document Type");
  });

  it("sorts by friendly description regardless of direction label", () => {
    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0,case:78&limit=999",
      "GET",
      {
        status: 200,
        body: {
          limit: 999,
          metadata: {
            doctype: {
              correspondence: 2,
              order: 0,
              report: 0,
              visit: 0,
              finance: 0,
              other: 0,
            },
            direction: {
              Incoming: 1,
              Outgoing: 1,
              Internal: 0,
            },
          },
          pages: {
            current: 1,
            total: 1,
          },
          total: 2,
          documents: [
            {
              id: 201,
              uuid: "3bb33aa0-1234-4891-8739-1aa0fbcdee01",
              type: "Save",
              friendlyDescription: "A Document",
              createdDate: "02/09/2024 09:00:00",
              direction: "Outgoing",
              notifyStatus: "posted",
              filename: "a.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-5678-5678",
                  caseSubtype: "hw",
                  caseType: "LPA",
                },
              ],
              persons: [],
              systemType: "A",
            },
            {
              id: 202,
              uuid: "7cc44bb0-5678-4891-8739-5bb1facdef02",
              type: "Correspondence",
              friendlyDescription: "Z Document",
              createdDate: "03/09/2024 10:00:00",
              direction: "Incoming",
              filename: "z.pdf",
              mimeType: "application/pdf",
              caseItems: [
                {
                  uId: "7000-5678-5678",
                  caseSubtype: "hw",
                  caseType: "LPA",
                },
              ],
              persons: [],
              subtype: "hw",
            },
          ],
        },
      },
    );

    cy.visit("/donor/1/documents?uid[]=7000-5678-5678");

    cy.contains("Documents (2)");

    cy.contains("th", "Name").find("button").click();
    cy.contains("th", "Name").should("have.attr", "aria-sort", "ascending");

    cy.get("tbody tr td:nth-child(2) a").then(($links) => {
      const names = [...$links].map((link) => link.textContent.trim());
      expect(names).to.deep.equal(["A Document", "Z Document"]);
    });

    cy.contains("th", "Name").find("button").click();
    cy.contains("th", "Name").should("have.attr", "aria-sort", "descending");

    cy.get("tbody tr td:nth-child(2) a").then(($links) => {
      const names = [...$links].map((link) => link.textContent.trim());
      expect(names).to.deep.equal(["Z Document", "A Document"]);
    });
  });

  it("does not show selection checkbox for non-PDF documents", () => {
    cy.visit("/donor/1/documents");

    cy.contains("a", "email.msg")
      .closest("tr")
      .find('input[name="document"]')
      .should("not.exist");

    cy.get('input[name="document"]').should("have.length", 5);
  });

  it("an EPA document", () => {
    cy.addMock(
      "/lpa-api/v1/persons/1/documents?filter=draft:0,preview:0,case:990&limit=999",
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
              id: 928,
              uuid: "d9e12f73-3ab2-4d24-9a63-6b0b3e49b1c5",
              type: "Application Related",
              friendlyDescription: "EPA.pdf",
              createdDate: "08/01/2025 10:36:41",
              direction: "Incoming",
              filename: "EPA.pdf",
              mimeType: "application/pdf",
              note: {
                description: "Manual Upload",
              },
              caseItems: [
                {
                  uId: "7001-0000-5678",
                  caseSubtype: "pfa",
                  caseType: "EPA",
                },
              ],
              persons: [],
              subtype: "pfa",
            },
          ],
        },
      },
    );

    cy.visit("/donor/1/documents?uid[]=7001-0000-5678");

    cy.contains("Documents (1)");
    cy.contains("Viewing 1 POA:");
    cy.contains("EPA");
    cy.contains("Name");
    cy.contains("Case number").should("not.exist");
    cy.contains("Date created");
    cy.contains("Document Type");
  });

  it("downloads selected document", () => {
    cy.addMock(
      "/lpa-api/v1/documents/download-multiple?id%5B%5D=5b4f0ad3-1e4a-4a55-b4a7-3f8e3d2bc3b9",
      "GET",
      {
        status: 200,
        headers: {
          "Content-Type": "application/pdf",
          "Content-Disposition": "attachment; filename=document.pdf",
        },
        body: "document",
      },
    );

    cy.intercept("POST", "/donor/1/documents*").as("download");

    cy.visit("/donor/1/documents");

    cy.get('input[name="document"]').first().check({ force: true });
    cy.contains("button", "Download").click();

    cy.wait("@download").then(({ request, response }) => {
      const requestContent = [request.url, request.body || ""].join("&");
      expect(decodeURIComponent(requestContent)).to.contain(
        "document=5b4f0ad3-1e4a-4a55-b4a7-3f8e3d2bc3b9",
      );
      expect(response.statusCode).to.eq(200);
      expect(response.headers["content-type"]).to.contain("application/pdf");
    });
  });

  it("downloads multiple selected documents", () => {
    cy.addMock(
      "/lpa-api/v1/documents/download-multiple?id%5B%5D=5b4f0ad3-1e4a-4a55-b4a7-3f8e3d2bc3b9&id%5B%5D=c8e3a1df-7b9b-4d45-94d9-2b8fc0d9e0fd",
      "GET",
      {
        status: 200,
        headers: {
          "Content-Type": "application/pdf",
          "Content-Disposition": "attachment; filename=document.pdf",
        },
        body: "document",
      },
    );

    cy.intercept("POST", "/donor/1/documents*").as("download");

    cy.visit("/donor/1/documents");

    cy.get('input[name="document"]').eq(0).check({ force: true });
    cy.get('input[name="document"]').eq(1).check({ force: true });
    cy.contains("button", "Download").click();

    cy.wait("@download").then(({ request, response }) => {
      const requestContent = [request.url, request.body || ""].join("&");
      expect(decodeURIComponent(requestContent)).to.contain(
        "document=5b4f0ad3-1e4a-4a55-b4a7-3f8e3d2bc3b9",
      );
      expect(decodeURIComponent(requestContent)).to.contain(
        "document=c8e3a1df-7b9b-4d45-94d9-2b8fc0d9e0fd",
      );
      expect(response.statusCode).to.eq(200);
    });
  });

  it("shows a validation error when nothing is selected", () => {
    cy.intercept("POST", "/donor/1/documents*").as("submit");

    cy.visit("/donor/1/documents");

    cy.contains("button", "Download").click();

    cy.wait("@submit");

    cy.contains("No document selected");
    cy.contains("Select one or more documents and try again.");
  });

  it("dismisses the validation error", () => {
    cy.intercept("POST", "/donor/1/documents*").as("submit");

    cy.visit("/donor/1/documents");

    cy.contains("button", "Download").click();
    cy.wait("@submit");
    cy.contains("No document selected");

    cy.contains("button", "Dismiss").click();
    cy.wait("@submit");

    cy.contains("No document selected").should("not.exist");
  });

  it("shows success notification when document is deleted", () => {
    cy.visit(
      "/donor/1/documents?success=true&documentFriendlyName=LP1F%20-%20Finance%20Instrument&documentCreatedTime=29/05/2022%2010:07:38",
    );

    cy.get("#govuk-notification-banner-title").contains("File deleted");
    cy.contains("29/05/2022 LP1F - Finance Instrument");
  });
});
