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
}
