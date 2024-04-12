import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";

export default function rightAlignedDropDownMenu() {
  let $container = $('[data-module="right-aligned-app-button-menu"]');
  if ($container.length > 0) {
    new MOJFrontend.ButtonMenu({
      container: $container,
      mq: "(min-width: 200em)",
      buttonText: "Donor record actions",
      buttonClasses:
        "govuk-button--secondary moj-button-menu__toggle-button--secondary",
    });
  }

  document
    .querySelectorAll("[data-module='right-dropdown-menu']")
    .forEach((e) => {
      let rightDropDownMenuToggle = e.querySelector(
        "[data-role='right-dropdown-menu-toggle']",
      );

      if (rightDropDownMenuToggle === null) {
        return;
      }

      rightDropDownMenuToggle.addEventListener("change", (e) => {
        const url = encodeURI(e.target.value);
        if (url.startsWith("//") || !url.startsWith("/")) {
          return;
        }

        window.location.href = url;
      });
    });
}
