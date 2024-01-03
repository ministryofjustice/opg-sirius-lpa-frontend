const fullWidthContainer = () => {
  if (window.location.pathname.includes("/search")) {
    const mainContainers = document.querySelectorAll(".govuk-width-container");

    if (mainContainers) {
      mainContainers.forEach((container) => {
        container.style.maxWidth = "none";
        container.classList.add(
          "govuk-!-margin-left-5",
          "govuk-!-margin-right-5",
        );
      });
    }
  }
};
export default fullWidthContainer;
