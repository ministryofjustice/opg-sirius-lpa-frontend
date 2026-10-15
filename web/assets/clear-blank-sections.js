export default function clearBlankSections(scope) {
  scope = scope || document;

  const checkbox = scope.querySelector(
    '[data-module="app-clear-blank-sections"] input[type="checkbox"]',
  );
  const radios = scope.querySelectorAll(
    '[data-module="app-clear-blank-sections"] input[type="radio"][name="blankSections"]',
  );
  const conditionalSections = scope.querySelectorAll(
    '[data-module="app-clear-blank-sections"] .govuk-radios__conditional',
  );

  if (!checkbox || !radios.length) return;

  checkbox.addEventListener("change", handleClearBlankSections);

  function handleClearBlankSections() {
    if (checkbox.checked) return;

    radios.forEach((radio) => {
      radio.checked = false;
      if (radio.ariaExpanded) {
        radio.setAttribute("aria-expanded", "false");
      }
    });

    conditionalSections.forEach((section) => {
      section.classList.add("govuk-radios__conditional--hidden");
    });
  }
}
