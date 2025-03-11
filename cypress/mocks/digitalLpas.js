import { addMock } from "./wiremock";

async function get(uid, body) {
  var defaultBody = {
    uId: uid,
    "opg.poas.sirius": {
      id: 1111,
      uId: uid,
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
      certificateProvider: {
        uid: "c362e307-71b9-4070-bdde-c19b4cdf5c1a",
        channel: "online",
        firstNames: "Rhea",
        lastName: "Vandervort",
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
    },
  };

  var opgPoasSirius = defaultBody["opg.poas.sirius"];
  var opgPoasSiriusDonor = defaultBody["opg.poas.sirius"].donor;
  var opgPoasSiriusApplication = defaultBody["opg.poas.sirius"].application;
  var opgPoasLpastore = defaultBody["opg.poas.lpastore"];
  var opgPoasLpastoreDonor = defaultBody["opg.poas.lpastore"].donor;
  var opgPoasLpastoreAttorneys = defaultBody["opg.poas.lpastore"].attorneys;
  var opgPoasLpastoreCertificateProvider =
    defaultBody["opg.poas.lpastore"].certificateProvider;

  if (body !== undefined) {
    if (body.hasOwnProperty("opg.poas.sirius")) {
      opgPoasSirius = Object.assign(
        defaultBody["opg.poas.sirius"],
        body["opg.poas.sirius"],
      );
      opgPoasSiriusDonor = Object.assign(
        defaultBody["opg.poas.sirius"].donor,
        body["opg.poas.sirius"].donor ?? null,
      );
      opgPoasSiriusApplication = Object.assign(
        defaultBody["opg.poas.sirius"].application,
        body["opg.poas.sirius"].application ?? null,
      );
    }

    if (body.hasOwnProperty("opg.poas.lpastore")) {
      opgPoasLpastore = Object.assign(
        defaultBody["opg.poas.lpastore"],
        body["opg.poas.lpastore"],
      );
      opgPoasLpastoreDonor = Object.assign(
        defaultBody["opg.poas.lpastore"].donor,
        body["opg.poas.lpastore"].donor ?? null,
      );
      opgPoasLpastoreAttorneys = Object.assign(
        defaultBody["opg.poas.lpastore"].attorneys,
        body["opg.poas.lpastore"].attorneys ?? null,
      );
      opgPoasLpastoreCertificateProvider = Object.assign(
        defaultBody["opg.poas.lpastore"].certificateProvider,
        body["opg.poas.lpastore"].certificateProvider ?? null,
      );
    }
  }

  var updatedBody = {
    uid: uid,
    "opg.poas.sirius": opgPoasSirius,
    "opg.poas.lpastore": opgPoasLpastore,
  };

  updatedBody["opg.poas.sirius"].donor = opgPoasSiriusDonor;
  updatedBody["opg.poas.sirius"].application = opgPoasSiriusApplication;
  updatedBody["opg.poas.lpastore"].donor = opgPoasLpastoreDonor;
  updatedBody["opg.poas.lpastore"].attorneys = opgPoasLpastoreAttorneys;
  updatedBody["opg.poas.lpastore"].certificateProvider =
    opgPoasLpastoreCertificateProvider;

  await addMock(`/lpa-api/v1/digital-lpas/${uid}`, "GET", {
    status: 200,
    body: updatedBody,
  });
}

const progressIndicators = {
  async feesInProgress(digitalLpaUid) {
    await addMock(
      `/lpa-api/v1/digital-lpas/${digitalLpaUid}/progress-indicators`,
      "GET",
      {
        status: 200,
        body: {
          digitalLpaUid: digitalLpaUid,
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
  },
};

export { get, progressIndicators };
