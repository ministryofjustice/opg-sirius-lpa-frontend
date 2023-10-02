export default function dropdownMenu() {
  document.querySelectorAll("[data-module='dropdown-menu']").forEach((e) => {
    let dropdownMenuToggle = e.querySelector(
      "[data-role='dropdown-menu-toggle']",
    );

    if (dropdownMenuToggle === null) {
      return;
    }

    dropdownMenuToggle.addEventListener("click", () => {
      const dropdownOpen = e.getAttribute("data-dropdown-open") === "true";
      e.setAttribute("data-dropdown-open", dropdownOpen ? "false" : "true");
    });
  });
}
