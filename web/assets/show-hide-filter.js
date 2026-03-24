export default function showHideFilter(prefix) {
  const button = document.querySelector("button[data-filter-toggle]");
  const filters = document.querySelector(".moj-filter-layout__filter");
  const timeline = document.querySelector("div[data-filter-timeline]");
  const summary = document.querySelector("div[data-filter-summary]");

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
