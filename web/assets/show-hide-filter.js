export default function showHideFilter(scope) {
  scope = scope || document;
  const button = scope.querySelector("button[data-filter-toggle]");
  const filters = scope.querySelector(".moj-filter-layout__filter");
  const timeline = scope.querySelector("div[data-filter-timeline]");
  const summary = scope.querySelector("div[data-filter-summary]");

  if (button && filters) {
    button.onclick = () => {
      if (button.innerText === "Hide filters") {
        button.innerText = "Show filters";
        filters.classList.add("govuk-!-display-none");
        if (timeline) {
          timeline.classList.remove("govuk-grid-column-three-quarters");
          timeline.classList.add("govuk-grid-column-full");
        }
        if (summary) {
          summary.classList.remove("govuk-!-display-none");
        }
      } else {
        button.innerText = "Hide filters";
        filters.classList.remove("govuk-!-display-none");
        if (timeline) {
          timeline.classList.add("govuk-grid-column-three-quarters");
          timeline.classList.remove("govuk-grid-column-full");
        }
        if (summary) {
          summary.classList.add("govuk-!-display-none");
        }
      }
    };
  }
}
