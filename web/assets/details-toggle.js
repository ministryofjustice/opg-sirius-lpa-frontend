export default function initDetailsToggle() {
  const toggleButtons = document.querySelectorAll(".app-details-toggle");

  toggleButtons.forEach((button) => {
    const detailsId = button.getAttribute("data-details-id");
    const contentElement = document.getElementById(detailsId);
    const isExpanded = button.getAttribute("aria-expanded") !== "true";
    contentElement.style.display = isExpanded ? "none" : "block";

    button.setAttribute("aria-expanded", (!isExpanded).toString());
    button.addEventListener("click", function (event) {
      event.preventDefault();
      const detailsId = button.getAttribute("data-details-id");
      const contentElement = document.getElementById(detailsId);
      const isExpanded = button.getAttribute("aria-expanded") === "true";

      contentElement.style.display = isExpanded ? "none" : "block";

      button.setAttribute("aria-expanded", (!isExpanded).toString());
    });
  });
}
