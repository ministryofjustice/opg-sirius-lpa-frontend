// Base event structure with common fields
const baseEvent = {
  id: 222,
  owningCase: {
    id: 105,
    uId: "7000-9000-7000",
    caseSubtype: "pfa",
    caseType: "LPA",
  },
  user: {
    id: 6619,
    phoneNumber: "03001234567",
    teams: [],
    displayName: "Team Less",
    deleted: false,
    email: "teamless@test.uk",
  },
  changeSet: [],
  hash: "N7R",
  createdOn: "2026-01-23T16:23:29+00:00",
};

// create a custom event by merging with base event
const mockEventHistory = (eventOverides) => {
  cy.addMock("/lpa-api/v1/persons/1/events?&sort=id:desc&limit=999", "GET", {
    status: 200,
    body: {
      events: [
        {
          ...baseEvent,
          ...eventOverides,
        },
      ],
    },
  });
};

describe("Show correct event content", () => {
  it("can view payment deleted event", () => {
    mockEventHistory({
      sourceType: "Payment",
      type: "DEL",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\PowerOfAttorney\Payment\Payment`,
        amount: 2345,
        source: "CHEQUE",
        paymentDate: "2006-01-02T15:04:05+00:00",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Deleted - £23.45 paid by cheque on 02/01/2006");
  });

  it("can view payment added event", () => {
    mockEventHistory({
      sourceType: "Payment",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\PowerOfAttorney\Payment\Payment`,
        amount: 8200,
        source: "CHEQUE",
        paymentDate: "2006-01-02T15:04:05+00:00",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "£82.00 paid by cheque on 02/01/2006");
  });

  it("can view payment updated event", () => {
    mockEventHistory({
      sourceType: "Payment",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\PowerOfAttorney\Payment\Payment`,
        amount: 9200,
        source: "PHONE",
        paymentDate: "2007-01-02T15:04:05+00:00",
      },
      changeSet: {
        amount: [8200, 9200],
        source: ["CHEQUE", "PHONE"],
        paymentDate: [
          {
            date: "2006-01-02 15:04:05.000000",
            timezone_type: 3,
            timezone: "UTC",
          },
          {
            date: "2007-01-02 15:04:05.000000",
            timezone_type: 3,
            timezone: "UTC",
          },
        ],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Amount: £82.00 changed to: £92.00")
      .should(
        "contain.text",
        "Payment method: paid by cheque changed to: paid over the phone",
      )
      .should(
        "contain.text",
        "Payment date: 02/01/2006 changed to: 02/01/2007",
      );
  });

  it("can view outbound document event", () => {
    mockEventHistory({
      sourceType: "OutgoingDocument",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\PowerOfAttorney\Document\OutgoingDocument`,
        friendlyDescription: "Joe Bloggs - Letter sent to donor",
      },
      sourceDocument: {
        UUID: "123e4567-e89b-12d3-a456-426614174000",
        friendlyDescription: "Joe Bloggs - Letter sent to donor",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Joe Bloggs - Letter sent to donor")
      .find("a")
      .should(
        "have.attr",
        "href",
        "/lpa#/donor/1/documents?docUuid=123e4567-e89b-12d3-a456-426614174000",
      );
  });

  it("can view case statement updated event", () => {
    mockEventHistory({
      sourceType: "Lpa",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseItem\PowerOfAttorney\Lpa`,
      },
      changeSet: {
        status: ["Pending", "Withdrawn"],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Status changed from Pending to Withdrawn");
  });
});
