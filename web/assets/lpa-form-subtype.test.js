import lpaFormSubtype from "./lpa-form-subtype.js";

describe("LPA form subtype", () => {
  let hwRadio;
  let pfaRadio;
  let hwSection;
  let pfaSection;

  beforeEach(async () => {
    document.body.innerHTML = `
    <input type="radio" data-module="app-lpa-form-subtype-input" name="caseSubtype" value="hw" />
    <input type="radio" data-module="app-lpa-form-subtype-input" name="caseSubtype" value="pfa" />
    <div data-module="app-lpa-form-subtype-hw"></div>
    <div data-module="app-lpa-form-subtype-pfa"></div>
  `;

    hwRadio = document.querySelector(
      'input[data-module="app-lpa-form-subtype-input"][value="hw"]',
    );
    pfaRadio = document.querySelector(
      'input[data-module="app-lpa-form-subtype-input"][value="pfa"]',
    );
    hwSection = document.querySelector(
      '[data-module="app-lpa-form-subtype-hw"]',
    );
    pfaSection = document.querySelector(
      '[data-module="app-lpa-form-subtype-pfa"]',
    );

    lpaFormSubtype();
  });

  afterEach(() => {
    hwRadio = null;
    pfaRadio = null;
    hwSection = null;
    pfaSection = null;
    document.body.innerHTML = "";
  });

  test("should hide both sections when no radios are checked", async () => {
    expect(hwSection.hidden).toBe(true);
    expect(pfaSection.hidden).toBe(true);
  });

  test("should show the pfa section when the pfa radio is checked", async () => {
    pfaRadio.checked = true;
    pfaRadio.dispatchEvent(new Event("change"));
    expect(hwSection.hidden).toBe(true);
    expect(pfaSection.hidden).toBe(false);
  });

  test("should show the hw section when the hw radio is checked", async () => {
    hwRadio.checked = true;
    hwRadio.dispatchEvent(new Event("change"));
    expect(hwSection.hidden).toBe(false);
    expect(pfaSection.hidden).toBe(true);
  });
});
