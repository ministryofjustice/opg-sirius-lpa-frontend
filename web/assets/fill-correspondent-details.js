export default function fillCorrespondentDetails(scope) {
  scope = scope || document;

  const radios = scope.querySelectorAll(
    '[data-module="app-fill-correspondent-details"] input[type="radio"]',
  );

  const fieldsContainer = scope.querySelector(
    '[data-module="app-fill-correspondent-details--fields"]',
  );

  if (!radios.length) return;

  radios.forEach((radio) => {
    radio.addEventListener("change", handleSubmitOnRadioChange);
  });

  function handleSubmitOnRadioChange(event) {
    if (fieldsContainer?.classList.contains("govuk-!-display-none")) {
      fieldsContainer.classList.remove("govuk-!-display-none");
    }

    const form = event.target.form;
    const fields = event.target.dataset.fields;
    if (!fields) {
      form.reset();
      // Resetting the form will uncheck the selected radio button, so we need to re-check it
      event.target.checked = true;
      return;
    }

    const data = JSON.parse(fields);
    for (const [key, value] of Object.entries(data)) {
      const field = form.elements.namedItem(key);
      field && (field.value = value);
    }
  }
}
