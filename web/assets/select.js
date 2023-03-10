import accessibleAutocomplete from "accessible-autocomplete";

export default function select(prefix) {
  enhanceUserPersonSearchElement(
    document.querySelector("[data-select-user]"),
    fetchUser(prefix)
  );
  enhanceUserPersonSearchElement(
    document.querySelector("[data-select-person]"),
    fetchPerson(prefix)
  );
  enhanceTemplateSearchElement(
    document.querySelector("[data-select-template]")
  );
}

function enhanceTemplateSearchElement(element) {
  if (element) {
    accessibleAutocomplete.enhanceSelectElement({
      selectElement: element,
      showAllValues: true,
      onConfirm: function (value) {
        // Provide default behaviour, which is normally overridden by `onConfirm`
        const requestedOption = [].filter.call(this.selectElement.options, option => (option.textContent || option.innerText) === value)[0]
        if (requestedOption) { requestedOption.selected = true }

        this.selectElement.dispatchEvent(new CustomEvent("confirm"));
      },
    });
  }
}

function enhanceUserPersonSearchElement(element, source) {
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

function fetchForAutocomplete(url, mapFunction) {
  let controller = { abort: () => {} };
  const fetchOptions = {
    headers: {
      Accept: "application/json",
    },
  };

  return (query, callback) => {
    controller.abort();

    if ("AbortController" in window) {
      controller = new AbortController();
      fetchOptions.signal = controller.signal;
    }

    fetch(url(query), fetchOptions)
      .then((response) => response.json())
      .then((json) => {
        callback(json.map(mapFunction));
      })
      .catch(() => {
        callback([]);
      });
  };
}

function fetchUser(prefix) {
  return fetchForAutocomplete(
    (query) => `${prefix}/search-users?q=${encodeURIComponent(query)}`,
    ({ id, displayName }) => ({ id, text: displayName })
  );
}

function fetchPerson(prefix) {
  return fetchForAutocomplete(
    (query) => `${prefix}/search-persons?q=${encodeURIComponent(query)}`,
    ({ uId, firstname, surname }) => ({
      id: uId,
      text: `${firstname} ${surname} (${uId})`,
    })
  );
}
