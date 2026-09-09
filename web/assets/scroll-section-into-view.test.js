import scrollSectionIntoView from "./scroll-section-into-view";

describe("Scroll accordion section into view", () => {
  let scrollTarget;

  beforeEach(() => {
    document.body.innerHTML = `
      <div class="govuk-accordion--scroll-sections" data-module="govuk-accordion">
        <div class="govuk-accordion__section">
          <div class="govuk-accordion__section-header"></div>
        </div>
      </div>
    `;

    scrollTarget = document.querySelector(
      '[data-module="govuk-accordion"].govuk-accordion--scroll-sections .govuk-accordion__section-header',
    );

    scrollTarget.scrollIntoView = jest.fn();
  });

  afterEach(() => {
    scrollTarget = null;
    document.body.innerHTML = "";
  });

  test("do nothing if the accordion section was collapsed", async () => {
    scrollSectionIntoView();
    scrollTarget.dispatchEvent(new Event("click"));
    expect(scrollTarget.scrollIntoView).not.toHaveBeenCalled();
  });

  test("scroll the section into view if the section was expanded", async () => {
    document.body
      .querySelector(".govuk-accordion__section")
      .classList.add("govuk-accordion__section--expanded");

    scrollSectionIntoView();
    scrollTarget.dispatchEvent(new Event("click"));
    expect(scrollTarget.scrollIntoView).toHaveBeenCalled();
  });
});
