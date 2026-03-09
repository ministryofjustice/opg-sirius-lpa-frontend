describe("Search filter", () => {
  // Helper to build a minimal mock with a given aggregations object
  const mockWithAggregations = (personType, results = []) => {
    cy.addMock("/lpa-api/v1/search/persons", "POST", {
      status: 200,
      body: {
        aggregations: { personType },
        results,
        total: { count: results.length },
      },
    });
  };

  // Minimal result record factory
  const makeResult = (personType, id = 1) => ({
    id,
    uId: `7000-0000-000${id}`,
    personType,
    firstname: "Test",
    surname: "Person",
    cases: [],
  });

  describe("Filter panel visibility", () => {
    beforeEach(() => {
      mockWithAggregations({ Donor: 1 }, [makeResult("Donor")]);
      cy.visit("/search?term=test");
    });

    it("shows the filter panel by default", () => {
      cy.get(".moj-filter-layout__filter").should("be.visible");
    });

    it("hides the filter panel when Hide filters is clicked", () => {
      cy.contains(".govuk-button", "Hide filters").click();
      cy.get(".moj-filter-layout__filter").should("not.be.visible");
    });

    it("shows the filter panel again when Show filters is clicked", () => {
      cy.contains(".govuk-button", "Hide filters").click();
      cy.contains(".govuk-button", "Show filters").click();
      cy.get(".moj-filter-layout__filter").should("be.visible");
    });
  });

  describe("Visible filter options (VEGA-3584)", () => {
    it("only shows filter options that have matching results", () => {
      mockWithAggregations({ Donor: 2, Deputy: 1 }, [
        makeResult("Donor", 1),
        makeResult("Donor", 2),
        makeResult("Deputy", 3),
      ]);
      cy.visit("/search?term=bob");

      cy.contains("label", "Donor (2)").should("exist");
      cy.contains("label", "Deputy (1)").should("exist");

      cy.contains("label", /Client/).should("not.exist");
      cy.contains("label", /Attorney/).should("not.exist");
      cy.contains("label", /Replacement attorney/).should("not.exist");
      cy.contains("label", /Trust corporation/).should("not.exist");
      cy.contains("label", /Notified person/).should("not.exist");
      cy.contains("label", /Certificate provider/).should("not.exist");
      cy.contains("label", /Correspondent/).should("not.exist");
    });

    it("shows all nine filter options when all person types have results", () => {
      mockWithAggregations({
        Donor: 1,
        Client: 1,
        Deputy: 1,
        Attorney: 1,
        "Replacement Attorney": 1,
        "Trust Corporation": 1,
        "Notified Person": 1,
        "Certificate Provider": 1,
        Correspondent: 1,
      });
      cy.visit("/search?term=test");

      cy.contains("label", /Donor/).should("exist");
      cy.contains("label", /Client/).should("exist");
      cy.contains("label", /Deputy/).should("exist");
      cy.contains("label", /Attorney \(/).should("exist");
      cy.contains("label", /Replacement attorney/).should("exist");
      cy.contains("label", /Trust corporation/).should("exist");
      cy.contains("label", /Notified person/).should("exist");
      cy.contains("label", /Certificate provider/).should("exist");
      cy.contains("label", /Correspondent/).should("exist");
    });

    it("shows no filter options when the search returns no results", () => {
      mockWithAggregations({});
      cy.visit("/search?term=zzznomatch");

      cy.contains("label", /Donor/).should("not.exist");
      cy.contains("label", /Deputy/).should("not.exist");
      cy.contains("label", /Attorney/).should("not.exist");
    });

    it("displays the correct count next to each filter option", () => {
      mockWithAggregations({
        Donor: 5,
        Attorney: 3,
        Correspondent: 1,
      });
      cy.visit("/search?term=test");

      cy.contains("label", "Donor (5)").should("exist");
      cy.contains("label", "Attorney (3)").should("exist");
      cy.contains("label", "Correspondent (1)").should("exist");
    });
  });

  describe("Active filters with zero results (VEGA-3584)", () => {
    it("keeps a selected filter visible even when its count drops to zero", () => {
      // Donor was selected, but refined search returns only Deputies
      mockWithAggregations({ Deputy: 2 }, [
        makeResult("Deputy", 1),
        makeResult("Deputy", 2),
      ]);
      cy.visit("/search?term=test&person-type=Donor");

      // Donor checkbox should still render so the user can uncheck it
      cy.get("#f-person-type-donor").should("exist").and("be.checked");
      cy.contains("label", "Donor (0)").should("exist");

      // Deputy should also appear as it has results
      cy.contains("label", "Deputy (2)").should("exist");
    });

    it("shows the selected filter tag so it can be removed", () => {
      mockWithAggregations({ Deputy: 2 });
      cy.visit("/search?term=test&person-type=Donor");

      cy.contains(".moj-filter__tag", "Donor").should("exist");
      cy.contains(".moj-filter__selected-heading", "Clear filters").should(
        "exist",
      );
    });

    it("removes a zero-result active filter when its tag is clicked", () => {
      mockWithAggregations({ Deputy: 2 });
      cy.visit("/search?term=test&person-type=Donor");

      cy.contains(".moj-filter__tag", "Donor").click();
      cy.get(".moj-filter__tag").should("not.exist");
    });
  });

  describe("Applying and removing filters", () => {
    beforeEach(() => {
      mockWithAggregations({
        Attorney: 4,
        "Trust Corporation": 2,
        Donor: 1,
      });
      cy.visit("/search?term=test");
    });

    it("hides the Apply filters button when JS is enabled (auto-apply)", () => {
      cy.contains("Apply filters").should("not.be.visible");
    });

    it("adds a filter tag when a checkbox is clicked", () => {
      cy.contains("label", /Attorney \(/).click();
      cy.contains(".moj-filter__tag", "Attorney").should("exist");
    });

    it("adds multiple filter tags when multiple checkboxes are clicked", () => {
      cy.contains("label", /Attorney \(/).click();
      cy.contains("label", /Trust corporation/).click();

      cy.contains(".moj-filter__tag", "Attorney").should("exist");
      cy.contains(".moj-filter__tag", "Trust Corporation").should("exist");
    });

    it("shows the Selected filters heading and Clear filters link when a filter is active", () => {
      cy.contains("label", /Attorney \(/).click();

      cy.get(".moj-filter__selected").should("exist");
      cy.contains(".moj-filter__selected-heading", "Clear filters").should(
        "exist",
      );
    });

    it("removes all filter tags when Clear filters is clicked", () => {
      cy.contains("label", /Attorney \(/).click();
      cy.contains("label", /Trust corporation/).click();

      cy.contains(".moj-filter__selected-heading", "Clear filters")
        .find("a")
        .click();

      cy.get(".moj-filter__tag").should("not.exist");
    });

    it("removes a specific filter when its tag is clicked", () => {
      cy.contains("label", /Attorney \(/).click();
      cy.contains("label", /Trust corporation/).click();

      cy.contains(".moj-filter__tag", "Attorney").click();

      cy.get(".moj-filter__tag").should("not.contain", "Attorney");
      cy.contains(".moj-filter__tag", "Trust Corporation").should("exist");
    });
  });
});
