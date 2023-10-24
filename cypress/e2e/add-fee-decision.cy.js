describe("Add a fee decision to a non-digital LPA", () => {
  beforeEach(() => {
    cy.visit("/add-fee-decision?id=801");
  });

  it("adds a fee decision to the case", () => {
    cy.contains("Record why a fee reduction will not be applied");
    cy.contains("7000-0000-0001");
    cy.get(".moj-banner").should("not.exist");
    cy.get("#f-decisionType").select("Declined exemption");
    cy.get("#f-decisionReason").type("Invalid evidence");
    cy.get("#f-decisionDate").type("2023-10-09");
    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
    cy.get(".moj-banner").contains("Fee decision added");
  });

  it("sets the applied date to today", () => {
    cy.clock(Date.UTC(2022, 1, 25), ["Date"]); // months in Date starts from 0 so February = 1
    cy.get('[data-module="select-todays-date"]').click();
    cy.get("#f-decisionDate").should("have.value", "2022-02-25");
  });
});

describe("Add a fee decision to a digital LPA", () => {
  it("adds a fee decision to the case", () => {
    cy.visit("/add-fee-decision?id=9456");
    cy.contains("M-9999-4567-AAAA");
    cy.get("#f-decisionType").select("Declined remission");
    cy.get("#f-decisionReason").type("Insufficient evidence");
    cy.get("#f-decisionDate").type("2023-10-09");
    cy.get("button[type=submit]").click();
    cy.url().should("contain", "/lpa/M-9999-4567-AAAA/payments");
  });
});
