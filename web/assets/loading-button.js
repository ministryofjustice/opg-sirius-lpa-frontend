const loadingButton = () => {
  /** @type HTMLAnchorElement|null loadingButton */
  const loadingButton = document.querySelector(
    '[data-module="app-loading-button"]'
  );

  if (loadingButton) {
    loadingButton.addEventListener(
      "click",
      (e) => {
        if (loadingButton.classList.contains("govuk-button--disabled")) {
          e.preventDefault();
          return false;
        }

        loadingButton.ariaDisabled = "true";
        loadingButton.classList.add("govuk-button--disabled");

        const messageSelector =
          loadingButton.getAttribute("data-loading-button-message") ?? "";
        const message = document.querySelector(messageSelector);

        if (message !== null) {
          message.classList.remove("govuk-!-display-none");
        }
      },
      false
    );
  }
};

export default loadingButton;
