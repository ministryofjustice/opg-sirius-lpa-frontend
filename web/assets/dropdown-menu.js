export default function dropdownMenu() {
  document.querySelectorAll("[data-module='dropdown-menu']").forEach((e) => {
    let dropdownMenuToggle = e.querySelector(
      "[data-role='dropdown-menu-toggle']",
    );

    if (dropdownMenuToggle === null) {
      return;
    }

    dropdownMenuToggle.addEventListener("change", (e) => {
      window.location.href = e.target.value;
    });
  });
}
