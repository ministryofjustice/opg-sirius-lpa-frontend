const STORAGE_KEY = "actionPanelTab";
const WIDE_PATHS = ["create-document", "edit-document"];

let bound = false;

const getStoredTab = () => {
  try {
    return window.sessionStorage.getItem(STORAGE_KEY);
  } catch {
    return null;
  }
};

const storeTab = (tab) => {
  try {
    window.sessionStorage.setItem(STORAGE_KEY, tab);
  } catch {
    return false;
  }
  return true;
};

const getTabs = () => [...document.querySelectorAll("[data-action-panel-tab]")];

const getActivePane = () =>
  document.querySelector("[data-action-panel-pane]:not([hidden])");

const syncWide = () => {
  const panel = document.querySelector(".action-panel");
  const pane = getActivePane();
  if (!panel || !pane) {
    return;
  }
  panel.classList.toggle("wide", pane.dataset.actionPanelWide === "true");
};

const selectTab = (name, { focus = false } = {}) => {
  const tabs = getTabs();
  if (tabs.length === 0) {
    return;
  }

  const target = tabs.some((tab) => tab.dataset.actionPanelTab === name)
    ? name
    : tabs[0].dataset.actionPanelTab;

  tabs.forEach((tab) => {
    const selected = tab.dataset.actionPanelTab === target;
    tab.setAttribute("aria-selected", String(selected));
    tab.tabIndex = selected ? 0 : -1;
    if (selected && focus) {
      tab.focus();
    }
  });

  document.querySelectorAll("[data-action-panel-pane]").forEach((pane) => {
    pane.hidden = pane.dataset.actionPanelPane !== target;
  });

  storeTab(target);
  syncWide();
};

const restoreTab = () => {
  const stored = getStoredTab();
  const current = getTabs().find(
    (tab) => tab.getAttribute("aria-selected") === "true",
  );
  if (stored && current && current.dataset.actionPanelTab !== stored) {
    selectTab(stored);
  }
};

const openForm = (container, requestPath) => {
  const pane = container.closest("[data-action-panel-pane]");
  if (!pane) {
    return;
  }

  const heading = pane.querySelector("[data-action-panel-heading]");
  const alreadyOpen = pane.dataset.actionPanelFormOpen === "true";

  pane.querySelector("[data-action-panel-buttons]").hidden = true;
  pane.querySelector("[data-action-panel-back]").hidden = false;
  container.hidden = false;

  if (pane.dataset.actionPanelPendingLabel) {
    heading.textContent = pane.dataset.actionPanelPendingLabel;
  }

  pane.dataset.actionPanelWide = String(
    WIDE_PATHS.some((path) => (requestPath || "").includes(path)),
  );
  pane.dataset.actionPanelFormOpen = "true";
  syncWide();

  if (!alreadyOpen) {
    heading.focus();
  }
};

const closeForm = (pane) => {
  const container = pane.querySelector("[data-action-panel-form]");
  const heading = pane.querySelector("[data-action-panel-heading]");
  const trigger = pane.dataset.actionPanelTrigger
    ? document.getElementById(pane.dataset.actionPanelTrigger)
    : null;

  container.innerHTML = "";
  container.hidden = true;
  pane.querySelector("[data-action-panel-buttons]").hidden = false;
  pane.querySelector("[data-action-panel-back]").hidden = true;
  heading.textContent = pane.dataset.actionPanelTitle;

  delete pane.dataset.actionPanelFormOpen;
  delete pane.dataset.actionPanelPendingLabel;
  delete pane.dataset.actionPanelTrigger;
  pane.dataset.actionPanelWide = "false";
  syncWide();

  if (trigger && pane.contains(trigger)) {
    trigger.focus();
  }

  if (document.activeElement !== trigger) {
    heading.focus();
  }
};

const handleClick = (event) => {
  const tab = event.target.closest("[data-action-panel-tab]");
  if (tab) {
    selectTab(tab.dataset.actionPanelTab);
    return;
  }

  const back = event.target.closest("[data-action-panel-back]");
  if (back) {
    closeForm(back.closest("[data-action-panel-pane]"));
  }
};

const handleKeydown = (event) => {
  const tab = event.target.closest("[data-action-panel-tab]");
  if (!tab) {
    return;
  }

  const tabs = getTabs();
  const index = tabs.indexOf(tab);
  const moves = {
    ArrowRight: (index + 1) % tabs.length,
    ArrowLeft: (index - 1 + tabs.length) % tabs.length,
    Home: 0,
    End: tabs.length - 1,
  };

  if (event.key in moves) {
    event.preventDefault();
    selectTab(tabs[moves[event.key]].dataset.actionPanelTab, { focus: true });
  }
};

const handleBeforeRequest = (event) => {
  const elt = event.detail.elt;
  if (!elt.closest || !elt.closest("[data-action-panel-buttons]")) {
    return;
  }

  const pane = elt.closest("[data-action-panel-pane]");
  pane.dataset.actionPanelPendingLabel = elt.dataset.actionPanelLabel || "";
  pane.dataset.actionPanelTrigger = elt.id;
};

const handleAfterSwap = (event) => {
  const target = event.detail.target;
  if (target && target.matches && target.matches("[data-action-panel-form]")) {
    openForm(
      target,
      event.detail.pathInfo && event.detail.pathInfo.requestPath,
    );
  }
};

export default function actionPanelTabs() {
  if (!bound) {
    document.addEventListener("click", handleClick);
    document.addEventListener("keydown", handleKeydown);
    document.addEventListener("htmx:beforeRequest", handleBeforeRequest);
    document.addEventListener("htmx:afterSwap", handleAfterSwap);
    document.addEventListener("htmx:afterSettle", restoreTab);
    document.addEventListener("htmx:oobAfterSwap", restoreTab);
    bound = true;
  }

  restoreTab();
}
