export default function clearPaymentValue(scope) {
  scope = scope || document;

  const radios = scope.querySelectorAll(
    '[data-module="app-clear-payment-trigger"]',
  );
  const input = scope.querySelector('[data-module="app-clear-payment-target"]');

  if (!radios.length || !input) return;

  radios.forEach((radio) =>
    radio.addEventListener("change", () => {
      input.value = "";
    }),
  );
}
