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

          if (button.type == "submit") {
            button.form.submit();
          }
        },
        false,
      );
    });
};

export default disableAfterClick;
