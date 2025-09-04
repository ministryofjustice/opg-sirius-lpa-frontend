import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

// helper function to create mock data for digital LPA - allows overriding of sirius and lpaStore data
// to allow testing of different scenarios without having to create a new mock file for each one
function createDigitalLpaM3333MockData(
  lpaStoreOverrides = {},
  siriusOverrides = {},
) {
  return {
    uId: "M-DIGI-LPA3-3333",
    "opg.poas.sirius": {
      id: 333,
      uId: "M-DIGI-LPA3-3333",
      status: "Draft",
      caseSubtype: "property-and-affairs",
      createdDate: "31/10/2023",
      investigationCount: 2,
      complaintCount: 1,
      taskCount: 2,
      warningCount: 4,
      dueDate: "01/12/2023",
      donor: {
        id: 33,
      },
      application: {
        donorFirstNames: "Agnes",
        donorLastName: "Hartley",
        donorDob: "27/05/1998",
        donorEmail: "agnes@host.example",
        donorPhone: "073656249524",
        donorAddress: {
          addressLine1: "Apartment 3",
          addressLine2: "Gherkin Building",
          addressLine3: "33 London Road",
          country: "GB",
          postcode: "B15 3AA",
          town: "Birmingham",
        },
        correspondentFirstNames: "Kendrick",
        correspondentLastName: "Lamar",
        correspondentAddress: {
          addressLine1: "Flat 3",
          addressLine2: "Digital LPA Lane",
          addressLine3: "Somewhere",
          country: "GB",
          postcode: "SW1 1AA",
          town: "London",
        },
      },
      linkedDigitalLpas: [
        {
          uId: "M-DIGI-LPA3-3334",
          caseSubtype: "personal-welfare",
          status: "Draft",
          createdDate: "01/11/2023",
        },
        {
          uId: "M-DIGI-LPA3-3335",
          caseSubtype: "personal-welfare",
          status: "Registered",
          createdDate: "02/11/2023",
        },
      ],
      ...siriusOverrides,
    },
    "opg.poas.lpastore": {
      donor: {
        uid: "donor-uid-333",
        firstNames: "Agnes",
        lastName: "Hartley",
        dateOfBirth: "1998-05-27",
        otherNamesKnownBy: "",
        address: {
          line1: "Apartment 3",
          line2: "Gherkin Building",
          line3: "33 London Road",
          town: "Birmingham",
          postcode: "B15 3AA",
          country: "GB",
        },
        email: "agnes@host.example",
        contactLanguagePreference: "en",
      },
      attorneys: [
        {
          firstNames: "Esther",
          lastName: "Greenwood",
          status: "active",
          appointmentType: "original",
          cannotMakeJointDecisions: true,
        },
        {
          firstNames: "Volo",
          lastName: "McSpolo",
          status: "active",
          appointmentType: "original",
          signedAt: "2022-12-20T12:02:43Z",
        },
        {
          firstNames: "Susanna",
          lastName: "Kaysen",
          status: "removed",
          appointmentType: "original",
        },
        {
          firstNames: "Philomena",
          lastName: "Guinea",
          status: "inactive",
          appointmentType: "replacement",
        },
        {
          firstNames: "Rico",
          lastName: "Welch",
          status: "inactive",
          appointmentType: "replacement",
          signedAt: "2022-12-19T09:12:59Z",
        },
        {
          firstNames: "Anne",
          lastName: "Rice",
          status: "active",
          appointmentType: "replacement",
          signedAt: "2022-12-19T07:18:59Z",
          cannotMakeJointDecisions: true,
        },
      ],
      trustCorporations: [
        {
          Name: "Trust Me Ltd.",
          CompanyNumber: "123456789",
          status: "active",
          appointmentType: "original",
        },
        {
          Name: "Trust Me Again Ltd.",
          CompanyNumber: "987654321",
          status: "inactive",
          appointmentType: "replacement",
        },
      ],
      certificateProvider: {
        uid: "e4d5e24e-2a8d-434e-b815-9898620acc71",
        firstNames: "Timothy",
        lastName: "Turner",
        signedAt: "2022-12-18T11:46:24Z",
      },
      lpaType: "pf",
      channel: "online",
      status: "draft",
      registrationDate: "2022-12-18",
      peopleToNotify: [],
      restrictionsAndConditions: "Do not do this",
      lifeSustainingTreatmentOption: "option-a",
      howAttorneysMakeDecisions: "jointly-for-some-severally-for-others",
      howAttorneysMakeDecisionsDetails:
        "My attorneys must act jointly to decide whether...",
      howReplacementAttorneysMakeDecisions:
        "jointly-for-some-severally-for-others",
      signedAt: "2022-12-15T09:00:00Z",
      authorisedSignatory: {
        firstNames: "John",
        lastName: "Signatory",
        signedAt: "2022-12-15T10:30:00Z",
      },
      witnessedByCertificateProviderAt: "2022-12-15T11:00:00Z",
      witnessedByIndependentWitnessAt: "2022-12-15T11:30:00Z",
      independentWitness: {
        firstNames: "Jane",
        lastName: "Witness",
        address: {
          line1: "123 Witness Street",
          line2: "",
          line3: "",
          town: "London",
          postcode: "SW1A 1AA",
          country: "GB",
        },
        email: "jane.witness@example.com",
      },
      ...lpaStoreOverrides,
    },
  };
}

