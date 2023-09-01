describe("Apply a fee reduction to a non-digital LPA", () => {
  beforeEach(() => {
    cy.visit("/apply-fee-reduction?id=801");
  });

  it("adds a fee reduction to the case", () => {
    cy.contains("Apply a fee reduction");
    cy.contains("7000-0000-0001");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-feeReductionType").select("Remission");
    cy.get("#f-paymentEvidence").type("Test evidence");
    cy.get("#f-paymentDate").type("2022-04-25");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
  });

  it("sets the applied date to today", () => {
    cy.clock(Date.UTC(2022, 1, 25), ["Date"]); // months in Date starts from 0 so February = 1
    cy.contains("Apply a fee reduction");
    cy.get('[data-module="select-todays-date"]').click();
    cy.get("#f-paymentDate").should("have.value", "2022-02-25");
  });
});

describe("Apply a fee reduction to a digital LPA", () => {
  it("adds a fee reduction to the case", () => {
    cy.visit("/apply-fee-reduction?id=9456");
    cy.contains("Apply a fee reduction");
    cy.contains("M-999-456-AAA");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-feeReductionType").select("Remission");
    cy.get("#f-paymentEvidence").type("Test evidence");
    cy.get("#f-paymentDate").type("2023-09-01");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
    cy.get(".moj-banner").contains("Remission approved");
    cy.url().should("contain", "/lpa/M-999-456-AAA/payments");
  });
});
