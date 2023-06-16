export default function showHideFilter(prefix) {
  const button = document.querySelector("button[data-filter-toggle]");
  const filters = document.querySelector(".moj-filter-layout__filter");

  if (button && filters) {
    button.onclick = () => {
      if (button.innerText === "Hide filters") {
        button.innerText = "Show filters";
        filters.classList.add("govuk-!-display-none");
      } else {
        button.innerText = "Hide filters";
        filters.classList.remove("govuk-!-display-none");
      }
    };
  }
}
