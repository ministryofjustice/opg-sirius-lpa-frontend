import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

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
              appointmentType: "original",
              status: "active",
              signedAt: "2025-09-30T10:25:47.881438922Z",
            },
            {
              uid: "replacement-attorney-1-uid",
              appointmentType: "replacement",
              status: "inactive",
            },
          ],
          certificateProvider: {
            uid: "certificate-provider",
            signedAt: "2025-09-30T10:25:47.881438922Z",
          },
          donor: {
            uid: "donor",
          },
        },
      },
    });

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
          {
            id: 129,
            status: "detected",
            fieldName: "address",
            ruleType: "Invalid address",
            fieldOwnerUid: "certificate-provider",
          },
          {
            id: 130,
            status: "detected",
            fieldName: "address",
            ruleType: "Invalid address",
            fieldOwnerUid: "attorney-1-uid",
          },
          {
            id: 131,
            status: "detected",
            fieldName: "address",
            ruleType: "Invalid address",
            fieldOwnerUid: "replacement-attorney-1-uid",
          },
          {
            id: 132,
            status: "detected",
            fieldName: "firstNames",
            ruleType: "empty",
            fieldOwnerUid: "certificate-provider",
          },
          {
            id: 133,
            status: "detected",
            fieldName: "lastName",
            ruleType: "empty",
            fieldOwnerUid: "certificate-provider",
          },
          {
            id: 138,
            status: "detected",
            fieldName: "address",
            ruleType: "no-country",
            fieldOwnerUid: "donor",
          },
          {
            id: 139,
            status: "detected",
            fieldName: "signedAt",
            ruleType:
              "Donor: date of ID and date of signature more than 6-months apart",
            fieldOwnerUid: "donor",
          },
          {
            id: 140,
            status: "detected",
            fieldName: "signedAt",
            ruleType:
              "Certificate-provider: date of ID and date of signature more than 6-months apart",
            fieldOwnerUid: "certificate-provider",
          },
          {
            id: 141,
            status: "detected",
            fieldName: "signedAt",
            ruleType:
              "certificate-provider signature more than 2-years after donor",
            fieldOwnerUid: "certificate-provider",
          },
          {
            id: 142,
            status: "detected",
            fieldName: "signedAt",
            ruleType:
              "attorney signature more than 2-years after donor",
            fieldOwnerUid: "attorney-1-uid",
          },
        ],
      },
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-TTTT-3333", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-TTTT-3333",
        "opg.poas.sirius": {
          id: 111,
          uId: "M-DIGI-TTTT-3333",
          status: "Processing",
          caseSubtype: "property-and-affairs",
        },
        "opg.poas.lpastore": {
          channel: "online",
          attorneys: [
            {
              uid: "attorney-1-uid",
              appointmentType: "original",
              status: "active",
            },
            {
              uid: "replacement-attorney-1-uid",
              appointmentType: "replacement",
              status: "inactive",
            },
          ],
          certificateProvider: {
            uid: "certificate-provider",
          },
        },
      },
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-TTTT-3333/anomalies", "GET", {
      status: 200,
      body: {
        uid: "M-DIGI-TTTT-3333",
        anomalies: [
          {
            id: 136,
            status: "detected",
            fieldName: "lastName",
            ruleType: "last-name-matches-donor",
            fieldOwnerUid: "certificate-provider",
          },
          {
            id: 137,
            status: "detected",
            fieldName: "lastName",
            ruleType: "last-name-matches-attorney",
            fieldOwnerUid: "certificate-provider",
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

    const mocks = Promise.allSettled([
      cases.warnings.empty("111"),
      cases.warnings.empty("222"),
      cases.tasks.empty("111"),
      cases.tasks.empty("222"),
      digitalLpas.objections.empty("M-DIGI-QQQQ-1111"),
      digitalLpas.objections.empty("M-DIGI-SSSS-3333"),
      digitalLpas.objections.empty("M-DIGI-TTTT-3333"),
    ]);

    cy.wrap(mocks);

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
    cy.contains("Review address as there is no country");
    cy.contains(
      "Review signature date - check this is within 6 months either side of the donor’s ID check",
    );
    cy.contains("Review attorney's first names");
    cy.contains("Review attorney's last name");
    cy.contains("Review attorney's address");
    cy.contains("Review signature date - check this is within 2 years of the donor signing the LPA");
    cy.contains("Review replacement attorney's first names");
    cy.contains("Review replacement attorney's last name");
    cy.contains("Review replacement attorney's address");
    cy.contains("Review how attorneys can make decisions");
    cy.contains("Review when the LPA can be used");
    cy.contains("Review certificate provider's first names");
    cy.contains("Review certificate provider's last name");
    cy.contains("Review certificate provider's address");
    cy.contains(
      "Review signature date - check this is within 6 months either side of the certificate provider’s ID check and within 2 years of the donor signing the LPA",
    );
  });

  it("shows anomalies for pa LPA", () => {
    cy.visit("/lpa/M-DIGI-SSSS-3333/lpa-details");
    cy.contains("Some LPA details have been identified for review.");
    cy.contains("For review");
    cy.contains("Review life sustaining treatment");
  });

  it("shows anomalies for last name matching both donor and at least one attorney", () => {
    cy.visit("/lpa/M-DIGI-TTTT-3333/lpa-details");
    cy.contains(
      "Review last name - this matches the donor and at least one of the attorneys. Check certificate provider's eligibility",
    );
  });
});
