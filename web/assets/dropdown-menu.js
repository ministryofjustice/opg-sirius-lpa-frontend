export default function dropdownMenu() {
  document.querySelectorAll("[data-module='dropdown-menu']").forEach((e) => {
    let dropdownMenuToggle = e.querySelector(
      "[data-role='dropdown-menu-toggle']",
    );

    if (dropdownMenuToggle === null) {
      return;
    }

    dropdownMenuToggle.addEventListener("change", (e) => {
      const url = encodeURI(e.target.value);
      if (url.startsWith("//") || !url.startsWith("/")) {
        return;
      }
      
      window.location.href = url
    });
  });
}
