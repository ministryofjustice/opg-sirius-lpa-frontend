import * as cases from "../mocks/cases";

describe("Remove an attorney", () => {
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
              appointmentType: "original",
              signedAt: "2022-12-19T09:12:59Z",
              dateOfBirth: "1971-11-27",
              mobile: "0500133447",
              email: "K.Collins@example.com",
            },
            {
              uid: "active-attorney-2",
              firstNames: "Rachel",
              lastName: "Jones",
              address: {
                line1: "10 O'Reilly Rise",
                line2: "Upton",
                town: "Williamsonborough",
                postcode: "ZZ24 4JM",
                country: "GB",
              },
              status: "active",
              appointmentType: "replacement",
              signedAt: "2022-12-20T09:12:59Z",
              dateOfBirth: "1971-11-29",
              mobile: "0500133447",
              email: "K.Collins@example.com",
            },
            {
              uid: "inactive-attorney-1",
              firstNames: "Barry",
              lastName: "Smith",
              address: {
                line1: "11 O'Reilly Rise",
                line2: "Upton",
                town: "Williamsonborough",
                postcode: "ZZ24 4JM",
                country: "GB",
              },
              status: "inactive",
              appointmentType: "replacement",
              signedAt: "2022-12-22T09:12:59Z",
              dateOfBirth: "1971-11-30",
              mobile: "0500133447",
              email: "K.Collins@example.com",
            },
          ],
        },
      },
    });

    cases.warnings.empty("1111");

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

    cy.addMock("/lpa-api/v1/digital-lpas/M-2222-2222-2222", "GET", {
      status: 200,
      body: {
        uId: "M-2222-2222-2222",
        "opg.poas.sirius": {
          id: 2222,
          uId: "M-2222-2222-2222",
          donor: {
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
        },
      },
    });

    cy.addMock("/lpa-api/v1/cases/2222", "GET", {
      status: 200,
      body: {
        id: 1111,
        uId: "M-2222-2222-2222",
        caseType: "DIGITAL_LPA",
        donor: {
          id: 2222,
        },
        status: "Processing",
      },
    });

    cy.addMock(
      "/lpa-api/v1/cases/2222/tasks?filter=status%3ANot+started%2Cactive%3Atrue&limit=99&sort=duedate%3AASC",
      "GET",
      {
        status: 200,
        body: {
          tasks: [],
        },
      },
    );

    cases.warnings.empty("2222");

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-2222-2222-2222/progress-indicators",
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-2222-2222-2222",
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

    cy.visit("/lpa/M-1111-1111-1111/remove-an-attorney");
  });

  it("shows the Remove an attorney page and clicking Cancel returns to the Application progress page", () => {
    cy.contains("Remove an attorney");
    cy.get("a").contains("Cancel").click();
    cy.url()
      .should("include", "/lpa/M-1111-1111-1111")
      .should("not.include", "remove-an-attorney");
  });

  it("shows the Remove an attorney page and clicking Back returns to the Application progress page", () => {
    cy.contains("Remove an attorney");
    cy.get("a").contains("Back").click();
    cy.url()
      .should("include", "/lpa/M-1111-1111-1111/manage-attorneys")
      .should("not.include", "remove-an-attorney");
  });

  it("shows the Remove an attorney page with active attorneys", () => {
    cy.contains("Remove an attorney");
    cy.get('input[name="confirmRemoval"]').should("not.exist");
    cy.get(".govuk-label").contains("Katheryn Collins");
    cy.get(".govuk-label").contains("Rachel Jones");
    cy.get(".govuk-label").contains("Barry Smith").should("not.exist");
  });

  it("shows an error when submitting a blank Remove an attorney form", () => {
    cy.get("button").contains("Continue").click();
    cy.contains("There is a problem");
  });

  it("shows the Confirm removal of attorney page when submitting the Remove an attorney form with an active attorney selected", () => {
    cy.contains("Remove an attorney");
    cy.get('input[name="confirmRemoval"]').should("not.exist");
    cy.get("#f-attorney-1").click();
    cy.get("button").contains("Continue").click();
    cy.url().should("include", "/lpa/M-1111-1111-1111/remove-an-attorney");
    cy.contains("Confirm removal of attorney");
    cy.get(".govuk-summary-list__value").contains("Katheryn Collins");
  });
});
