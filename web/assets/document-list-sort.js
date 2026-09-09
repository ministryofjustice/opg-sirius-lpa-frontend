export default function documentListSort() {
  const select = document.querySelector("[data-document-sort]");
  if (!select) return;

  select.addEventListener("change", function () {
    const value = this.value;
    if (!value) return;

    const table = document.querySelector('[data-module="moj-sortable-table"]');
    if (!table) return;

    const tbody = table.querySelector(".govuk-table__body");
    if (!tbody) return;

    // Reset all checkbox selections
    table.querySelectorAll('input[type="checkbox"]').forEach((cb) => {
      cb.checked = false;
    });

    // Reset column header sort indicators
    table.querySelectorAll("th[aria-sort]").forEach((th) => {
      th.setAttribute("aria-sort", "none");
    });

    const rows = Array.from(tbody.querySelectorAll("tr.govuk-table__row"));

    rows.sort((a, b) => {
      switch (value) {
        case "date-new-old":
          return getSortValue(b, "date").localeCompare(getSortValue(a, "date"));
        case "date-old-new":
          return getSortValue(a, "date").localeCompare(getSortValue(b, "date"));
        case "direction-incoming":
          // "incoming" < "outgoing" alphabetically, so ascending puts Incoming first
          return getDirection(a).localeCompare(getDirection(b));
        case "direction-outgoing":
          return getDirection(b).localeCompare(getDirection(a));
        case "type-az":
          return getSortValue(a, "type").localeCompare(getSortValue(b, "type"));
        case "type-za":
          return getSortValue(b, "type").localeCompare(getSortValue(a, "type"));
        case "name-az":
          return getSortValue(a, "name").localeCompare(getSortValue(b, "name"));
        case "name-za":
          return getSortValue(b, "name").localeCompare(getSortValue(a, "name"));
        case "case-low-high":
          return normaliseUID(getSortValue(a, "case")).localeCompare(
            normaliseUID(getSortValue(b, "case")),
          );
        case "case-high-low":
          return normaliseUID(getSortValue(b, "case")).localeCompare(
            normaliseUID(getSortValue(a, "case")),
          );
        default:
          return 0;
      }
    });

    rows.forEach((row) => tbody.appendChild(row));
  });
}

function getSortValue(row, type) {
  const cell = row.querySelector(`[data-sort-type="${type}"]`);
  if (!cell) return "";
  return (cell.dataset.sortValue || cell.textContent).trim().toLowerCase();
}

function getDirection(row) {
  return (row.dataset.direction || "").toLowerCase();
}

function normaliseUID(uid) {
  return uid.replace(/-/g, "");
}
