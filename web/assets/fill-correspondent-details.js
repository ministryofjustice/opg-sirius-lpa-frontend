export default function fillCorrespondentDetails(scope) {
  scope = scope || document;

  const radios = scope.querySelectorAll(
    '[data-module="app-fill-correspondent-details"] input[type="radio"]',
  );

  if (!radios.length) return;

  radios.forEach((radio) => {
    radio.addEventListener("change", handleSubmitOnRadioChange);
  });

  function handleSubmitOnRadioChange(event) {
    const data = JSON.parse(event.target.dataset.fields);
    for (const [key, value] of Object.entries(data)) {
      const field = event.target.form.elements.namedItem(key);
      field && (field.value = value);
    }
  }
}
