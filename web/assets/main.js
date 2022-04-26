import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import accessibleAutocomplete from "accessible-autocomplete";
import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import $ from "jquery";

// Expose jQuery on window so MOJFrontend can use it
window.$ = $;

// we aren't using the JS tabs, but they try to initialise this will stop them breaking
GOVUKFrontend.Tabs.prototype.setup = () => {};

GOVUKFrontend.initAll();
MOJFrontend.initAll();

const prefix = document.body.getAttribute("data-prefix");

const selectUser = document.querySelector("[data-select-user]");
if (selectUser) {
  accessibleAutocomplete.enhanceSelectElement({
    selectElement: selectUser,
    minLength: 3,
    confirmOnBlur: false,
    source(query, callback) {
      fetch(`${prefix}/search-users?q=${encodeURIComponent(query)}`)
        .then((response) => response.json())
        .then((json) => {
          callback(
            json.map(({ id, displayName }) => ({ id, text: displayName }))
          );
        })
        .catch(() => {
          callback([]);
        });
    },
    templates: {
      inputValue(value) {
        return !value ? "" : value.text;
      },
      suggestion(value) {
        return value.text;
      },
    },
    onConfirm(selected) {
      selectUser.innerHTML = `<option value="${selected.id}:${selected.text}" selected>${selected.text}</option>`;
    },
  });
}

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
