export default function showHideTrustCorpActiveRadios(scope) {
  scope = scope || document;

  const radios = scope.querySelectorAll(
    '[data-module="trust-corp-replacement-input"] input',
  );
  const activeSection = scope.querySelector(
    '[data-module="trust-corp-active-input"]',
  );

  if (!radios.length || !activeSection) return;

  activeSection.hidden = false;

  radios.forEach((radio) =>
    radio.addEventListener("change", handleShowHideActiveInputs),
  );

  function handleShowHideActiveInputs() {
    const checkedRadios = Array.from(radios).filter((radio) => radio.checked);
    if (checkedRadios.length === 0) return;
    const value = checkedRadios[0].value;

    activeSection.hidden = value === "true";
  }

  handleShowHideActiveInputs();
}
