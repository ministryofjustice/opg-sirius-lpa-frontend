import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import * as GOVUKFrontend from "govuk-frontend";
import $ from "jquery";
import autoResizeTextArea from "./auto-resize-text-area.js";
import loadingButton from "./loading-button.js";
import select from "./select.js";
import todaysDate from "./todays-date.js";
import showHideFilter from "./show-hide-filter";
import fullWidthContainer from "./full-width-container";
import searchController from "./search-controller";
import textEditor from "./text-editor";
import selectTab from "./select-tab";
import handleInsertCheckboxes from "./handle-insert-checkboxes";
import autoClick from "./auto-click";
import handleCreateDocumentButton from "./handle-create-document-button";
import insertSelector from "./insert-selector";
import addressFinder from "./address-finder";
import autoApplyFilter from "./auto-apply-filter";

// Expose jQuery on window so MOJFrontend can use it
window.$ = $;

const prefix = document.body.getAttribute("data-prefix");

GOVUKFrontend.initAll();
MOJFrontend.initAll();
autoResizeTextArea();
select(prefix);
todaysDate();
loadingButton();
searchController();
showHideFilter();
fullWidthContainer();
textEditor();
selectTab();
handleInsertCheckboxes();
autoClick();
handleCreateDocumentButton();
insertSelector();
addressFinder(prefix);
autoApplyFilter();

if (window.self !== window.parent) {
  const success = document.querySelector('[data-app-reload~="page"]');
  if (success) {
    window.parent.postMessage(
      "form-done",
      `${window.location.protocol}//${window.location.host}`
    );
  }

  document.querySelectorAll("[data-app-iframe-cancel]").forEach((el) => {
    el.addEventListener("click", (event) => {
      window.parent.postMessage(
        "form-cancel",
        `${window.location.protocol}//${window.location.host}`
      );
      event.preventDefault();
    });
  });

  const saveAndExit = document.querySelector(
    '[data-app-reload~="saveAndExit"]'
  );
  if (saveAndExit) {
    window.parent.postMessage(
      "form-cancel",
      `${window.location.protocol}//${window.location.host}`
    );
  }

  const reloadTimeline = document.querySelector(
    '[data-app-reload~="reload-timeline"]'
  );
  if (reloadTimeline) {
    window.parent.postMessage(
      "reload-timeline",
      `${window.location.protocol}//${window.location.host}`
    );
  }
}
