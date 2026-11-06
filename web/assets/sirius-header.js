/**
 * Header module
 * Handles search relocation, dropdown page offset and people list scroll position
 */

const DROPDOWN_OFFSET_PROPERTY = "--app-sirius-header-dropdown-height";

const moveSearchIntoSiriusHeader = () => {
  const slot = document.querySelector("[data-header-search-slot]");
  if (!(slot instanceof HTMLElement)) {
    return;
  }

  const desktopSearch = document.querySelector(
    ".app-search-inline-phase-banner",
  );
  if (!(desktopSearch instanceof HTMLElement)) {
    return;
  }

  if (desktopSearch.parentElement === slot) {
    return;
  }

  slot.appendChild(desktopSearch);
};

const offsetPageForOpenDropdown = () => {
  const header = document.querySelector(".sirius-header");
  if (!(header instanceof HTMLElement)) {
    return;
  }

  if (header.dataset.dropdownOffsetInitialised === "true") {
    return;
  }
  header.dataset.dropdownOffsetInitialised = "true";

  const root = document.documentElement;

  const update = () => {
    const openPanel = header.querySelector(
      ".sirius-header__details[open] > .sirius-header__dropdown-panel",
    );
    const height =
      openPanel instanceof HTMLElement
        ? openPanel.getBoundingClientRect().height
        : 0;
    root.style.setProperty(DROPDOWN_OFFSET_PROPERTY, `${height}px`);
  };

  const observer = new ResizeObserver(update);
  header
    .querySelectorAll(".sirius-header__dropdown-panel")
    .forEach((panel) => observer.observe(panel));

  header.addEventListener("toggle", update, true);
  update();
};

const PEOPLE_LIST_SELECTOR = ".app-person-info__people";

let peopleListScrollInitialised = false;

const preservePeopleListScroll = () => {
  if (peopleListScrollInitialised) {
    return;
  }
  peopleListScrollInitialised = true;

  document.addEventListener("htmx:beforeSwap", (event) => {
    const target = event.detail.target;
    if (!(target instanceof HTMLElement)) {
      return;
    }

    const list = target.querySelector(PEOPLE_LIST_SELECTOR);
    if (list instanceof HTMLElement) {
      target.dataset.peopleListScrollTop = String(list.scrollTop);
    }
  });

  document.addEventListener("htmx:afterSettle", (event) => {
    const target = event.detail.target;
    if (!(target instanceof HTMLElement)) {
      return;
    }

    const savedScrollTop = target.dataset.peopleListScrollTop;
    if (savedScrollTop === undefined) {
      return;
    }
    delete target.dataset.peopleListScrollTop;

    const list = target.querySelector(PEOPLE_LIST_SELECTOR);
    if (list instanceof HTMLElement) {
      list.scrollTop = Number(savedScrollTop);
    }
  });
};

export default function initSiriusHeader() {
  moveSearchIntoSiriusHeader();
  offsetPageForOpenDropdown();
  preservePeopleListScroll();
}
