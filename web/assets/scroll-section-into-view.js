export default function scrollSectionIntoView(scope) {
  scope = scope || document;

  const scrollTargets = scope.querySelectorAll(
    '[data-module="govuk-accordion"].govuk-accordion--scroll-sections .govuk-accordion__section-header',
  );
  Array.from(scrollTargets).forEach((scrollTarget) => {
    scrollTarget.addEventListener("click", () => {
      if (
        !scrollTarget
          .closest(".govuk-accordion__section")
          .classList.contains("govuk-accordion__section--expanded")
      ) {
        return;
      }

      scrollTarget.scrollIntoView({ behavior: "smooth" });
    });
  });
}
