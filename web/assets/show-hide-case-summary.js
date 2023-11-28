export default function showHideCaseSummary(prefix) {
    const button = document.querySelector('[data-module="app-case-summary-toggle"]',);
    const buttonText = document.querySelector(".case-summary-toggle-text")
    const buttonIcon = document.querySelector(".case-summary-toggle-icon")
    const caseSummary = document.querySelector(".case-summary");

    if (button && caseSummary) {
        button.onclick = () => {
            if (buttonText.innerText === "Hide") {
                buttonText.innerText = "Show";
                buttonIcon.classList.replace("govuk-accordion-nav__chevron--up", "govuk-accordion-nav__chevron--down");
                caseSummary.hidden = true;
            } else {
                buttonText.innerText = "Hide";
                buttonIcon.classList.replace("govuk-accordion-nav__chevron--down", "govuk-accordion-nav__chevron--up");
                caseSummary.hidden = false;
            }
        };
    }
}