describe("View a digital LPA", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333", "GET", {
      status: 200,
      body: createDigitalLpaM3333MockData(),
    });

    cy.addMock("/lpa-api/v1/cases/333", "GET", {
      status: 200,
      body: {
        id: 333,
        uId: "M-DIGI-LPA3-3333",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 33,
        },
        status: "Processing",
      },
    });

    cy.addMock(
      "/lpa-api/v1/persons/33/events?filter=case:333&sort=id:desc",
      "GET",
      {
        status: 200,
        body: {
          events: [
            {
              id: 111111,
              user: {
                id: 11,
                phoneNumber: "12345678",
                teams: [],
                displayName: "system admin",
                deleted: false,
                email: "system.admin@opgtest.com",
              },
              sourceType: "Donor",
              sourcePerson: {
                id: 111111,
                uId: "7000-1111-1111",
                firstname: "John",
                surname: "Smith",
              },
              type: "INS",
              changeSet: [],
              entity: {
                _class: "Opg\\Core\\Model\\Entity\\CaseActor\\Donor",
                email: "",
                firstname: "John",
                id: 111111,
                salutation: "",
                surname: "Smith",
                uId: 700011111111,
              },
              createdOn: "2024-01-02T12:13:14+00:00",
              hash: "5555",
            },
          ],
        },
      },
    );

    cy.addMock("/lpa-api/v1/cases/333/warnings", "GET", {
      status: 200,
      body: [
        {
          id: 44,
          warningType: "Court application in progress",
          warningText: "Court notified",
          dateAdded: "24/08/2022 13:13:13",
          caseItems: [
            { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
          ],
        },
        {
          id: 22,
          warningType: "Complaint Received",
          warningText: "Complaint from donor",
          dateAdded: "12/12/2023 12:12:12",
          caseItems: [
            { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
            { uId: "M-DIGI-LPA3-5555", caseSubtype: "property-and-affairs" },
          ],
        },
        {
          id: 24,
          warningType: "Donor Deceased",
          warningText: "Advised of donor death",
          dateAdded: "05/01/2022 10:10:00",
          caseItems: [
            { uId: "M-DIGI-LPA3-3333", caseSubtype: "personal-welfare" },
            { uId: "M-DIGI-LPA3-5555", caseSubtype: "property-and-affairs" },
            { uId: "M-DIGI-LPA3-6666", caseSubtype: "personal-welfare" },
          ],
        },
      ],
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333/objections", "GET", {
      status: 200,
      body: [
        {
          id: 12,
          notes: "",
          objectionType: "factual",
          receivedDate: "2025-01-01",
        },
      ],
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333/anomalies", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      "/lpa-api/v1/lpas/333/documents?type[-][]=Draft&type[-][]=Preview",
      "GET",
      {
        status: 200,
        body: [
          {
            id: 1,
            uuid: "7327f57d-e3d5-4300-95a8-67b3337c7231",
            friendlyDescription: "Mr Test Person - Blank Template",
            direction: "Outgoing",
            createdDate: "24/08/2023 15:27:16",
            systemType: "EP-BB",
            correspondent: {
              uId: "7000-0000-0013",
              firstname: "Test",
              surname: "Person",
              personType: "Donor",
            },
          },
          {
            id: 2,
            uuid: "40fa2847-27ae-4976-a93a-9f45ec0a4e98",
            friendlyDescription: "Mr John Doe - Reduced fee evidence",
            direction: "Incoming",
            createdDate: "15/05/2023 11:09:28",
            receivedDateTime: "15/05/2023 11:09:28",
            type: "Application Related",
            subType: "Reduced fee request evidence",
            correspondent: {
              uId: "7000-0000-0013",
              firstname: "John",
              surname: "Doe",
              personType: "Correspondent",
            },
          },
        ],
      },
    );

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3334", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-LPA3-3334",
        "opg.poas.sirius": {
          id: 334,
          uId: "M-DIGI-LPA3-3334",
          status: "Draft",
          caseSubtype: "property-and-affairs",
          createdDate: "31/10/2023",
          investigationCount: 2,
          complaintCount: 1,
          taskCount: 2,
          warningCount: 4,
          donor: {
            id: 33,
          },
          application: {
            donorFirstNames: "Agnes",
            donorLastName: "Hartley",
            donorDob: "27/05/1998",
            donorEmail: "agnes@host.example",
            donorPhone: "073656249524",
            donorAddress: {
              addressLine1: "Apartment 3",
              addressLine2: "Gherkin Building",
              addressLine3: "33 London Road",
              country: "GB",
              postcode: "B15 3AA",
              town: "Birmingham",
            },
            correspondentFirstNames: "Kendrick",
            correspondentLastName: "Lamar",
            correspondentAddress: {
              addressLine1: "Flat 3",
              addressLine2: "Digital LPA Lane",
              addressLine3: "Somewhere",
              country: "GB",
              postcode: "SW1 1AA",
              town: "London",
            },
          },
          linkedDigitalLpas: [
            {
              uId: "M-DIGI-LPA3-3333",
              caseSubtype: "property-and-affairs",
              status: "Draft",
              createdDate: "01/11/2023",
            },
            {
              uId: "M-DIGI-LPA3-3335",
              caseSubtype: "personal-welfare",
              status: "Registered",
              createdDate: "02/11/2023",
            },
          ],
        },
        "opg.poas.lpastore": null,
      },
    });

    cy.addMock("/lpa-api/v1/cases/334", "GET", {
      status: 200,
      body: {
        id: 334,
        uId: "M-DIGI-LPA3-3334",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 33,
        },
        status: "Draft",
      },
    });

    const mocks = Promise.allSettled([
      cases.warnings.empty("333"),
      cases.warnings.empty("334"),
      cases.warnings.empty("336"),
      cases.tasks.empty("334"),
      cases.tasks.empty("336"),
      digitalLpas.objections.empty("M-DIGI-LPA3-3334"),
      digitalLpas.objections.empty("M-DIGI-LPA3-3336"),
    ]);

    cy.wrap(mocks);

    cy.addMock(
      `/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333/progress-indicators`,
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-DIGI-LPA3-3333",
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

    cy.addMock("/lpa-api/v1/tasks/1", "GET", {
      status: 200,
      body: {
        caseItems: [
          {
            caseType: "DIGITAL_LPA",
            uId: "M-DIGI-LPA3-3333",
          },
        ],
        dueDate: "10/01/2022",
        id: 1,
        name: "Create physical case file",
        status: "Not Started",
      },
    });

    cy.addMock(
      `/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3334/progress-indicators`,
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-DIGI-LPA3-3334",
          progressIndicators: [
            { indicator: "FEES", status: "COMPLETE" },
            { indicator: "DONOR_ID", status: "IN_PROGRESS" },
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

    const siriusPaperLpa = {
      id: 336,
      uId: "M-DIGI-LPA3-3336",
      status: "Draft",
      caseSubtype: "property-and-affairs",
      createdDate: "31/10/2023",
      investigationCount: 2,
      complaintCount: 1,
      taskCount: 2,
      warningCount: 4,
      donor: {
        id: 33,
      },
      application: {
        donorFirstNames: "Agnes",
        donorLastName: "Hartley",
        donorDob: "27/05/1998",
        donorEmail: "agnes@host.example",
        donorPhone: "073656249524",
        donorAddress: {
          addressLine1: "Apartment 3",
          addressLine2: "Gherkin Building",
          addressLine3: "33 London Road",
          country: "GB",
          postcode: "B15 3AA",
          town: "Birmingham",
        },
        correspondentFirstNames: "Kendrick",
        correspondentLastName: "Lamar",
        correspondentAddress: {
          addressLine1: "Flat 3",
          addressLine2: "Digital LPA Lane",
          addressLine3: "Somewhere",
          country: "GB",
          postcode: "SW1 1AA",
          town: "London",
        },
      },
      linkedDigitalLpas: [
        {
          uId: "M-DIGI-LPA3-3333",
          caseSubtype: "property-and-affairs",
          status: "Draft",
          createdDate: "01/11/2023",
        },
        {
          uId: "M-DIGI-LPA3-3334",
          caseSubtype: "personal-welfare",
          status: "Draft",
          createdDate: "01/11/2023",
        },
        {
          uId: "M-DIGI-LPA3-3335",
          caseSubtype: "personal-welfare",
          status: "Registered",
          createdDate: "02/11/2023",
        },
      ],
    };

    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3336", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-LPA3-3336",
        "opg.poas.sirius": siriusPaperLpa,
        "opg.poas.lpastore": {
          channel: "paper",
          restrictionsAndConditionsImages: [
            {
              path: "just-an-unsigned-url.jpg",
            },
          ],
        },
      },
    });

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3336?presignImages",
      "GET",
      {
        status: 200,
        body: {
          uId: "M-DIGI-LPA3-3336",
          "opg.poas.sirius": siriusPaperLpa,
          "opg.poas.lpastore": {
            channel: "paper",
            restrictionsAndConditionsImages: [
              {
                path: "some-presigned-url.jpg",
              },
            ],
          },
        },
      },
    );

    cy.addMock("/lpa-api/v1/cases/336", "GET", {
      status: 200,
      body: {
        id: 336,
        uId: "M-DIGI-LPA3-3336",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 33,
        },
        status: "Draft",
      },
    });

    cy.addMock(
      `/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3336/progress-indicators`,
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-DIGI-LPA3-3336",
          progressIndicators: [
            { indicator: "FEES", status: "COMPLETE" },
            { indicator: "DONOR_ID", status: "IN_PROGRESS" },
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

    const defaultDigitalLpaMocks = Promise.allSettled([
      digitalLpas.get("M-1111-1111-1111"),
      cases.warnings.empty("1111"),
      cases.tasks.empty("1111"),
      digitalLpas.progressIndicators.feesInProgress("M-1111-1111-1111"),
      digitalLpas.objections.empty("M-1111-1111-1111"),
    ]);

    cy.wrap(defaultDigitalLpaMocks);

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3337/progress-indicators",
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-DIGI-LPA3-3337",
          progressIndicators: [
            { indicator: "FEES", status: "COMPLETE" },
            { indicator: "DONOR_ID", status: "COMPLETE" },
          ],
        },
      },
    );

    cy.addMock("/lpa-api/v1/cases/337/tasks", "GET", {
      status: 200,
      body: { tasks: [] },
    });

    cy.addMock("/lpa-api/v1/cases/337/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3338/progress-indicators",
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-DIGI-LPA3-3338",
          progressIndicators: [
            { indicator: "FEES", status: "COMPLETE" },
            { indicator: "DONOR_ID", status: "COMPLETE" },
          ],
        },
      },
    );

    cy.addMock("/lpa-api/v1/cases/338/tasks", "GET", {
      status: 200,
      body: { tasks: [] },
    });

    cy.addMock("/lpa-api/v1/cases/338/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3339/progress-indicators",
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-DIGI-LPA3-3339",
          progressIndicators: [
            { indicator: "FEES", status: "COMPLETE" },
            { indicator: "DONOR_ID", status: "COMPLETE" },
          ],
        },
      },
    );

    cy.addMock("/lpa-api/v1/cases/339/tasks", "GET", {
      status: 200,
      body: { tasks: [] },
    });

    cy.addMock("/lpa-api/v1/cases/339/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-CERT-ONLY-1111", "GET", {
      status: 200,
      body: {
        uId: "M-CERT-ONLY-1111",
        "opg.poas.sirius": {
          id: 1111,
          uId: "M-CERT-ONLY-1111",
          status: "Draft",
          caseSubtype: "property-and-affairs",
          createdDate: "31/10/2023",
          investigationCount: 0,
          complaintCount: 0,
          taskCount: 0,
          warningCount: 0,
          donor: { id: 111 },
          application: {
            donorFirstNames: "Test",
            donorLastName: "Donor",
            donorDob: "01/01/1980",
            donorEmail: "test@example.com",
            donorPhone: "07123456789",
            donorAddress: {
              addressLine1: "1 Test Street",
              addressLine2: "",
              addressLine3: "",
              country: "GB",
              postcode: "T1 1ST",
              town: "Testtown",
            },
          },
        },
        "opg.poas.lpastore": {
          donor: {
            uid: "donor-uid-111",
            firstNames: "Test",
            lastName: "Donor",
            dateOfBirth: "1980-01-01",
            address: {
              line1: "1 Test Street",
              town: "Testtown",
              postcode: "T1 1ST",
              country: "GB",
            },
            email: "test@example.com",
            contactLanguagePreference: "en",
          },
          attorneys: [],
          trustCorporations: [],
          certificateProvider: {
            uid: "cert-111",
            firstNames: "Certificate",
            lastName: "Provider",
          },
          peopleToNotify: [],
          channel: "online",
          status: "draft",
          lpaType: "pf",
          restrictionsAndConditions: "",
          signedAt: "2024-01-15T10:30:00Z",
          authorisedSignatory: {
            firstNames: "Helper",
            lastName: "Person",
            signedAt: "2024-01-15T10:30:00Z",
          },
          witnessedByCertificateProviderAt: "2024-01-15T10:31:00Z",
          witnessedByIndependentWitnessAt: "",
          independentWitness: null,
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/1111", "GET", {
      status: 200,
      body: {
        id: 1111,
        uId: "M-CERT-ONLY-1111",
        caseType: "DIGITAL_LPA",
        donor: { id: 111 },
        status: "Processing",
      },
    });

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-CERT-ONLY-1111/progress-indicators",
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-CERT-ONLY-1111",
          progressIndicators: [
            { indicator: "FEES", status: "COMPLETE" },
            { indicator: "DONOR_ID", status: "COMPLETE" },
          ],
        },
      },
    );

    cy.addMock("/lpa-api/v1/cases/1111/tasks", "GET", {
      status: 200,
      body: { tasks: [] },
    });

    cy.addMock("/lpa-api/v1/cases/1111/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-CERT-ONLY-1111/objections", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock("/lpa-api/v1/digital-lpas/M-CERT-ONLY-1111/anomalies", "GET", {
      status: 200,
      body: [],
    });
  });

  it("shows case information", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.get("h1").contains("Agnes Hartley");

    cy.contains("M-DIGI-LPA3-3333");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3333'] .govuk-tag").contains("Draft");

    cy.contains("PW M-DIGI-LPA3-3334");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3334'] .govuk-tag").contains("Draft");

    cy.contains("PW M-DIGI-LPA3-3335");
    cy.get("a[href='/lpa/M-DIGI-LPA3-3335'] .govuk-tag").contains("Registered");
  });

  it("shows payment information", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("M-DIGI-LPA3-3333");
    cy.get("h1").contains("Agnes Hartley");

    cy.contains("Fees").click();
    cy.contains("£41.00 expected");
  });

  it("shows document information", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("M-DIGI-LPA3-3333");
    cy.get("h1").contains("Agnes Hartley");
    cy.contains("Documents").click();

    cy.contains("Mr Test Person - Blank Template");
    cy.contains("[OUT]");
    cy.contains("24 August 2023");
    cy.contains("EP-BB");

    cy.contains("Mr John Doe - Reduced fee evidence");
    cy.contains("[IN]");
    cy.contains("15 May 2023");
    cy.contains("Application Related");
    cy.contains("Reduced fee request evidence");
  });

  it("shows task table", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.get(
      "table[data-role=tasks-table] [data-role=tasks-table-header] tr th",
    ).should((elts) => {
      expect(elts).to.contain("Tasks");
      expect(elts).to.contain("Due date");
      expect(elts).to.contain("Actions");
    });
    cy.get(
      "table[data-role=tasks-table] tr[data-role=tasks-table-task-row]",
    ).should((elts) => {
      expect(elts).to.have.length(3);
      expect(elts).to.contain("Review reduced fee eligibility");
      expect(elts).to.contain("Review application correspondence");
      expect(elts).to.contain("Another task");
      expect(elts).to.contain("Reassign task");
    });
  });

  it("shows warnings list", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.get(".app-caseworker-summary > div:nth-child(2) li").should((elts) => {
      expect(elts).to.have.length(4);

      expect(elts[0]).to.contain("Donor Deceased");
      expect(elts[0]).to.contain(
        "this case, PA M-DIGI-LPA3-5555 and PW M-DIGI-LPA3-6666",
      );

      expect(elts[1]).to.contain("Complaint Received");
      expect(elts[1]).to.contain("this case and PA M-DIGI-LPA3-5555");

      expect(elts[2]).to.contain("Court application in progress");
      expect(elts[2]).not.to.contain("this case");
    });
  });

  it("creates a task via case actions", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains(".govuk-button", "Case actions").click();
    cy.contains("Create a task").click();
    cy.url().should("include", "/create-task?id=333");
    cy.contains("M-DIGI-LPA3-3333");
    cy.get("#f-taskType").select("Check Application");
    cy.get("#f-name").type("Do this task");
    cy.get("#f-description").type("This task, do");
    cy.contains("label", "Team").click();
    cy.get("#f-assigneeTeam").select("Cool Team");
    cy.get("#f-dueDate").type("2035-01-01");
    cy.get("button[type=submit]").click();

    cy.get(".moj-alert").should("exist");
    cy.get(".moj-alert").contains("Task created");
    cy.get("h1").contains("Agnes Hartley");
    cy.location("pathname").should("eq", "/lpa/M-DIGI-LPA3-3333");
  });

  it("creates a warning via case actions", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains(".govuk-button", "Case actions").click();
    cy.contains("Create a warning").click();
    cy.url().should("include", "/create-warning?id=33");
    cy.get("#f-warningType").select("Complaint Received");
    cy.get("#f-warningText").type("Be warned!");
    cy.get("button[type=submit]").click();

    cy.get(".moj-alert").should("exist");
    cy.get(".moj-alert").contains("Warning created");
    cy.get("h1").contains("Agnes Hartley");
    cy.location("pathname").should("eq", "/lpa/M-DIGI-LPA3-3333");
  });

  it("shows lpa details from store when status is Processing", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("LPA details").click();
    cy.contains("Attorneys (4)");
    cy.contains("Replacement attorneys (3)");
    cy.contains("Removed attorneys (1)");
    cy.contains("Notified people (0)");
    cy.contains("Correspondent");

    cy.contains("Review and confirm if severance is required").should(
      "not.exist",
    );
  });

  it("allows changing lpa decisions", () => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333/decisions", "PUT", {
      status: 200,
    });

    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("LPA details").click();
    cy.contains("Decisions").click();
    cy.get("#f-update-decisions").click();
    cy.contains("Jointly for some").click();
    cy.get("#f-howAttorneysMakeDecisionsDetails").type("This way");
    cy.contains("Continue").click();
    cy.contains("Update saved");
  });

  it("shows channel for donor", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details").then(() => {
      Cypress.$("span:contains('Donor')").closest("button")[0].click();
      cy.contains("Online");
    });
  });

  it("shows due date", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details").then(() => {
      cy.contains("Registration due: 1 December 2023");
    });
  });

  it("shows attorney signed on date and label if set", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");

    cy.contains("Attorneys (4)")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("Greenwood")
          .parents(".govuk-summary-list")
          .contains("Signed on")
          .should("not.exist");
        cy.contains("McSpolo")
          .parents(".govuk-summary-list")
          .should("contain", "Signed on")
          .and("contain", "20 December 2022");
      });
  });

  it("shows replacement attorney signed on date and label if set", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");

    cy.contains("Replacement attorneys (3)")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("Guinea")
          .parents(".govuk-summary-list")
          .contains("Signed on")
          .should("not.exist");
        cy.contains("Welch")
          .parents(".govuk-summary-list")
          .should("contain", "Signed on")
          .and("contain", "19 December 2022");
      });
  });

  it("shows certificate provider signed on date, label and change link", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");

    cy.contains("Certificate provider")
      .click()
      .parents(".govuk-accordion__section")
      .should("contain", "Signed on")
      .and("contain", "18 December 2022")
      .find("#f-change-certificate-provider-details")
      .should("contain", "Update");
  });

  it("shows application details when store is empty", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3334");

    cy.contains("LPA details").click();
    cy.contains("Application format");
    cy.contains("Paper");
  });

  it("can cancel reassign task", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("Reassign task").click();
    cy.url().should("include", "/assign-task?id=1");
    cy.contains("Assign Task");

    cy.contains("Cancel").click();
    cy.url().should("include", "/lpa/M-DIGI-LPA3-3333");
    cy.contains("Case summary");
  });

  it("can cancel clear task", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("Clear task").click();
    cy.url().should("include", "/clear-task?id=1");
    cy.contains("Save and clear task");

    cy.contains("Cancel").click();
    cy.url().should("include", "/lpa/M-DIGI-LPA3-3333");
    cy.contains("Case summary");
  });

  it("can cancel creating a warning", () => {
    cy.addMock("/lpa-api/v1/persons/33/cases", "GET", {
      status: 200,
      body: {
        cases: [
          {
            caseSubtype: "property-and-affairs",
            id: 333,
            uId: "M-DIGI-LPA3-3333",
            status: "Processing",
          },
        ],
      },
    });

    cy.visit("/lpa/M-DIGI-LPA3-3333");
    cy.contains("Case actions").click();
    cy.contains("Create a warning").click();

    cy.url().should("include", "/create-warning?id=33");
    cy.contains("Create Warning");
    cy.contains("Cancel").click();

    cy.url().should("include", "/lpa/M-DIGI-LPA3-3333");
    cy.contains("Case summary");
  });

  it("can cancel changing the status", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");
    cy.contains("Case actions").click();
    cy.contains("Change case status").click();

    cy.url().should("include", "/change-case-status?uid=M-DIGI-LPA3-3333");
    cy.contains("Change case status");
    cy.get(".govuk-button-group").contains("Cancel").click();

    cy.url().should("include", "/lpa/M-DIGI-LPA3-3333");
    cy.contains("Case summary");
  });

  it("can clear a task", () => {
    cy.addMock("/lpa-api/v1/tasks/1/mark-as-completed", "PUT", {
      status: 200,
    });
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("Clear task").click();
    cy.url().should("include", "/clear-task?id=1");
    cy.get("button[type=submit]").click();

    cy.get(".moj-alert").should("exist");
    cy.get(".moj-alert").contains("Task completed");

    cy.url().should("contain", "/lpa/M-DIGI-LPA3-3333");
    cy.contains("Case summary");
  });

  it("review severance messages appears when review restrictions tasks is open", () => {
    cy.addMock(
      "/lpa-api/v1/cases/333/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [
            {
              id: 2,
              name: "Review restrictions and conditions",
              duedate: "10/12/2023",
              status: "OPEN",
              assignee: { displayName: "Super Team" },
            },
          ],
        },
      },
    );

    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");
    cy.contains("Review and confirm if severance is required");
  });

  it("shows restrictions and conditions", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333");

    cy.contains("a", "LPA details").click();
    cy.contains("button", "Restrictions and conditions").click();

    cy.contains("Do not do this");
  });

  it("shows certificate provider witness only (1 witness)", () => {
    cy.visit("/lpa/M-CERT-ONLY-1111/lpa-details");

    cy.get(".govuk-accordion__section")
      .contains("Donor")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("LPA signed on behalf of the donor by").should("exist");
        cy.contains("Helper Person").should("exist");
      });
  });

  it("shows both witnesses when both are present (2 witnesses)", () => {
    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");

    cy.get(".govuk-accordion__section")
      .contains("Donor")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("LPA signed on behalf of the donor by").should("exist");
        cy.contains("John Signatory").should("exist");

        cy.contains(".govuk-details__summary", "View witness details").click();

        cy.get(".govuk-details__text").within(() => {
          cy.contains("Signed by witness 1 (certificate provider)").should(
            "exist",
          );
          cy.contains("Signed by witness 2").should("exist");
          cy.contains("Jane Witness").should("exist");
          cy.contains("123 Witness Street").should("exist");
        });
      });
  });

  it("shows donor signed directly (no signed on behalf)", () => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA3-3333", "GET", {
      status: 200,
      body: createDigitalLpaM3333MockData({
        authorisedSignatory: null,
        witnessedByIndependentWitnessAt: "",
        independentWitness: null,
      }),
    });

    cy.visit("/lpa/M-DIGI-LPA3-3333/lpa-details");

    cy.get(".govuk-accordion__section")
      .contains("Donor")
      .click()
      .parents(".govuk-accordion__section")
      .within(() => {
        cy.contains("LPA signed on");
        cy.contains("15 December 2022");

        cy.contains("LPA signed on behalf of the donor by").should("not.exist");
        cy.contains("View witness details").should("not.exist");
      });
  });
});
