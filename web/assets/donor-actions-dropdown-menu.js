import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";

export default function donorDropDownMenu() {
  let $container = $('[data-module="donor-actions-button-menu"]');
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
    .querySelectorAll("[data-module='donor-dropdown-menu']")
    .forEach((e) => {
      let donorDropDownMenuToggle = e.querySelector(
        "[data-role='donor-dropdown-menu-toggle']",
      );

      if (donorDropDownMenuToggle === null) {
        return;
      }

      donorDropDownMenuToggle.addEventListener("change", (e) => {
        const url = encodeURI(e.target.value);
        if (url.startsWith("//") || !url.startsWith("/")) {
          return;
        }

        window.location.href = url;
      });
    });
}
