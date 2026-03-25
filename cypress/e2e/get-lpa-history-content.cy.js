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

  it("can view incoming document event", () => {
    mockEventHistory({
      sourceType: "IncomingDocument",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\Document\IncomingDocument`,
        friendlyDescription: "Incoming Document",
        subType: "Application related",
      },
      sourceDocument: {
        UUID: "123e4567-e89b-12d3-a456-426614174000",
        friendlyDescription: "Incoming document",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Incoming Document")
      .should("contain.text", "Application related")
      .find("a")
      .should(
        "have.attr",
        "href",
        "/lpa#/donor/1/documents?docUuid=123e4567-e89b-12d3-a456-426614174000",
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

  it("can view warning created event", () => {
    mockEventHistory({
      sourceType: "Warning",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\Warning\Warning`,
        warningType: "Complaint Received",
        warningText: "Test",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Complaint Received")
      .should("contain.text", "Test");
  });

  it("can view warning deleted event", () => {
    mockEventHistory({
      sourceType: "Warning",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\Warning\Warning`,
        warningType: "Complaint Received",
        warningText: "Test",
      },
      changeSet: {
        systemStatus: [false, true],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Complaint Received")
      .should("contain.text", "Test")
      .should("contains.text", "Warning removed by Team Less");
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

  it("can view complaint added event", () => {
    mockEventHistory({
      sourceType: "Complaint",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\PowerOfAttorney\Complaint\Complaint`,
        summary: "Test Complaint",
        severity: "Medium",
        description: "123",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Medium - Test Complaint")
      .should("contain.text", "123");
  });

  it("can view complaint updated event", () => {
    mockEventHistory({
      sourceType: "Complaint",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\PowerOfAttorney\Complaint\Complaint`,
        summary: "Test Complaint",
        severity: "Medium",
        description: "123",
      },
      changeSet: {
        receivedDate: [
          { date: "2006-12-01 00:00:00.000000" },
          { date: "2006-12-02 00:00:00.000000" },
        ],
        resolutionDate: {
          1: {
            date: "2006-12-04 00:00:00.000000",
          },
        },
        category: ["01", "02"],
        origin: { 1: "CONTACT_CENTRE" },
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should(
        "contain.text",
        "Received date: 01/12/2006 changed to: 02/12/2006",
      )
      .should("contain.text", "Resolution date: 04/12/2006")
      .should(
        "contain.text",
        "Category: Correspondence changed to: OPG Decisions",
      )
      .should("contain.text", "Origin: Contact centre");
  });

  it("can view task created event", () => {
    mockEventHistory({
      sourceType: "Task",
      type: "INS",
      entity: {
        type: "Test Type",
        name: "Some Task",
        description: "some description",
        assignee: {
          displayName: "Some User",
        },
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .within(() => {
        cy.get("p.govuk-body strong").should("have.text", "Test Type");
        cy.root()
          .should("contain.text", "Test Type now assigned to Some User")
          .should("contain.text", "Some Task — some description");
      });
  });

  it("can view task created event, no type", () => {
    mockEventHistory({
      sourceType: "Task",
      type: "INS",
      entity: {
        name: "Some Task",
        description: "some description",
        assignee: {
          displayName: "Some User",
        },
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Some Task — some description");
  });

  it("can view task completed event", () => {
    mockEventHistory({
      sourceType: "Task",
      type: "UPD",
      entity: {
        name: "Some Task",
      },
      changeSet: {
        completedDate: [
          null,
          {
            date: "2026-03-16 15:00:00.000000",
          },
        ],
        status: ["Not Started", "Completed"],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .within(() => {
        cy.get("p.govuk-body strong").should("have.text", "Completed");
        cy.root()
          .should("contain.text", "Some Task")
          .should("contain.text", "Date completed: 16/03/2026");
      });
  });

  it("can view task reassigned event", () => {
    mockEventHistory({
      sourceType: "Task",
      type: "UPD",
      entity: {
        name: "Some Task",
      },
      changeSet: {
        assignee: [
          {
            details: {
              displayName: "Some User",
            },
          },
          {
            details: {
              displayName: "Another User",
            },
          },
        ],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Some Task")
      .should(
        "contain.text",
        "Task was assigned to Some User now assigned to Another User",
      );
  });

  it("can view fee reduction approved event", () => {
    mockEventHistory({
      sourceType: "Payment",
      type: "INS",
      entity: {
        source: "FEE_REDUCTION",
        feeReductionType: "REMISSION",
        paymentDate: "2026-01-02T15:04:05+00:00",
      },
    });
    // RefData is mocked in ../mocks/common.json
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Remission approved on 02/01/2026");
  });

  it("can view fee reduction updated event", () => {
    mockEventHistory({
      sourceType: "Payment",
      type: "UPD",
      changeSet: {
        feeReductionType: ["REMISSION", "EXEMPTION"],
      },
      entity: {
        source: "FEE_REDUCTION",
      },
    });
    // RefData is mocked in ../mocks/common.json
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Reduction type: Remission changed to Exemption");
  });

  it("can view fee reduction deleted event", () => {
    mockEventHistory({
      sourceType: "Payment",
      type: "DEL",
      entity: {
        source: "FEE_REDUCTION",
        feeReductionType: "REMISSION",
        paymentDate: "2026-01-02T15:04:05+00:00",
      },
    });
    // RefData is mocked in ../mocks/common.json
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Deleted - Remission approved on 02/01/2026");
  });

  it("can view attorney added event", () => {
    mockEventHistory({
      sourceType: "Attorney",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\ReplacementAttorney`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Some User")
      .should("contain.text", "ACME");
  });

  it("can view attorney updated event with prior values", () => {
    mockEventHistory({
      sourceType: "Attorney",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\ReplacementAttorney`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
      changeSet: {
        dob: [
          { date: "2006-12-01 00:00:00.000000" },
          { date: "2006-12-02 00:00:00.000000" },
        ],
        systemStatus: [false, true],
        CorrespondenceByPhone: [true, false],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should(
        "contain.text",
        "Date of birth: 01/12/2006 changed to: 02/12/2006",
      )
      .should("contain.text", "Changed from inactive to: active")
      .should(
        "contain.text",
        "Correspondence by phone: true changed to: false",
      );
  });

  it("can view attorney updated event without prior values", () => {
    mockEventHistory({
      sourceType: "Attorney",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\ReplacementAttorney`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
      changeSet: {
        dob: { 1: { date: "2006-12-01 00:00:00.000000" } },
        systemStatus: { 1: false },
        CorrespondenceByEmail: { 1: true },
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Date of birth: 01/12/2006")
      .should("contain.text", "Changed to: inactive")
      .should("contain.text", "Correspondence by email: true");
  });

  it("can view replacement attorney added event", () => {
    mockEventHistory({
      sourceType: "ReplacementAttorney",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\ReplacementAttorney`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Some User")
      .should("contain.text", "ACME");
  });

  it("can view replacement attorney updated event with prior values", () => {
    mockEventHistory({
      sourceType: "ReplacementAttorney",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\ReplacementAttorney`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
      changeSet: {
        dob: [
          { date: "2006-12-01 00:00:00.000000" },
          { date: "2006-12-02 00:00:00.000000" },
        ],
        systemStatus: [false, true],
        CorrespondenceByPhone: [true, false],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should(
        "contain.text",
        "Date of birth: 01/12/2006 changed to: 02/12/2006",
      )
      .should("contain.text", "Changed from inactive to: active")
      .should(
        "contain.text",
        "Correspondence by phone: true changed to: false",
      );
  });

  it("can view replacement attorney updated event without prior values", () => {
    mockEventHistory({
      sourceType: "ReplacementAttorney",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\ReplacementAttorney`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
      changeSet: {
        dob: { 1: { date: "2006-12-01 00:00:00.000000" } },
        systemStatus: { 1: false },
        CorrespondenceByEmail: { 1: true },
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Date of birth: 01/12/2006")
      .should("contain.text", "Changed to: inactive")
      .should("contain.text", "Correspondence by email: true");
  });

  it("can view correspondent added event", () => {
    mockEventHistory({
      sourceType: "Correspondent",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\Correspondent`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Some User")
      .should("contain.text", "ACME");
  });

  it("can view correspondent updated event with prior values", () => {
    mockEventHistory({
      sourceType: "Correspondent",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\Correspondent`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
      changeSet: {
        dob: [
          { date: "2006-12-01 00:00:00.000000" },
          { date: "2006-12-02 00:00:00.000000" },
        ],
        systemStatus: [false, true],
        CorrespondenceByPhone: [true, false],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should(
        "contain.text",
        "Date of birth: 01/12/2006 changed to: 02/12/2006",
      )
      .should("contain.text", "Changed from inactive to: active")
      .should(
        "contain.text",
        "Correspondence by phone: true changed to: false",
      );
  });

  it("can view correspondent updated event without prior values", () => {
    mockEventHistory({
      sourceType: "Correspondent",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\Correspondent`,
        firstname: "Some",
        surname: "User",
        companyName: "ACME",
      },
      changeSet: {
        dob: { 1: { date: "2006-12-01 00:00:00.000000" } },
        systemStatus: { 1: false },
        CorrespondenceByEmail: { 1: true },
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Date of birth: 01/12/2006")
      .should("contain.text", "Changed to: inactive")
      .should("contain.text", "Correspondence by email: true");
  });

  it("can view trust corporation added event", () => {
    mockEventHistory({
      sourceType: "TrustCorporation",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\TrustCorporation`,
        companyName: "ACME",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item").eq(0).should("contain.text", "ACME");
  });

  it("can view trust corporation updated event", () => {
    mockEventHistory({
      sourceType: "TrustCorporation",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\CaseActor\TrustCorporation`,
        companyName: "ACME",
      },
      changeSet: {
        dob: [
          { date: "2006-12-01 00:00:00.000000" },
          { date: "2006-12-02 00:00:00.000000" },
        ],
        systemStatus: [false, true],
        CorrespondenceByPhone: [true, false],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "ACME")
      .should(
        "contain.text",
        "Date of birth: 01/12/2006 changed to: 02/12/2006",
      )
      .should("contain.text", "Changed from inactive to: active")
      .should(
        "contain.text",
        "Correspondence by phone: true changed to: false",
      );
  });

  it("can view investigation created event", () => {
    mockEventHistory({
      sourceType: "Investigation",
      type: "INS",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\Investigation\Investigation`,
        type: "Aspect",
        investigationTitle: "Test Investigation",
        additionalInformation: "Some extra info",
        investigationReceivedDate: "2026-03-01T00:00:00+00:00",
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Aspect - Test Investigation")
      .should("contain.text", "Some extra info")
      .should("contain.text", "Investigation received on 01/03/2026");
  });

  it("can view investigation updated event with prior values", () => {
    mockEventHistory({
      sourceType: "Investigation",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\Investigation\Investigation`,
        investigationTitle: "Test Investigation",
        additionalInformation: "Some extra info",
      },
      changeSet: {
        riskAssessmentDate: [
          { date: "2026-01-01 00:00:00.000000" },
          { date: "2026-02-01 00:00:00.000000" },
        ],
        reportApprovalOutcome: ["Outcome A", "Outcome B"],
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should(
        "contain.text",
        "Risk assessment date: 01/01/2026 changed to 01/02/2026",
      )
      .should(
        "contain.text",
        "Report approval outcome: Outcome A changed to Outcome B",
      );
  });

  it("can view investigation updated event without prior values", () => {
    mockEventHistory({
      sourceType: "Investigation",
      type: "UPD",
      entity: {
        _class: String.raw`Opg\Core\Model\Entity\Investigation\Investigation`,
        investigationTitle: "Test Investigation",
      },
      changeSet: {
        investigationClosureDate: {
          1: { date: "2026-06-15 00:00:00.000000" },
        },
        reportApprovalOutcome: { 1: "Approved" },
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Investigation closed: 15/06/2026")
      .should("contain.text", "Report approval outcome: Approved");
  });

  it("can view a manual event", () => {
    mockEventHistory({
      sourceType: "Note",
      type: "Application processing",
      entity: {
        type: "Application processing",
        name: "Test note",
        description: "This is a test note",
        document: {
          UUID: "123e4567-e89b-12d3-a456-426614174000",
          friendlyDescription: "Test document",
        },
      },
    });
    cy.visit("/donor/1/history");
    cy.get(".moj-timeline__item")
      .eq(0)
      .should("contain.text", "Test note - This is a test note");
  });
});
