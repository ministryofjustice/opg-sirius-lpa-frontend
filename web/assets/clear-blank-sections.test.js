import clearBlankSections from "./clear-blank-sections.js";

describe("clear blank sections when the checkbox is unchecked", () => {
  let checkbox;
  let radio;
  let conditionalSection;

  beforeEach(async () => {
    document.body.innerHTML = `
    <div data-module="app-clear-blank-sections">
      <input type="checkbox" />
      <input type="radio" name="blankSections" aria-expanded="true" checked />
      <div class="govuk-radios__conditional"></div>
    </div>
  `;

    checkbox = document.querySelector(
      '[data-module="app-clear-blank-sections"] input[type="checkbox"]',
    );
    radio = document.querySelector(
      '[data-module="app-clear-blank-sections"] input[type="radio"][name="blankSections"]',
    );
    conditionalSection = document.querySelector(
      '[data-module="app-clear-blank-sections"] .govuk-radios__conditional',
    );

    clearBlankSections();
  });

  afterEach(() => {
    checkbox = null;
    radio = null;
    conditionalSection = null;
    document.body.innerHTML = "";
  });

  test("uncheck the radios and hide the conditional sections when the checkbox is checked", async () => {
    expect(radio.checked).toBe(true);
    expect(radio.ariaExpanded).toBe("true");
    expect(
      conditionalSection.classList.contains(
        "govuk-radios__conditional--hidden",
      ),
    ).toBe(false);
    checkbox.dispatchEvent(new Event("change"));
    expect(radio.checked).toBe(false);
    expect(radio.ariaExpanded).toBe("false");
    expect(
      conditionalSection.classList.contains(
        "govuk-radios__conditional--hidden",
      ),
    ).toBe(true);
  });
});
