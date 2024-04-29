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
        "opg.poas.lpastore": null,
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

    cy.addMock(`/lpa-api/v1/digital-lpas/${uid}/progress-indicators`, "GET", {
      status: 200,
      body: {
        digitalLpaUid: uid,
        progressIndicators: [
          { indicator: "FEES", status: "NOT_STARTED" },
          { indicator: "FEES", status: "COMPLETE" },
          { indicator: "FEES", status: "IN_PROGRESS" },
        ],
      },
    });
  });

  it("shows application progress", () => {
    cy.visit(`/lpa/${uid}`);
  });
});
