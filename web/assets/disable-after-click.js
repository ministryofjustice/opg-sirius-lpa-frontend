const disableAfterClick = () => {
  document
    .querySelectorAll('[data-disable-after-click="true"]')
    .forEach((button) => {
      button.addEventListener(
        "click",
        (e) => {
          if (button.hasAttribute("disabled")) {
            e.preventDefault();
            return false;
          }

          button.ariaDisabled = "true";
          button.setAttribute("disabled", "true");
        },
        false,
      );
    });
};

export default disableAfterClick;
