describe("View and edit anomalies for a digital LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-QQQQ-1111", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-QQQQ-1111",
        "opg.poas.sirius": {
          id: 111,
          uId: "M-DIGI-QQQQ-1111",
          status: "Processing",
          caseSubtype: "property-and-affairs",
        },
        "opg.poas.lpastore": {
          channel: "online",
          attorneys: [
            {
              uid: "attorney-1-uid",
              status: "active",
            },
            {
              uid: "replacement-attorney-1-uid",
              status: "replacement",
            },
          ],
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/111/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      "/lpa-api/v1/cases/111/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-QQQQ-1111/anomalies", "GET", {
      status: 200,
      body: {
        uid: "M-DIGI-QQQQ-1111",
        anomalies: [
          {
            id: 123,
            status: "detected",
            fieldName: "firstNames",
            ruleType: "empty",
            fieldOwnerUid: "attorney-1-uid",
          },
          {
            id: 124,
            status: "detected",
            fieldName: "lastName",
            ruleType: "empty",
            fieldOwnerUid: "attorney-1-uid",
          },
          {
            id: 125,
            status: "detected",
            fieldName: "firstNames",
            ruleType: "empty",
            fieldOwnerUid: "replacement-attorney-1-uid",
          },
          {
            id: 126,
            status: "detected",
            fieldName: "lastName",
            ruleType: "empty",
            fieldOwnerUid: "replacement-attorney-1-uid",
          },
          {
            id: 127,
            status: "detected",
            fieldName: "howAttorneysMakeDecisions",
            ruleType: "empty",
            fieldOwnerUid: "",
          },
          {
            id: 128,
            status: "detected",
            fieldName: "whenTheLpaCanBeUsed",
            ruleType: "empty",
            fieldOwnerUid: "",
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-SSSS-3333", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-SSSS-3333",
        "opg.poas.sirius": {
          id: 222,
          uId: "M-DIGI-SSSS-3333",
          status: "Processing",
          caseSubtype: "personal-welfare",
        },
        "opg.poas.lpastore": {
          channel: "online",
          attorneys: [
            {
              uid: "attorney-1-uid",
              status: "active",
            },
          ],
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/222/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      "/lpa-api/v1/cases/222/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-SSSS-3333/anomalies", "GET", {
      status: 200,
      body: {
        uid: "M-DIGI-SSSS-3333",
        anomalies: [
          {
            id: 220,
            status: "detected",
            fieldName: "lifeSustainingTreatmentOption",
            ruleType: "empty",
            fieldOwnerUid: "",
          },
        ],
      },
    });
  });

  it("shows anomalies for pfa LPA", () => {
    cy.visit("/lpa/M-DIGI-QQQQ-1111/lpa-details");
    cy.contains("Some LPA details have been identified for review.");
    cy.contains("For review");
    cy.contains("Review attorney's first names");
    cy.contains("Review attorney's last name");
    cy.contains("Review replacement attorney's first names");
    cy.contains("Review replacement attorney's last name");
    cy.contains("Review how attorney's can make decisions");
    cy.contains("Review when the LPA can be used");
  });

  it("shows anomalies for pa LPA", () => {
    cy.visit("/lpa/M-DIGI-SSSS-3333/lpa-details");
    cy.contains("Some LPA details have been identified for review.");
    cy.contains("For review");
    cy.contains("Review life sustaining treatment");
  });
});
