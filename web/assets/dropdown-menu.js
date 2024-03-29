import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";

export default function dropdownMenu() {
  let $container = $('[data-module="app-button-menu"]');
  if ($container.length > 0) {
    new MOJFrontend.ButtonMenu({
      container: $container,
      mq: "(min-width: 200em)",
      buttonText: "Case actions",
      buttonClasses:
        "govuk-button--secondary moj-button-menu__toggle-button--secondary",
    });
  }

  document.querySelectorAll("[data-module='dropdown-menu']").forEach((e) => {
    let dropdownMenuToggle = e.querySelector(
      "[data-role='dropdown-menu-toggle']",
    );

    if (dropdownMenuToggle === null) {
      return;
    }

    dropdownMenuToggle.addEventListener("change", (e) => {
      const url = encodeURI(e.target.value);
      if (url.startsWith("//") || !url.startsWith("/")) {
        return;
      }

      window.location.href = url;
    });
  });
}
