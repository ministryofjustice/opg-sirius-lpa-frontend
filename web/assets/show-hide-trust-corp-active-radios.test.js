import showHideTrustCorpActiveRadios from "./show-hide-trust-corp-active-radios.js";

describe("show or hide the trust corporation active inputs", () => {
  let attorneyRadio;
  let replacementAttorneyRadio;
  let activeInputs;

  beforeEach(async () => {
    document.body.innerHTML = `
    <div data-module="trust-corp-replacement-input">
      <input type="radio" name="isReplacement" value="true" />
      <input type="radio" name="isReplacement" value="false" />
    </div>
    <div data-module="trust-corp-active-input"></div>
  `;

    attorneyRadio = document.querySelector(
      '[data-module="trust-corp-replacement-input"] input[value="false"]',
    );
    replacementAttorneyRadio = document.querySelector(
      '[data-module="trust-corp-replacement-input"] input[value="true"]',
    );
    activeInputs = document.querySelector(
      '[data-module="trust-corp-active-input"]',
    );

    showHideTrustCorpActiveRadios();
  });

  afterEach(() => {
    attorneyRadio = null;
    replacementAttorneyRadio = null;
    activeInputs = null;
    document.body.innerHTML = "";
  });

  test("should show the active inputs section when no radios are checked", async () => {
    expect(activeInputs.hidden).toBe(false);
  });

  test("should show the active inputs section if the attorney radio is already checked", async () => {
    attorneyRadio.checked = true;
    showHideTrustCorpActiveRadios();

    expect(activeInputs.hidden).toBe(false);
  });

  test("should hide the active inputs section if the replacement attorney radio is already checked", async () => {
    replacementAttorneyRadio.checked = true;
    showHideTrustCorpActiveRadios();

    expect(activeInputs.hidden).toBe(true);
  });

  test("should show the active inputs section when the attorney radio is checked", async () => {
    attorneyRadio.checked = true;
    attorneyRadio.dispatchEvent(new Event("change"));
    expect(activeInputs.hidden).toBe(false);
  });

  test("should hide the active inputs section when the replacement attorney radio is checked", async () => {
    replacementAttorneyRadio.checked = true;
    replacementAttorneyRadio.dispatchEvent(new Event("change"));
    expect(activeInputs.hidden).toBe(true);
  });
});
