import accessibleAutocomplete from "accessible-autocomplete";

export default function select(prefix) {
  enhanceElement(
    document.querySelector("[data-select-user]"),
    fetchUser(prefix)
  );
  enhanceElement(
    document.querySelector("[data-select-person]"),
    fetchPerson(prefix)
  );
}

function enhanceElement(element, source) {
  if (element) {
    accessibleAutocomplete.enhanceSelectElement({
      selectElement: element,
      minLength: 3,
      confirmOnBlur: false,
      source,
      templates: {
        inputValue(value) {
          return !value ? "" : value.text;
        },
        suggestion(value) {
          return value.text;
        },
      },
      onConfirm(selected) {
        element.innerHTML = `<option value="${selected.id}:${selected.text}" selected>${selected.text}</option>`;
      },
    });
  }
}

function fetchUser(prefix) {
  let controller = { abort: () => {} };
  const fetchOptions = {};

  return (query, callback) => {
    controller.abort();

    if ("AbortController" in window) {
      controller = new AbortController();
      fetchOptions.signal = controller.signal;
    }

    fetch(`${prefix}/search-users?q=${encodeURIComponent(query)}`, fetchOptions)
      .then((response) => response.json())
      .then((json) => {
        callback(
          json.map(({ id, displayName }) => ({ id, text: displayName }))
        );
      })
      .catch(() => {
        callback([]);
      });
  };
}

function fetchPerson(prefix) {
  let controller = { abort: () => {} };
  const fetchOptions = {};

  return (query, callback) => {
    controller.abort();

    if ("AbortController" in window) {
      controller = new AbortController();
      fetchOptions.signal = controller.signal;
    }

    fetch(
      `${prefix}/search-persons?q=${encodeURIComponent(query)}`,
      fetchOptions
    )
      .then((response) => response.json())
      .then((json) => {
        callback(
          json.map(({ uId, firstname, surname }) => ({
            id: uId,
            text: `${firstname} ${surname} (${uId})`,
          }))
        );
      })
      .catch(() => {
        callback([]);
      });
  };
}
