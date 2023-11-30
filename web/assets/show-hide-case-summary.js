export default function showHideCaseSummary(prefix) {
  const button = document.querySelector(
    '[data-module="app-case-summary-toggle"]',
  );

  if (button) {
    const container = button.closest(
      '[data-id="caseworker-summary-container"]',
    );
    const buttonText = container.querySelector(
      '[data-id="case-summary-toggle-text"]',
    );
    const buttonIcon = container.querySelector(
      '[data-id="case-summary-toggle-icon"]',
    );
    const caseSummary = container.querySelector('[data-id="case-summary"]');

    if (container && buttonText && buttonIcon && caseSummary) {
      container.onclick = () => {
        if (buttonText.innerText === "Hide") {
          buttonText.innerText = "Show";
          buttonIcon.classList.replace(
            "govuk-accordion-nav__chevron--up",
            "govuk-accordion-nav__chevron--down",
          );
          caseSummary.hidden = true;
          caseSummary.ariaExpanded = "false";
        } else {
          buttonText.innerText = "Hide";
          buttonIcon.classList.replace(
            "govuk-accordion-nav__chevron--down",
            "govuk-accordion-nav__chevron--up",
          );
          caseSummary.hidden = false;
          caseSummary.ariaExpanded = "true";
        }
      };
    }
  }
}
