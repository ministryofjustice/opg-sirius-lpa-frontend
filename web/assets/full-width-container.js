const fullWidthContainer = () => {
  if (window.location.pathname.includes("/search")) {
    const headerContainer = document.getElementsByClassName(
      "moj-header__container"
    )[0];
    const mainContainer = document.getElementsByClassName(
      "govuk-width-container"
    )[0];
    if (headerContainer && mainContainer) {
      headerContainer.style.maxWidth = "none";
      headerContainer.className =
        headerContainer.className +
        " govuk-!-margin-left-5 govuk-!-margin-right-5";
      mainContainer.className = "govuk-!-margin-left-5 govuk-!-margin-right-5";
    }
  }
};
export default fullWidthContainer;
