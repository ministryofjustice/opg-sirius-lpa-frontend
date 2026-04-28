import "cypress-axe";
import "./commands";

const TerminalLog = (violations) => {
  cy.task("log", `${violations.length} accessibility violation(s) detected on this page`);

  // Temporary log for spike investigation
  violations.forEach(({ id, nodes }) => {
    nodes.forEach(node => {
      cy.task("log", `Rule: ${id} | Target: ${JSON.stringify(node.target)} | HTML: ${node.html}`);
    });
  });

  const violationData = violations.map(({ id, impact, description, nodes }) => ({
    id, impact, description, nodes: nodes.length,
  }));
  cy.task("table", violationData);
};

afterEach(() => {
  cy.resetMocks();
  cy.injectAxe();
  cy.configureAxe({
    rules: [
      { id: "region", selector: "*:not(.govuk-back-link, .moj-search__label, .moj-search__input)", },
      // region: moj-search label and input render outside a landmark via opg-sirius-search-ui

      // Remove this suppression once the upstream package is fixed or replaced.
      {
        id: "aria-allowed-attr",
        selector: "*:not(input[type='radio'][aria-expanded])",
      },
    ],
  });
  cy.checkA11y(null, null, TerminalLog, true);
});
