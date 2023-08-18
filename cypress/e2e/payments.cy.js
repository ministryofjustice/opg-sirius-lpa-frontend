describe("View a payment", () => {
  describe("No payments on case", () => {
    it("displays default message when there are no payments on the case", () => {
      cy.visit("/payments/801");
      cy.contains("7000-0000-0001");
      cy.contains("There is currently no fee data available to display.");
    });

    it("displays add payment but not apply fee reduction button", () => {
      cy.addMock("/lpa-api/v1/users/current", "GET", {
        status: 200,
        body: {
          roles: ["OPG User"],
        },
      });

      cy.visit("/payments/801");
      cy.contains(".govuk-button", "Add payment");
      cy.contains(".govuk-button", "Apply fee reduction").should("not.exist");
    });

    it("displays add payment and apply fee reduction buttons", () => {
      cy.addMock("/lpa-api/v1/users/current", "GET", {
        status: 200,
        body: {
          roles: ["OPG User", "Reduced Fees User"],
        },
      });

      cy.visit("/payments/801");
      cy.contains(".govuk-button", "Add payment");
      cy.contains(".govuk-button", "Apply fee reduction");
    });
  });

  describe("Payments on case", () => {
    it("displays payment information if there is a payment on the case", () => {
      cy.visit("/payments/800");
      cy.contains("7000-0000-0000");
      cy.contains("Total paid");
      cy.contains("£41.00");
      cy.contains("Outstanding fee due");
      cy.contains("£41.00");
      cy.contains("Fee details");
      cy.contains("Payments");
      cy.get("#f-payments-tab").click();
      cy.contains("Amount");
      cy.contains("£41.00");
      cy.contains("Date of payment:");
      cy.contains("23/01/2022");
      cy.contains("Method");
      cy.get(".govuk-link").contains("Edit payment");
      cy.get(".govuk-link").contains("Delete payment");
    });

    it("displays fee reduction information", () => {
      cy.visit("/payments/802");
      cy.contains("Fee reductions");
      cy.get("#f-fee-reductions-tab").click();
      cy.contains("Outstanding fee due");
      cy.contains("£41.00");
      cy.contains("Fee reduction type");
      cy.contains("Remission");
      cy.contains("Reduction type");
      cy.contains("Remission");
      cy.contains("Date reduction approved:");
      cy.contains("24/01/2022");
      cy.contains("Evidence:");
      cy.contains("Test multiple line evidence");
      cy.get(".govuk-link").contains("Edit fee reduction");
      cy.get(".govuk-link").contains("Delete fee reduction");
      cy.get("f-apply-fee-reduction-button").should("not.exist");
    });
  });
});
