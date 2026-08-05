export default function lpaFormSubtype(scope) {
  scope = scope || document;

  const radios = scope.querySelectorAll(
    '[data-module="app-lpa-form-subtype-input"]',
  );
  const hwSection = scope.querySelector(
    '[data-module="app-lpa-form-subtype-hw"]',
  );
  const pfaSection = scope.querySelector(
    '[data-module="app-lpa-form-subtype-pfa"]',
  );

  if (!radios.length || !hwSection || !pfaSection) return;

  hwSection.hidden = true;
  pfaSection.hidden = true;

  radios.forEach((radio) =>
    radio.addEventListener("change", handleFormSubtype),
  );

  function handleFormSubtype() {
    const checkedRadios = Array.from(radios).filter((radio) => radio.checked);
    if (checkedRadios.length === 0) return;
    const value = checkedRadios[0].value;

    hwSection.hidden = value !== "hw";
    pfaSection.hidden = value !== "pfa";
  }

  handleFormSubtype();
}
