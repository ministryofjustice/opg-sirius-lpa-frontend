import autoCheckSingleAttorneyApplicant from "./auto-check-single-attorney-applicant.js";

describe("auto-check the applicant checkbox when there is only one attorney", () => {
  let attorneyRadio;

  const buildMarkup = (checkboxesHtml) => {
    document.body.innerHTML = `
    <div data-module="applicant-attorney-radio">
      <input type="radio" name="applicantType" value="attorney" />
    </div>
    <div data-module="applicant-attorney-checkboxes">
      ${checkboxesHtml}
    </div>
  `;

    attorneyRadio = document.querySelector(
      '[data-module="applicant-attorney-radio"] input',
    );
  };

  afterEach(() => {
    attorneyRadio = null;
    document.body.innerHTML = "";
  });

  test("should not check the checkbox if the attorney radio is not selected", () => {
    buildMarkup('<input type="checkbox" value="1" />');
    autoCheckSingleAttorneyApplicant();

    const checkbox = document.querySelector('input[type="checkbox"]');
    expect(checkbox.checked).toBe(false);
  });

  test("should check the single checkbox if the attorney radio is already selected on load", () => {
    buildMarkup('<input type="checkbox" value="1" />');
    attorneyRadio.checked = true;
    autoCheckSingleAttorneyApplicant();

    const checkbox = document.querySelector('input[type="checkbox"]');
    expect(checkbox.checked).toBe(true);
  });

  test("should check the single checkbox when the attorney radio is selected", () => {
    buildMarkup('<input type="checkbox" value="1" />');
    autoCheckSingleAttorneyApplicant();

    attorneyRadio.checked = true;
    attorneyRadio.dispatchEvent(new Event("change"));

    const checkbox = document.querySelector('input[type="checkbox"]');
    expect(checkbox.checked).toBe(true);
  });

  test("should not check any checkboxes when there is more than one attorney", () => {
    buildMarkup(
      '<input type="checkbox" value="1" /><input type="checkbox" value="2" />',
    );
    autoCheckSingleAttorneyApplicant();

    attorneyRadio.checked = true;
    attorneyRadio.dispatchEvent(new Event("change"));

    const checkboxes = document.querySelectorAll('input[type="checkbox"]');
    checkboxes.forEach((checkbox) => {
      expect(checkbox.checked).toBe(false);
    });
  });

  test("should do nothing if there are no attorneys to select", () => {
    buildMarkup("");
    autoCheckSingleAttorneyApplicant();

    attorneyRadio.checked = true;
    attorneyRadio.dispatchEvent(new Event("change"));

    expect(document.querySelectorAll('input[type="checkbox"]').length).toBe(0);
  });
});
