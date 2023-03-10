describe("Create donor", () => {
  beforeEach(() => {
    cy.visit("/edit-donor?id=188");
  });

  it("edits a donor", () => {
    cy.contains("Edit Donor");
    cy.get(".moj-banner").should("not.exist");

    cy.get("#f-firstname").should("have.value", "John");
    cy.get("#f-surname").should("have.value", "Doe");
    cy.get("#f-dob").should("have.value", "1970-05-05");

    cy.get("#f-salutation").type("Dr");
    cy.get("#f-firstname").clear().type("Will");
    cy.get("#f-middlenames").clear().type("Oswald");
    cy.get("#f-surname").clear().type("Niesborella");
    cy.get("#f-dob").type("1995-07-01");
    cy.get("#f-previousNames").type("Will Macphail");
    cy.get("#f-otherNames").type("Bill");
    cy.get("#f-addressLine1").type("47209 Stacey Plain");
    cy.get("#f-addressLine2").type("Suite 113");
    cy.get("#f-addressLine3").type("Devonburgh");
    cy.get("#f-town").type("Marquardtville");
    cy.get("#f-county").type("North Carolina");
    cy.get("#f-postcode").type("40936");
    cy.get("#f-country").type("United States");
    cy.get("#f-isAirmailRequired").click();
    cy.get("#f-phoneNumber").type("0841781784");
    cy.get("#f-email").type("docniesborella@mail.test");
    cy.get("#f-correspondenceBy-post").click();
    cy.get("#f-correspondenceBy-email").click();
    cy.get("#f-correspondenceBy-welsh").click();
    cy.get("#f-researchOptOut").click();

    cy.get("button[type=submit]").click();
    cy.get(".moj-banner").should("exist");
    cy.get(".moj-banner").contains("Donor was edited");
  });
});
