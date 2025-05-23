import * as digitalLpas from "../mocks/digitalLpas";

describe("Add a payment to a non-digital LPA", () => {
  beforeEach(() => {
    cy.visit("/add-payment?id=800");
  });

  it("adds a payment to the case", () => {
    cy.contains("Add a payment");
    cy.contains("7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");
    cy.get("#f-amount").type("41.00");
    cy.get("#f-source").select("PHONE");
    cy.get("#f-paymentDate").type("2022-04-25");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
  });

  it("sets the payment date to today", () => {
    cy.clock(Date.UTC(2022, 1, 25), ["Date"]); // months in Date starts from 0 so February = 1
    cy.contains("Add a payment");
    cy.contains("7000-0000-0000");
    cy.get(".moj-alert").should("not.exist");
    cy.get('[data-module="select-todays-date"]').click();
    cy.get("#f-paymentDate").should("have.value", "2022-02-25");
  });
});

describe("Add a payment to a digital LPA", () => {
  it("adds a payment to the case", () => {
    cy.wrap(digitalLpas.objections.empty("M-9999-4567-AAAA"));
    cy.visit("/add-payment?id=900");
    cy.contains("Add a payment");
    cy.contains("M-9999-4567-AAAA");
    cy.get(".moj-alert").should("not.exist");
    cy.get("#f-amount").type("82.00");
    cy.get("#f-source").select("PHONE");
    cy.get("#f-paymentDate").type("2023-08-31");
    cy.get("button[type=submit]").click();
    cy.get(".moj-alert").should("exist");
    cy.get(".moj-alert").contains("Payment added");
    cy.url().should("include", "/lpa/M-9999-4567-AAAA/payments");
  });
});
