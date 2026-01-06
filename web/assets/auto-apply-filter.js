export default function autoApplyFilter() {
  const filters = document.querySelectorAll(
    '[data-module="app-auto-apply-filter"]',
  );
  let timeout = null;

  filters.forEach((filter) => {
    filter.addEventListener("click", () => {
      if (timeout !== null) clearTimeout(timeout);

      timeout = setTimeout(() => {
        filter.closest("form").submit();
      }, 1000);
    });
  });
}
