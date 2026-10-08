describe("View a payment", () => {
    describe("No payments on case", () => {
        beforeEach(() => {
            cy.addMock("/lpa-api/v1/cases/801", "GET", {
                status: 200,
                body: {
                    id: 801,
                    uid: "7000-0000-0001",
                    caseType: "LPA",
                    status: "Registered",
                    expectedPaymentTotal: 0,
                },
            });

            cy.addMock("/lpa-api/v1/cases/801/fee-decisions", "GET", {
                status: 200,
                body: [],
            });
        })

        it("displays default message when there are no payments on the case", () => {
            cy.addMock("/lpa-api/v1/users/current", "GET", {
                status: 200,
                body: {
                    roles: ["OPG User"],
                },
            });

            cy.visit("/payments/801");
            cy.contains("7000-0000-0001");
            cy.contains("There is currently no fee data available to display.");
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
        beforeEach(() => {
            cy.addMock("/lpa-api/v1/cases/801", "GET", {
                status: 200,
                body: {
                    id: 801,
                    uid: "7000-0000-0000",
                    caseType: "LPA",
                    status: "Registered",
                    expectedPaymentTotal: 8200,
                },
            });

            cy.addMock("/lpa-api/v1/cases/801/fee-decisions", "GET", {
                status: 200,
                body: [
                    {
                        id:1,
                        decisionDate: "23/01/2022",
                        decisionType: "DECLINED_REMISSION",
                        decisionReason: "We will not let you have a fee reduction",
                    }
                ],
            });

            cy.addMock("/lpa-api/v1/cases/801/payments", "GET", {
                status: 200,
                body: [
                    {
                        id: 2,
                        amount: 4100,
                        paymentDate: "23/01/2022",

                    },
                    {
                        id: 3,
                        source: "FEE_REDUCTION",
                        feeReductionType: "REMISSION",
                        paymentEvidence: "Test multiple line evidence",
                        paymentDate: "24/01/2022",
                    }
                ],
            });
        });

        it("displays payment information if there is a payment on the case", () => {
            cy.visit("/payments/801");
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
            cy.visit("/payments/801");
            cy.contains("Fee reductions");
            cy.get("#f-fee-reductions-tab").click();
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

        it("displays decline fee reduction information", () => {
        cy.visit("/payments/801");
        cy.contains("Decline fee reductions");
        cy.get("#f-decline-fee-reductions-tab").click();
        cy.contains("Reduction declined:").should("be.visible");
        cy.contains("Remission declined");
        cy.contains("Date declined:").should("be.visible");
        cy.contains("24/01/2022");
        cy.contains("Reason:").should("be.visible");
        cy.contains("We will not let you have a fee reduction");
        });
    });

    describe("Online payments", () => {
        beforeEach(() => {
            cy.addMock("/lpa-api/v1/cases/800", "GET", {
                status: 200,
                body: {
                    id: 800,
                    uid: "7000-0000-0000",
                    caseType: "LPA",
                    status: "Registered",
                    expectedPaymentTotal: 200,
                },
            });

            cy.addMock("/lpa-api/v1/cases/800/payments", "GET", {
                status: 200,
                body: [
                    {
                        amount: 8200,
                        case: {
                            id: 800,
                        },
                        id: 234,
                        paymentDate: "12/02/2025",
                        source: "ONLINE",
                        locked: true,
                    },
                ],
            });

            cy.addMock("/lpa-api/v1/cases/800/fee-decisions", "GET", {
                status: 200,
                body: [],
            });
        });

        it("can be deleted by system admin users", () => {
            cy.addMock("/lpa-api/v1/users/current", "GET", {
                status: 200,
                body: {
                    roles: ["OPG User", "System Admin"],
                },
            });
            cy.visit("/payments/800");

            cy.contains("7000-0000-0000");
            cy.get("#f-payments-tab").click();
            cy.contains("Edit payment").should("not.exist");
            cy.get(".govuk-link").contains("Delete payment");
        });

        it("can not be deleted by other users", () => {
          cy.addMock("/lpa-api/v1/users/current", "GET", {
            status: 200,
            body: {
              roles: ["OPG User", "Reduced Fees User"],
            },
          });

          cy.visit("/payments/800");
          cy.contains("7000-0000-0000");
          cy.get("#f-payments-tab").click();
          cy.contains("Edit payment").should("not.exist");
          cy.contains("Delete payment").should("not.exist");
        });
    });
});
