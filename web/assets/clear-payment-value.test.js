import clearPaymentValue from "./clear-payment-value.js";

describe("Clear payment value", () => {
  let radio;
  let input;

  beforeEach(async () => {
    document.body.innerHTML = `
    <input type="radio" data-module="app-clear-payment-trigger" name="amount" value="92.00" />
    <input type="number" data-module="app-clear-payment-target" name="amountOther" value="11.50" />
  `;

    radio = document.querySelector('[data-module="app-clear-payment-trigger"]');
    input = document.querySelector('[data-module="app-clear-payment-target"]');

    clearPaymentValue();
  });

  afterEach(() => {
    radio = null;
    input = null;
    document.body.innerHTML = "";
  });

  test("it should clear the value when the radio is checked", async () => {
    expect(input.value).toBe("11.50");
    radio.dispatchEvent(new Event("change"));
    expect(input.value).toBe("");
  });
});
