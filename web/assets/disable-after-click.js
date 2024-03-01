const disableAfterClick = () => {
  document
    .querySelectorAll('button[type="submit"][data-disable-after-click="true"]')
    .forEach((button) => {
      if (!button.form) {
        return false;
      }

      button.form.addEventListener(
        "submit",
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
