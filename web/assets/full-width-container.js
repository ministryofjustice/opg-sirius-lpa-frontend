const fullWidthContainer = () => {
  if (window.location.pathname.includes("/search")) {
    const headerContainer = document.getElementsByClassName(
      "moj-header__container",
    )[0];
    const mainContainers = document.querySelectorAll(
      ".govuk-width-container",
    );

    if (headerContainer && mainContainers) {
      headerContainer.style.maxWidth = "none";
      headerContainer.className =
        headerContainer.className +
        " govuk-!-margin-left-5 govuk-!-margin-right-5";
      mainContainers.forEach((container) => container.className = "govuk-!-margin-left-5 govuk-!-margin-right-5");
    }
  }
};
export default fullWidthContainer;
