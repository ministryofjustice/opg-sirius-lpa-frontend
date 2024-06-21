describe("View the application progress for a digital LPA", () => {
  const uid = "M-QEQE-EEEE-WERT";
  const id = 113222;

  beforeEach(() => {
    cy.addMock(`/lpa-api/v1/digital-lpas/${uid}`, "GET", {
      status: 200,
      body: {
        uId: uid,
        "opg.poas.sirius": {
          id: id,
          uId: "M-QEQE-EEEE-WERT",
          status: "Processing",
          caseSubtype: "personal-welfare",
          createdDate: "24/04/2024",
          investigationCount: 2,
          complaintCount: 1,
          taskCount: 2,
          warningCount: 0,
          donor: {
            id: 33,
          },
          application: {
            donorFirstNames: "Peter",
            donorLastName: "Maaaabbbb",
            donorDob: "27/05/1968",
            donorEmail: "peter@bbsssnssssss.org",
            donorPhone: "073656249524",
            donorAddress: {
              addressLine1: "Flat 9999",
              addressLine2: "Flaim House",
              addressLine3: "33 Marb Road",
              country: "GB",
              postcode: "X15 3XX",
            },
            correspondentFirstNames: "Salty",
            correspondentLastName: "McNab",
            correspondentAddress: {
              addressLine1: "Flat 3",
              addressLine2: "Digital LPA Avenue",
              addressLine3: "Noplace",
              country: "GB",
              postcode: "SW1 1AA",
            },
          },
          linkedDigitalLpas: [],
        },
        "opg.poas.lpastore": {
          attorneys: [
            {
              firstNames: "Esther",
              lastName: "Greenwood",
              status: "active",
            },
          ],
          certificateProvider: {
            channel: "paper",
            firstNames: "Fake",
            lastName: "Provider",
          },
          lpaType: "pf",
          channel: "paper",
          registrationDate: "2022-12-18",
          peopleToNotify: [],
        },
      },
    });

    cy.addMock(`/lpa-api/v1/cases/${id}/warnings`, "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      `/lpa-api/v1/cases/${id}/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC`,
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );
  });

  it("shows application progress not started", () => {
    cy.addMock(`/lpa-api/v1/digital-lpas/${uid}/progress-indicators`, "GET", {
      status: 200,
      body: {
        digitalLpaUid: uid,
        progressIndicators: [
          { indicator: "FEES", status: "CANNOT_START" },
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
    });

    cy.visit(`/lpa/${uid}`);

    cy.get(".app-progress-indicator-summary").then((elts) => {
      expect(
        Cypress.$(elts[0]).find("svg[data-progress-indicator=not-started]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[1]).find("svg[data-progress-indicator=not-started]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[2]).find("svg[data-progress-indicator=not-started]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[3]).find("svg[data-progress-indicator=not-started]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[4]).find("svg[data-progress-indicator=not-started]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[5]).find("svg[data-progress-indicator=not-started]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[6]).find("svg[data-progress-indicator=not-started]")
          .length,
      ).to.equal(1);
    });

    cy.contains("Donor identity confirmation").click();
    cy.contains("Not started");
    cy.contains("Start donor identity check").should("not.exist");

    cy.contains("Certificate provider identity confirmation").click();
    cy.contains("Start certificate provider identity check").should(
      "not.exist",
    );
  });

  it("shows application progress in progress", () => {
    cy.addMock(`/lpa-api/v1/digital-lpas/${uid}/progress-indicators`, "GET", {
      status: 200,
      body: {
        digitalLpaUid: uid,
        progressIndicators: [
          { indicator: "FEES", status: "IN_PROGRESS" },
          { indicator: "DONOR_ID", status: "IN_PROGRESS" },
          { indicator: "CERTIFICATE_PROVIDER_ID", status: "IN_PROGRESS" },
          {
            indicator: "CERTIFICATE_PROVIDER_SIGNATURE",
            status: "IN_PROGRESS",
          },
          { indicator: "ATTORNEY_SIGNATURES", status: "IN_PROGRESS" },
          { indicator: "PREREGISTRATION_NOTICES", status: "IN_PROGRESS" },
          { indicator: "REGISTRATION_NOTICES", status: "IN_PROGRESS" },
        ],
      },
    });

    cy.visit(`/lpa/${uid}`);

    cy.get(".app-progress-indicator-summary").then((elts) => {
      expect(
        Cypress.$(elts[0]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[1]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[2]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[3]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[4]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[5]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[6]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
    });

    cy.contains("Donor identity confirmation").click();
    cy.contains("In progress");
    cy.contains("Start donor identity check").should("exist");
    cy.get(
      "a[href='/lpa/identity-check/start?personType=donor&lpas[]=M-QEQE-EEEE-WERT']",
    ).contains("Start donor identity check");

    cy.contains("Certificate provider identity confirmation").click();
    cy.contains("Certificate provider: Fake Provider");
    cy.contains("Start certificate provider identity check").should("exist");
    cy.get(
      "a[href='/lpa/identity-check/start?personType=certificate-provider&lpas[]=M-QEQE-EEEE-WERT']",
    ).contains("Start certificate provider identity check");
  });

  it("shows application progress completed", () => {
    const uid = "M-QEQE-EEEE-WERT";

    cy.addMock(`/lpa-api/v1/digital-lpas/${uid}/progress-indicators`, "GET", {
      status: 200,
      body: {
        digitalLpaUid: uid,
        progressIndicators: [
          { indicator: "FEES", status: "COMPLETE" },
          { indicator: "DONOR_ID", status: "COMPLETE" },
          { indicator: "CERTIFICATE_PROVIDER_ID", status: "COMPLETE" },
          {
            indicator: "CERTIFICATE_PROVIDER_SIGNATURE",
            status: "COMPLETE",
          },
          { indicator: "ATTORNEY_SIGNATURES", status: "COMPLETE" },
          { indicator: "PREREGISTRATION_NOTICES", status: "COMPLETE" },
          { indicator: "REGISTRATION_NOTICES", status: "COMPLETE" },
        ],
      },
    });

    cy.visit(`/lpa/${uid}`);

    cy.get(".app-progress-indicator-summary").then((elts) => {
      expect(
        Cypress.$(elts[0]).find("svg[data-progress-indicator=complete]").length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[1]).find("svg[data-progress-indicator=complete]").length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[2]).find("svg[data-progress-indicator=complete]").length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[3]).find("svg[data-progress-indicator=complete]").length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[4]).find("svg[data-progress-indicator=complete]").length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[5]).find("svg[data-progress-indicator=complete]").length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[6]).find("svg[data-progress-indicator=complete]").length,
      ).to.equal(1);
    });

    cy.contains("Donor identity confirmation").click();
    cy.contains("Complete");
    cy.contains("Start donor identity check").should("not.exist");

    cy.contains("Certificate provider identity confirmation").click();
    cy.contains("Certificate provider: Fake Provider");
    cy.contains("Start certificate provider identity check").should(
      "not.exist",
    );
  });

  it("does not show link to id check when actor is online", () => {
    cy.addMock(`/lpa-api/v1/digital-lpas/${uid}`, "GET", {
      status: 200,
      body: {
        uId: uid,
        "opg.poas.sirius": {
          id: id,
          uId: "M-QEQE-EEEE-WERT",
          status: "Processing",
          caseSubtype: "personal-welfare",
          createdDate: "24/04/2024",
          investigationCount: 2,
          complaintCount: 1,
          taskCount: 2,
          warningCount: 0,
          donor: {
            id: 33,
          },
          application: {
            donorFirstNames: "Peter",
            donorLastName: "Maaaabbbb",
            donorDob: "27/05/1968",
            donorEmail: "peter@bbsssnssssss.org",
            donorPhone: "073656249524",
            donorAddress: {
              addressLine1: "Flat 9999",
              addressLine2: "Flaim House",
              addressLine3: "33 Marb Road",
              country: "GB",
              postcode: "X15 3XX",
            },
            correspondentFirstNames: "Salty",
            correspondentLastName: "McNab",
            correspondentAddress: {
              addressLine1: "Flat 3",
              addressLine2: "Digital LPA Avenue",
              addressLine3: "Noplace",
              country: "GB",
              postcode: "SW1 1AA",
            },
          },
          linkedDigitalLpas: [],
        },
        "opg.poas.lpastore": {
          attorneys: [
            {
              firstNames: "Esther",
              lastName: "Greenwood",
              status: "active",
            },
          ],
          certificateProvider: {
            channel: "online",
            firstNames: "Fake",
            lastName: "Provider",
          },
          lpaType: "pf",
          channel: "online",
          registrationDate: "2022-12-18",
          peopleToNotify: [],
        },
      },
    });

    cy.addMock(`/lpa-api/v1/digital-lpas/${uid}/progress-indicators`, "GET", {
      status: 200,
      body: {
        digitalLpaUid: uid,
        progressIndicators: [
          { indicator: "DONOR_ID", status: "IN_PROGRESS" },
          { indicator: "CERTIFICATE_PROVIDER_ID", status: "IN_PROGRESS" },
        ],
      },
    });

    cy.visit(`/lpa/${uid}`);

    cy.get(".app-progress-indicator-summary").then((elts) => {
      expect(
        Cypress.$(elts[0]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
      expect(
        Cypress.$(elts[1]).find("svg[data-progress-indicator=in-progress]")
          .length,
      ).to.equal(1);
    });

    cy.contains("Donor identity confirmation").click();
    cy.contains("Start donor identity check").should("not.exist");

    cy.contains("Certificate provider identity confirmation").click();
    cy.contains("Start certificate provider identity check").should(
      "not.exist",
    );
  });
});
