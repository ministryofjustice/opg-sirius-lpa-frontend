import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Manage Attorneys", () => {
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
              signedAt: "2022-12-19T09:12:59Z",
              dateOfBirth: "1971-11-27",
              mobile: "0500133447",
              email: "K.Collins@example.com",
            },
          ],
        },
      },
    });

    const mocks = Promise.allSettled([
      cases.warnings.empty("1111"),
      cases.warnings.empty("2222"),
      cases.tasks.empty("1111"),
      cases.tasks.empty("2222"),
      digitalLpas.objections.empty("M-1111-1111-1111"),
      digitalLpas.objections.empty("M-2222-2222-2222"),
    ]);

    cy.wrap(mocks);

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

    cy.visit("/lpa/M-1111-1111-1111/manage-attorneys");
  });

  it("shows the Manage Attorneys page and clicking Cancel returns to the Application progress page", () => {
    cy.contains("Manage attorneys");
    cy.get("a").contains("Cancel").click();
    cy.url()
      .should("include", "/lpa/M-1111-1111-1111")
      .should("not.include", "manage-attorneys");
  });

  it("shows an error when submitting a blank form", () => {
    cy.get("button").contains("Continue").click();
    cy.contains("There is a problem");
  });

  it("shows the Manage attorneys Case actions link for an LPA with attorneys", () => {
    cy.visit("/lpa/M-1111-1111-1111");
    cy.get(".moj-button-menu").contains("Case actions").click();
    cy.get("a").contains("Manage attorneys");
  });

  it("does not show the Manage attorneys Case actions link for an LPA without attorneys", () => {
    cy.visit("/lpa/M-2222-2222-2222");
    cy.get(".moj-button-menu").contains("Case actions").click();
    cy.get("a").contains("Manage attorneys").should("not.exist");
  });
});
