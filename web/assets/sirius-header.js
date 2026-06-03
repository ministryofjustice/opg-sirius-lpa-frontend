/**
 * Header module
 * Handles search relocation
 */

const moveSearchIntoSiriusHeader = () => {
  const slot = document.querySelector("[data-header-search-slot]");
  if (!(slot instanceof HTMLElement)) {
    return;
  }

  const desktopSearch = document.querySelector(".app-search-inline-phase-banner");
  if (!(desktopSearch instanceof HTMLElement)) {
    return;
  }

  if (desktopSearch.parentElement === slot) {
    return;
  }

  slot.appendChild(desktopSearch);
};

export default function initSiriusHeader() {
  moveSearchIntoSiriusHeader();
}
