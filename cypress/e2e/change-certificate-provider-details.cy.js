describe("Change certificate provider details form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-1111-1111-1111", "GET", {
      status: 200,
      body: {
        uId: "M-1111-1111-1111",
        "opg.poas.sirius": {
          id: 1111,
          uId: "M-1111-1111-1111",
          status: "Draft",
          caseSubtype: "property-and-affairs",
          createdDate: "31/10/2023",
          investigationCount: 0,
          complaintCount: 0,
          taskCount: 0,
          warningCount: 0,
          dueDate: "01/12/2023",
          donor: {
            id: 1111,
            firstname: "Steven",
            surname: "Munnell",
            dob: "17/06/1982",
            addressLine1: "1 Scotland Street",
            addressLine2: "Netherton",
            addressLine3: "Glasgow",
            town: "Edinburgh",
            postcode: "EH6 18J",
            country: "GB",
            personType: "Donor",
          },
          certificateProvider: {
            uid: "c362e307-71b9-4070-bdde-c19b4cdf5c1a",
            channel: "online",
            firstNames: "Rhea",
            lastNames: "Vandervort",
            address: {
              line1: "290 Vivien Road",
              line2: "Lower Court",
              line3: "Tillman",
              town: "Oxfordshire",
              postcode: "JJ80 7QL",
              country: "GB",
            },
            email: "Rhea.Vandervort@example.com",
            phone: "0151 087 7256",
            signedAt: "2025-01-19T09:12:59Z",
          },
          application: {
            donorFirstNames: "Steven",
            donorLastName: "Munnell",
            donorDob: "17/06/1982",
            donorAddress: {
              addressLine1: "1 Scotland Street",
              postcode: "EH6 18J",
            },
          },
        },
        "opg.poas.lpastore": {
          lpaType: "pf",
          channel: "online",
          status: "draft",
          registrationDate: "2022-12-18",
          peopleToNotify: [],
          donor: {
            uid: "572fe550-e465-40b3-a643-ca9564fabab8",
            firstNames: "Steven",
            lastName: "Munnell",
            email: "Steven.Munnell@example.com",
            dateOfBirth: "17/06/1982",
            otherNamesKnownBy: "",
            contactLanguagePreference: "",
            address: {
              line1: "1 Scotland Street",
              line2: "Netherton",
              line3: "Glasgow",
              town: "Edinburgh",
              postcode: "EH6 18J",
              country: "GB",
            },
          },
          attorneys: [
            {
              uid: "active-attorney-1",
              firstNames: "Katheryn",
              lastName: "Collins",
              address: {
                line1: "9 O'Reilly Rise",
                line2: "Upton",
                town: "Williamsonborough",
                postcode: "ZZ24 4JM",
                country: "GB",
              },
              status: "active",
              signedAt: "2022-12-19T09:12:59Z",
              dateOfBirth: "1971-11-27",
              mobile: "0500133447",
              email: "K.Collins@example.com",
            },
          ],
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/1111/warnings", "GET", {
      status: 200,
      body: [],
    });

    cy.addMock(
      "/lpa-api/v1/cases/1111/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-1111-1111-1111/progress-indicators",
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-1111-1111-1111",
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
  });

  it("can be visited from the LPA details certificate provider Change link", () => {
    cy.visit("/lpa/M-1111-1111-1111/lpa-details").then(() => {
      cy.get(".govuk-accordion__section-button")
        .contains("Certificate provider")
        .click();
      cy.get("#f-change-certificate-provider-details").click();
      cy.contains("Change certificate provider details");
      cy.url().should(
        "contain",
        "/lpa/M-1111-1111-1111/certificate-provider/change-details",
      );
    });
  });

  it("can submit the change details form", () => {
    cy.get("#f-firstNames").should("have.value", "Rhea");
    cy.get("#f-lastName").should("have.value", "Rutherford");

    cy.get("#f-address\\.Line1").should("have.value", "15 Cameron Approach");
    cy.get("#f-address\\.Line2").should("have.value", "Lower Court");
    cy.get("#f-address\\.Line3").should("have.value", "Tillman");
    cy.get("#f-address\\.Town").should("have.value", "Oxfordshire");
    cy.get("#f-address\\.Postcode").should("have.value", "JJ80 7QL");
    cy.get("#f-address\\.Country").should("have.value", "GB");

    cy.get("#f-phoneNumber").should("have.value", "0151 087 7256");
    cy.get("#f-email").should("have.value", "Rhea.Vandervort@example.com");

    cy.get("#f-signedAt-day").should("have.value", "19");
    cy.get("#f-signedAt-month").should("have.value", "01");
    cy.get("#f-signedAt-year").should("have.value", "2025");

    cy.contains("Submit").click();
    cy.url().should("contain", "/lpa/M-1111-1111-1111/lpa-details");
  });

  it("can go Back to LPA details", () => {
    cy.contains("Back to LPA details").click();
    cy.url().should("contain", "/lpa/M-1111-1111-1111/lpa-details");
  });

  it("can be cancelled, returning to the LPA details", () => {
    cy.contains("Cancel").click();
    cy.url().should("contain", "/lpa/M-1111-1111-1111/lpa-details");
  });
});
