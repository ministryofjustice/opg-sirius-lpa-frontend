import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import * as GOVUKFrontend from "govuk-frontend";
import $ from "jquery";
import autoResizeTextArea from "./auto-resize-text-area.js";
import loadingButton from "./loading-button.js";
import select from "./select.js";
import todaysDate from "./todays-date.js";

// Expose jQuery on window so MOJFrontend can use it
window.$ = $;

// we aren't using the JS tabs, but they try to initialise this will stop them breaking
GOVUKFrontend.Tabs.prototype.setup = () => {};

const prefix = document.body.getAttribute("data-prefix");

GOVUKFrontend.initAll();
MOJFrontend.initAll();
autoResizeTextArea();
select(prefix);
todaysDate();
loadingButton();

if (window.self !== window.parent) {
  const success = document.querySelector(".moj-banner--success");
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
}
