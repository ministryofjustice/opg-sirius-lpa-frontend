const submitLoadingButton = () => {
  /** @type HTMLButtonElement|null submitLoadingButton */
  const submitLoadingButton = document.querySelector(
    '[data-module="app-submit-loading-button"]',
  );

  if (submitLoadingButton) {
    submitLoadingButton.addEventListener(
      "submit",
      (e) => {
        if (submitLoadingButton.hasAttribute("disabled")) {
          e.preventDefault();
          return false;
        }

        submitLoadingButton.ariaDisabled = "true";
        submitLoadingButton.setAttribute("disabled", "true");
      },
      false,
    );
  }
};

export default submitLoadingButton;
