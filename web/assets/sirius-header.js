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

const siriusHeaderDropdownId = "panel-tabs__dropdown";
const siriusHeaderDropdownSelector = `#${siriusHeaderDropdownId}`;
const siriusHeaderDropdownButtonSelector = `button[aria-controls="${siriusHeaderDropdownId}"]`;

const getSiriusHeaderButtons = () =>
  Array.from(document.querySelectorAll(siriusHeaderDropdownButtonSelector));

const getSiriusHeaderButtonFromElement = (element) => {
  if (!(element instanceof Element)) {
    return null;
  }

  const button = element.closest(siriusHeaderDropdownButtonSelector);

  return button instanceof HTMLButtonElement ? button : null;
};

const closeSiriusHeaderDropdown = () => {
  const dropdown = document.querySelector(siriusHeaderDropdownSelector);
  if (!(dropdown instanceof HTMLElement)) {
    return;
  }

  dropdown.replaceChildren();
  delete dropdown.dataset.openedBy;
  getSiriusHeaderButtons().forEach((button) => {
    button.setAttribute("aria-expanded", "false");
  });
};

const initSiriusHeaderDropdownToggle = () => {
  document.body.addEventListener("htmx:beforeRequest", (event) => {
    const detail = event.detail;

    const dropdown = event.detail?.target;
    if (!(dropdown instanceof HTMLElement) || dropdown.id !== siriusHeaderDropdownId) {
      return;
    }

    const requestSource = detail?.requestConfig?.elt ?? detail?.elt;
    const button = getSiriusHeaderButtonFromElement(requestSource);
    if (!button) {
      return;
    }

    const isOpen =
      dropdown.innerHTML.trim() !== "" && dropdown.dataset.openedBy === button.id;
    if (!isOpen) {
      return;
    }

    event.preventDefault();
    closeSiriusHeaderDropdown();
  });

  document.addEventListener("keydown", (event) => {
    if (event.key !== "Escape") {
      return;
    }

    const dropdown = document.querySelector(siriusHeaderDropdownSelector);
    if (!(dropdown instanceof HTMLElement) || dropdown.innerHTML.trim() === "") {
      return;
    }

    closeSiriusHeaderDropdown();
  });

  document.body.addEventListener("htmx:afterSwap", (event) => {
    const detail = event.detail;
    const dropdown = event.detail?.target;
    if (!(dropdown instanceof HTMLElement) || dropdown.id !== siriusHeaderDropdownId) {
      return;
    }

    const requestSource = detail?.requestConfig?.elt ?? detail.elt;
    const sourceButton = getSiriusHeaderButtonFromElement(requestSource);
    getSiriusHeaderButtons().forEach((button) => {
      const isSource = sourceButton && button.id === sourceButton.id;
      button.setAttribute("aria-expanded", isSource ? "true" : "false");
    });

    if (sourceButton?.id) {
      detail.target.dataset.openedBy = sourceButton.id;
      // Label the region by the button that opened it
      detail.target.setAttribute("aria-labelledby", sourceButton.id);
    }
  });
};

export default function initSiriusHeader() {
  moveSearchIntoSiriusHeader();
  initSiriusHeaderDropdownToggle();
}