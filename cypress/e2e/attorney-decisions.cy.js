import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("Manage attorney decisions form", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/digital-lpas/M-DIGI-LPA1-1111", "GET", {
      status: 200,
      body: {
        uId: "M-DIGI-LPA1-1111",
        "opg.poas.sirius": {
          id: 1111,
          uId: "M-DIGI-LPA1-1111",
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
          howAttorneysMakeDecisions: "jointly-for-some-severally-for-others",
          howReplacementAttorneysMakeDecisions:
            "jointly-for-some-severally-for-others",
          donor: {
            uid: "572fe550-e465-40b3-a643-ca9564fabab8",
            firstNames: "Steven",
            lastName: "Munnell",
            email: "Steven.Munnell@example.com",
            dateOfBirth: "17/06/1982",
            contactLanguagePreference: "",
            address: {
              line1: "1 Scotland Street",
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
                postcode: "ZZ24 4JM",
                country: "GB",
              },
              status: "active",
              appointmentType: "original",
              signedAt: "2022-12-19T09:12:59Z",
              dateOfBirth: "1971-11-27",
              howAttorneysMakeDecisions: false,
            },
            {
              uid: "active-attorney-2",
              firstNames: "Rachel",
              lastName: "Jones",
              address: {
                line1: "10 O'Reilly Rise",
                postcode: "ZZ24 4JM",
                country: "GB",
              },
              status: "active",
              appointmentType: "replacement",
              signedAt: "2022-12-20T09:12:59Z",
              dateOfBirth: "1971-11-29",
              howAttorneysMakeDecisions: false,
            },
            {
              uid: "inactive-attorney-1",
              firstNames: "Barry",
              lastName: "Smith",
              address: {
                line1: "11 O'Reilly Rise",
                postcode: "ZZ24 4JM",
                country: "GB",
              },
              status: "inactive",
              appointmentType: "replacement",
              signedAt: "2022-12-22T09:12:59Z",
              dateOfBirth: "1971-11-30",
              howAttorneysMakeDecisions: false,
            },
          ],
        },
      },
    });

    const mocks = Promise.allSettled([
      cases.warnings.empty("1111"),
      cases.tasks.empty("1111"),
      digitalLpas.objections.empty("M-DIGI-LPA1-1111"),
    ]);

    cy.wrap(mocks);

    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-DIGI-LPA1-1111/progress-indicators",
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: "M-DIGI-LPA1-1111",
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

    cy.visit("/lpa/M-DIGI-LPA1-1111/manage-attorney-decisions");
  });

  it("clicking Back returns to the manage attorneys page", () => {
    cy.contains("Manage decisions - attorneys who cannot act jointly");
    cy.get("a").contains("Back").click();
    cy.url()
      .should("include", "/lpa/M-DIGI-LPA1-1111/manage-attorneys")
      .should("not.include", "manage-attorney-decisions");
  });

  it("clicking Cancel returns to the lpa page", () => {
    cy.contains("Manage decisions - attorneys who cannot act jointly");
    cy.get("a").contains("Cancel").click();
    cy.url()
      .should("include", "/lpa/M-DIGI-LPA1-1111")
      .should("not.include", "manage-attorney-decisions");
  });

  it("shows the attorney appointment type and attorney details", () => {
    cy.contains("Manage decisions - attorneys who cannot act jointly");
    cy.contains("Attorneys appointment type");
    cy.contains("Replacement attorney appointment type");
    cy.contains("Jointly for some, severally for others");
    cy.contains("Select who cannot make joint decisions");
    cy.contains("Katheryn Collins (attorney)");
    cy.contains("Rachel Jones (previously replacement attorney)");
    cy.contains("Joint decisions can be made by all attorneys");
  });

  it("shows an error when the user doesn't select an option", () => {
    cy.get("button").contains("Continue").click();
    cy.contains("There is a problem");
  });

  it("Can complete form journey - attorneys who cannot act jointly", () => {
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-DIGI-LPA1-1111/attorney-decisions",
      "PUT",
      {
        status: 204,
        body: {
          attorneyDecisions: [
            { uid: "active-attorney-1", cannotMakeJointDecisions: false },
            { uid: "active-attorney-2", cannotMakeJointDecisions: false },
          ],
          receivedDate: "13/12/2024",
          objectionType: "prescribed",
          notes: "Test",
        },
      },
    );

    cy.contains("Manage decisions - attorneys who cannot act jointly");
    cy.get("#f-skipDecisionAttorney").click();
    cy.get("button").contains("Continue").click();
    cy.url().should(
      "include",
      "/lpa/M-DIGI-LPA1-1111/manage-attorney-decisions",
    );
    cy.contains("Confirm who cannot make joint decisions");
    cy.contains("Decisions");
    cy.contains("Attorney appointment type");
    cy.contains("Jointly for some, severally for others");
    cy.contains("Confirm joint decision making").click();
    cy.url()
      .should("include", "/lpa/M-DIGI-LPA1-1111")
      .should("not.include", "manage-attorney-decisions");
    cy.contains("Update saved");
  });

  it("Can complete form journey - attorneys selected", () => {
    cy.addMock(
      "/lpa-api/v1/digital-lpas/M-DIGI-LPA1-1111/attorney-decisions",
      "PUT",
      {
        status: 204,
        body: {
          attorneyDecisions: [
            { uid: "active-attorney-1", cannotMakeJointDecisions: true },
            { uid: "active-attorney-2", cannotMakeJointDecisions: true },
          ],
          receivedDate: "13/12/2024",
          objectionType: "prescribed",
          notes: "Test",
        },
      },
    );

    cy.contains("Manage decisions - attorneys who cannot act jointly");
    cy.get("#f-activeAttorney-1").click();
    cy.get("#f-activeAttorney-2").click();
    cy.get("button").contains("Continue").click();
    cy.url().should(
      "include",
      "/lpa/M-DIGI-LPA1-1111/manage-attorney-decisions",
    );
    cy.contains("Katheryn Collins (attorney)");
    cy.contains("Rachel Jones (previously replacement attorney)");
    cy.contains("Confirm joint decision making").click();
    cy.url()
      .should("include", "/lpa/M-DIGI-LPA1-1111")
      .should("not.include", "manage-attorney-decisions");
    cy.contains("Update saved");
  });
});
