import * as cases from "../mocks/cases";
import * as digitalLpas from "../mocks/digitalLpas";

describe("View digital LPA history", () => {
  beforeEach(() => {
    const uid = "M-1234-1234-1234";
    const id = "1111";

    const events = [{
      "id":"607631d0-586a-40a5-ae5f-cc9121b6eef1",
      "uid":"M-1234-1234-1234",
      "applied":"2025-12-16T12:06:50Z",
      "author":"urn:opg:sirius:users:1",
      "type":"CORRECTION",
      "source":"lpa_store",
      "changes":[{
        "key":"/donor/firstNames",
        "old":"Anne",
        "new":"Anna"
      },{
        "key":"/donor/address/postcode",
        "old":"M1 1AB",
        "new":"M1 1AC"
      }]
    },{
      "id":"43e3b449-909a-44a8-9fc9-37c90c3aab4c",
      "uid":"M-1234-1234-1234",
      "applied":"2025-12-16T12:19:16Z",
      "author":"urn:opg:sirius:users:1",
      "type":"CORRECTION",
      "source":"lpa_store",
      "changes":[{
        "key":"/attorneys/0/lastName",
        "old":"Dent",
        "new":"Branch"
      },{
        "key":"/attorneys/0/address/line1",
        "old":"1 Oak Road",
        "new":"1 Willow Avenue"
      },{
        "key":"/attorneys/0/mobile",
        "old":"07123 456789",
        "new":"07123 456710"
      }]
    },{
      "uuid": "c4b4dff2-e2e8-4b71-af0a-9fabdf648edc",
      "owningCase": {
        "id": 111,
        "uId": "M-NN8A-XMHL-GF69",
        "caseSubtype": "personal-welfare",
        "caseType": "DIGITAL_LPA"
      },
      "user": {
        "id": 51,
        "phoneNumber": "12345678",
        "teams": [],
        "displayName": "system admin",
        "deleted": false,
        "email": "system.admin@opgtest.com"
      },
      "sourceType": "OutgoingDocument",
      "sourceDocument": {
        "id": 9,
        "uuid": "ac488874-4e74-49f3-bc02-d95741b96aa3",
        "friendlyDescription": "Anne Barlow - TEMP - Digital LPA Registration Form",
        "createdDate": "17\/12\/2025 09:27:12",
        "filename": "DLPA-form.pdf",
        "mimeType": "application\/pdf"
      },
      "type": "INS",
      "changeSet": [],
      "entity": {
        "_class": "Opg\\Core\\Model\\Entity\\Document\\OutgoingDocument",
        "correspondent": {
          "firstname": "Anne",
          "surname": "Barlow"
        },
        "direction": "OUTGOING",
        "friendlyDescription": "Anne Barlow - TEMP - Digital LPA Registration Form",
        "id": 9,
        "type": "Save",
        "subType": ""
      },
      "createdOn": "2025-12-17T09:27:13+00:00",
      "hash": "IK",
      "source": "sirius"
    },{
      "uuid": "055679b3-9874-43f1-9c8c-fd0851c0aa55",
      "owningCase": {
        "id": 111,
        "uId": "M-NN8A-XMHL-GF69",
        "caseSubtype": "personal-welfare",
        "caseType": "DIGITAL_LPA"
      },
      "user": {
        "id": 51,
        "phoneNumber": "12345678",
        "teams": [],
        "displayName": "system admin",
        "deleted": false,
        "email": "system.admin@opgtest.com"
      },
      "sourceType": "Task",
      "sourceTask": {
        "id": 251,
        "assignee": {
          "id": 0,
          "phoneNumber": "03004560300",
          "displayName": "Unassigned",
          "deleted": false
        }
      },
      "type": "INS",
      "changeSet": [],
      "entity": {
        "_class": "Opg\\Core\\Model\\Entity\\Task\\Task",
        "assignee": {
          "displayName": "Unassigned"
        },
        "name": "Print and post donor form"
      },
      "createdOn": "2025-12-17T09:27:13+00:00",
      "hash": "II",
      "source": "sirius"
    },{
      "uuid": "458ef72a-1d43-4c69-8c6c-dac2d138d1d3",
      "owningCase": {
        "id": 111,
        "uId": "M-NN8A-XMHL-GF69",
        "caseSubtype": "personal-welfare",
        "caseType": "DIGITAL_LPA"
      },
      "user": {
        "id": 51,
        "phoneNumber": "12345678",
        "teams": [],
        "displayName": "system admin",
        "deleted": false,
        "email": "system.admin@opgtest.com"
      },
      "sourceType": "DigitalLpa",
      "sourceCase": {
        "id": 111,
        "uId": "M-NN8A-XMHL-GF69",
        "caseSubtype": "personal-welfare",
        "caseType": "DIGITAL_LPA"
      },
      "type": "INS",
      "changeSet": [],
      "entity": {
        "_class": "Opg\\Core\\Model\\Entity\\CaseItem\\PowerOfAttorney\\DigitalLpa"
      },
      "createdOn": "2025-12-17T09:27:13+00:00",
      "hash": "IE",
      "source": "sirius"
    }];

    const mocks = Promise.allSettled([
      digitalLpas.get(uid),
      cases.warnings.empty(id),
      cases.tasks.empty(id),
      digitalLpas.objections.empty(uid),
      digitalLpas.events.get(uid, events)
    ]);

    cy.wrap(mocks);

    cy.visit("/lpa/M-1234-1234-1234/history");
  });

  it("shows LPA store and events history", () => {
    cy.contains("LPA details: Donor");
    cy.contains("First names updated from Anne to Anna");
    cy.contains("Post code updated from M1 1AB to M1 1AC");

    cy.contains("LPA details: Attorneys");
    cy.contains("Last name updated from Dent to Branch");
    cy.contains("Address line 1 updated from 1 Oak Road to 1 Willow Avenue");
    cy.contains("Mobile updated from 07123 456789 to 07123 456710");

    cy.contains("Created: OutgoingDocument");
    cy.contains("Direction: OUTGOING");

    cy.contains("Created: Task");
    cy.contains("Name: Print and post donor form");
  });
});


