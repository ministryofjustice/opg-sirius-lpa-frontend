describe("Add a payment", () => {
  beforeEach(() => {
    cy.addMock("/lpa-api/v1/cases/724", "GET", {
      status: 200,
      body: {
        uId: "7000-0000-0000",
      },
    });

    cy.visit("/add-payment?id=724");
  });

  it("adds a payment to the case", () => {
    cy.addMock("/lpa-api/v1/cases/724/payments", "POST", {
      status: 201,
      body: {
        amount: 4100,
        paymentDate: "25/04/2022",
        source: "PHONE",
      },
    });

    cy.addMock("/lpa-api/v1/cases/724/payments", "GET", {
      status: 200,
      body: [
        {
          amount: 4100,
          case: {
            id: 800,
          },
          id: 2,
          paymentDate: "23/01/2022",
          source: "MAKE",
        },
      ],
    });

    cy.contains("Add a payment");
    cy.contains("7000-0000-0000");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-amount").type("41.00");
    cy.get("#f-source").select("PHONE");
    cy.get("#f-paymentDate").type("2022-04-25");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });

  it("sets the payment date to today", () => {
    cy.clock(Date.UTC(2022, 1, 25), ["Date"]); // months in Date starts from 0 so February = 1
    cy.contains("Add a payment");
    cy.contains("7000-0000-0000");
    cy.get(".moj-banner").should("not.exist");
    cy.get('[data-module="select-todays-date"]').click();
    cy.get("#f-paymentDate").should("have.value", "2022-02-25");
  });
});
