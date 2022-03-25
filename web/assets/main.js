import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import $ from "jquery";
import TomSelect from "tom-select";

document.body.className = document.body.className
  ? document.body.className + " js-enabled"
  : "js-enabled";

// Expose jQuery on window so MOJFrontend can use it
window.$ = $;

// we aren't using the JS tabs, but they try to initialise this will stop them breaking
GOVUKFrontend.Tabs.prototype.setup = () => {};

GOVUKFrontend.initAll();
MOJFrontend.initAll();

const prefix = document.body.getAttribute("data-prefix");

const selectUser = document.querySelector("[data-select-user]");
if (selectUser) {
  new TomSelect("[data-select-user]", {
    maxItems: 1,
    create: false,
    valueField: "id",
    labelField: "displayName",
    searchField: "displayName",
    load(query, callback) {
      fetch(`${prefix}/search-users?q=${encodeURIComponent(query)}`)
        .then((response) => response.json())
        .then((json) => {
          callback(
            json.map(({ id, displayName }) => ({
              id: id + ":" + displayName,
              displayName,
            }))
          );
        })
        .catch(() => {
          callback();
        });
    },
  });
}

if (window.self !== window.parent) {
  document.body.className += " app-!-embedded";

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
