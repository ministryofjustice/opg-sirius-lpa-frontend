import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Add objections form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0008", "GET", {
      status: 200,
      body: {
        uId: "M-0000-0000-0008",
        "opg.poas.sirius": {
          id: 8,
          uId: "M-0000-0000-0008",
          status: "Draft",
          caseSubtype: "personal-welfare",
          createdDate: "31/10/2023",
          investigationCount: 0,
          complaintCount: 0,
          taskCount: 0,
          warningCount: 0,
          donor: {
            id: 8,
          },
          application: {
            donorFirstNames: "James",
            donorLastName: "Rubin",
            donorDob: "22/02/1990",
            donorEmail: "jrubin@mail.example",
            donorPhone: "073656249524",
            donorAddress: {
              addressLine1: "Apartment 3",
              country: "GB",
              postcode: "B15 3AA",
              town: "Birmingham",
            },
            correspondentFirstNames: "Kendrick",
            correspondentLastName: "Lamar",
            correspondentAddress: {
              addressLine1: "Flat 3",
              country: "GB",
              postcode: "SW1 1AA",
              town: "London",
            },
          },
          linkedDigitalLpas: [
            {
              uId: "M-0000-0000-0009",
              caseSubtype: "property-and-affairs",
              status: "Draft",
              createdDate: "01/11/2023",
            },
            {
              uId: "M-0000-0000-0007",
              caseSubtype: "personal-welfare",
              status: "Registered",
              createdDate: "02/11/2023",
            },
            {
              uId: "M-0000-0000-0006",
              caseSubtype: "property-and-affairs",
              status: "In progress",
              createdDate: "01/11/2023",
            },
            {
              uId: "M-0000-0000-0005",
              caseSubtype: "personal-welfare",
              status: "Statutory waiting period",
              createdDate: "02/11/2023",
            },
            {
              uId: "M-0000-0000-0010",
              caseSubtype: "property-and-affairs",
              status: "Do not register",
              createdDate: "01/11/2023",
            },
            {
              uId: "M-0000-0000-0011",
              caseSubtype: "personal-welfare",
              status: "Expired",
              createdDate: "02/11/2023",
            },
            {
              uId: "M-0000-0000-0012",
              caseSubtype: "property-and-affairs",
              status: "Cannot register",
              createdDate: "01/11/2023",
            },
            {
              uId: "M-0000-0000-0013",
              caseSubtype: "personal-welfare",
              status: "Cancelled",
              createdDate: "02/11/2023",
            },
            {
              uId: "M-0000-0000-0014",
              caseSubtype: "property-and-affairs",
              status: "De-registered",
              createdDate: "01/11/2023",
            },
          ],
        },
        "opg.poas.lpastore": {
          donor: {
            uid: "5ff557dd-1e27-4426-9681-ed6e90c2c08d",
            firstNames: "James",
            lastName: "Rubin",
            otherNamesKnownBy: "Somebody",
            dateOfBirth: "1990-02-22",
            address: {
              line1: "Apartment 3",
              town: "Birmingham",
              country: "GB",
              postcode: "B15 3AA",
            },
            contactLanguagePreference: "en",
            email: "jrubin@mail.example",
          },
          attorneys: [
            {
              firstNames: "Esther",
              lastName: "Greenwood",
              status: "active",
            },
          ],
          certificateProvider: {
            uid: "e4d5e24e-2a8d-434e-b815-9898620acc71",
            firstNames: "Timothy",
            lastNames: "Turner",
            signedAt: "2022-12-18T11:46:24Z",
          },
          signedAt: "2024-10-18T11:46:24Z",
          lpaType: "pw",
          channel: "online",
          registrationDate: "2024-11-11",
          peopleToNotify: [],
        },
      },
    });

    cy.addMock(
      `/lpa-api/v1/digital-lpas/M-0000-0000-0008/progress-indicators`,
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-0000-0000-0008",
          progressIndicators: [
            { indicator: "FEES", status: "IN_PROGRESS" },
            { indicator: "DONOR_ID", status: "CANNOT_START" },
            { indicator: "CERTIFICATE_PROVIDER_ID", status: "CANNOT_START" },
            {
              indicator: "CERTIFICATE_PROVIDER_SIGNATURE",
              status: "CANNOT_START",
            },
            { indicator: "ATTORNEY_SIGNATURES", status: "CANNOT_START" },
            { indicator: "PREREGISTRATION_NOTICES", status: "CANNOT_START" },
            { indicator: "REGISTRATION_NOTICES", status: "CANNOT_START" },
          ],
        },
      },
    );

    cy.addMock("/lpa-api/v1/cases/8", "GET", {
      status: 200,
      body: {
        id: 8,
        uId: "M-0000-0000-0008",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 8,
        },
        status: "Draft",
      },
    });

    const mocks = Promise.allSettled([
      cases.warnings.empty("8"),
      cases.tasks.empty("8"),
      digitalLpas.objections.empty("M-0000-0000-0008"),
    ]);

    cy.wrap(mocks);

    cy.visit("/add-objection?uid=M-0000-0000-0008");
  });

  it("can be visited from the LPA details Change link", () => {
    cy.visit("/lpa/M-0000-0000-0008/lpa-details").then(() => {
      cy.contains("Update donor's record").click();
      cy.contains("Add an objection").click();

      cy.url().should("include", "/add-objection?uid=M-0000-0000-0008");
      cy.contains("Add Objection");
    });
  });

  it("objection shows in case summary", () => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-0000-0000-0008/objections", "GET", {
      status: 200,
      body: {
        uid: "M-0000-0000-0008",
        objections: [
          {
            id: 18,
            notes: "test",
            objectionType: "factual",
            receivedDate: "2025-01-01",
            lpaUids: ["M-0000-0000-0008"],
          },
        ],
      },
    });

    cy.visit("/lpa/M-0000-0000-0008/lpa-details").then(() => {
      cy.contains("Objection received");
      cy.contains("Received on 1 January 2025");
      cy.contains("Added to M-0000-0000-0008");
      cy.contains("Record objection outcome");
    });
  });

  it("can go Back to LPA details", () => {
    cy.contains("Back").click();
    cy.url().should("contain", "/lpa/M-0000-0000-0008");
  });

  it("can be cancelled, returning to the LPA details", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-0000-0000-0008");
  });

  it("Filters out cases for adding an objection to", () => {
    cy.contains("PW M-0000-0000-0005");
    cy.contains("PA M-0000-0000-0006");
    cy.contains("PW M-0000-0000-0008");
    cy.contains("PA M-0000-0000-0009");

    cy.contains("M-0000-0000-0007").should("not.exist");
    cy.contains("M-0000-0000-0010").should("not.exist");
    cy.contains("M-0000-0000-0011").should("not.exist");
    cy.contains("M-0000-0000-0012").should("not.exist");
    cy.contains("M-0000-0000-0013").should("not.exist");
    cy.contains("M-0000-0000-0014").should("not.exist");
  });

  it("Can add objection", () => {
    cy.addMock("/lpa-api/v1/objections", "POST", {
      status: 201,
      body: {
        lpaUids: ["M-0000-0000-0008", "M-0000-0000-0009"],
        receivedDate: "20/01/2025",
        objectionType: "factual",
        notes: "Test",
      },
    });

    cy.contains("PW M-0000-0000-0008").click();
    cy.contains("PA M-0000-0000-0009").click();

    cy.get("#f-receivedDate-day").type("20");
    cy.get("#f-receivedDate-month").type("01");
    cy.get("#f-receivedDate-year").type("2025");

    cy.contains("Factual").click();

    cy.get("#f-notes").type("Test");

    cy.get("button[type=submit]").click();

    cy.url().should("contain", "/lpa/M-0000-0000-0008");
  });
});
