import { nodeListForEach } from "./lib/nodeListForEach";

/**
 * @typedef {Object} Options
 * @property {string} prefix
 */

/**
 * @interface {} options
 * @param {HTMLElement} $module
 * @param {Options} options
 */
function AddressFinder($module, options) {
  this.$module = $module;
  this.results = [];
  this.baseUrl = options.prefix;

  this.$editMode =
    this.$module.querySelectorAll('input[value]:not([value=""])').length > 0;

  this.$inputs = this.$module.querySelectorAll(".govuk-form-group");
  const id = this.$module.id || Math.random().toString(36).substring(2);

  const $container = document.createElement("div");
  $container.innerHTML = this.$editMode
    ? AddressFinder.editModeTemplate(id)
    : AddressFinder.template(id);

  this.$module.appendChild($container);

  /** @type {HTMLInputElement} */
  this.$input = $container.querySelector("input");
  /** @type {HTMLDivElement} */
  this.$dropdownContainer = $container.querySelector("#dropdown-container");
  /** @type {HTMLSelectElement} */
  this.$dropdown = $container.querySelector("select");
  /** @type {HTMLDivElement} */
  this.$error = $container.querySelector(".govuk-error-message");
  /** @type {HTMLButtonElement} */
  this.$button = $container.querySelector("button");
  /** @type {HTMLParagraphElement} */
  this.$linkContainer = $container.querySelector(
    "#address-finder-link-container",
  );
  /** @type {HTMLDivElement} */
  this.$controls = $container.querySelector("#address-finder-controls");

  this.$input.addEventListener("keydown", this.handleKeydown.bind(this));
  this.$button.addEventListener("click", this.handleSearch.bind(this));
  this.$dropdown.addEventListener("change", this.handleSelect.bind(this));

  const $label = $container.querySelector(`[for="f-${id}-input"]`);
  $label.innerText = this.$module.getAttribute("data-app-address-finder-label");

  this.fillCountry = this.$module.getAttribute(
    "data-app-address-finder-fill-country",
  );

  const $link = $container.querySelector(".govuk-link");
  $link?.addEventListener("click", this.showControls.bind(this));

  this.$inputContainer = this.$editMode
    ? $container.querySelector("#address-finder-inputs")
    : $container.querySelector(".govuk-details__text");

  nodeListForEach(this.$inputs, ($input) => {
    this.$inputContainer.appendChild($input);
  });
}

AddressFinder.template = (id) => `
  <div class="govuk-form-group govuk-!-margin-bottom-4">
    <label class="govuk-label" for="f-${id}-input"></label>
    <div class="govuk-hint" id="f-${id}-hint">
      Enter a UK postcode
    </div>
    <p id="f-${id}-error" class="govuk-error-message govuk-!-display-none">
      <span class="govuk-visually-hidden">Error:</span>
      No matching address found. Please try again using a UK postcode or enter the address manually
    </p>
    <div class="address-finder__container">
      <input
        class="govuk-input"
        id="f-${id}-input"
        aria-describedby="f-${id}-hint"
      />
      <button class="govuk-button govuk-!-margin-bottom-0" type="button">
        Look up UK postcode
      </button>
    </div>
  </div>
  <div id="dropdown-container" class="govuk-form-group govuk-!-margin-bottom-4 govuk-!-display-none">
    <label class="govuk-label govuk-visually-hidden" for="f-${id}-select">
      Select an address
    </label>
    <select class="govuk-select govuk-!-padding-0" id="f-${id}-select"></select>
  </div>
  <details class="govuk-details govuk-!-margin-bottom-0">
    <summary class="govuk-details__summary">
      <span class="govuk-details__summary-text">
        Enter address manually
      </span>
    </summary>
    <div class="govuk-details__text">
    </div>
  </details>
`;

AddressFinder.editModeTemplate = (id) => `
  <label class="govuk-label" for="f-${id}-input"></label>
  <p id="address-finder-link-container" class="govuk-body govuk-!-margin-bottom-4 govuk-!-margin-top-4">
    <a href="#" class="govuk-link govuk-link--no-visited-state">
      Look up UK postcode
    </a>
  </p>
  <div id="address-finder-controls" class="govuk-!-display-none govuk-!-margin-bottom-6">
    <div class="govuk-form-group govuk-!-margin-bottom-4">
      <div class="govuk-hint" id="f-${id}-hint">
        Enter a UK postcode
      </div>
      <p id="f-${id}-error" class="govuk-error-message govuk-!-display-none">
        <span class="govuk-visually-hidden">Error:</span>
        No matching address found. Please try again using a UK postcode or enter the address manually
      </p>
      <div class="address-finder__container">
        <input
          class="govuk-input"
          id="f-${id}-input"
          aria-describedby="f-${id}-hint"
        />
        <button class="govuk-button govuk-!-margin-bottom-0" type="button">
          Look up UK postcode
        </button>
      </div>
    </div>
    <div id="dropdown-container" class="govuk-form-group govuk-!-margin-bottom-4 govuk-!-display-none">
      <label class="govuk-label govuk-visually-hidden" for="f-${id}-select">
        Select an address
      </label>
      <select class="govuk-select govuk-!-padding-0" id="f-${id}-select"></select>
    </div>
  </div>
  <div id="address-finder-inputs"></div>
`;

AddressFinder.prototype.showControls = function (e) {
  e.preventDefault();
  this.$controls?.classList.remove("govuk-!-display-none");
  this.$linkContainer?.classList.add("govuk-!-display-none");
};

AddressFinder.prototype.showError = function () {
  this.$input.classList.add("govuk-input--error");
  this.$input.setAttribute(
    "aria-describedby",
    this.$input.getAttribute("aria-describedby") + " " + this.$error.id,
  );
  this.$input
    .closest(".govuk-form-group")
    ?.classList.add("govuk-form-group--error");

  this.$error.classList.remove("govuk-!-display-none");
};

AddressFinder.prototype.resetError = function () {
  this.$input.classList.remove("govuk-input--error");

  const describedBy = this.$input.getAttribute("aria-describedby");
  if (describedBy) {
    this.$input.setAttribute(
      "aria-describedby",
      describedBy.replace(this.$error.id, "").trim(),
    );
  }

  this.$input
    .closest(".govuk-form-group")
    ?.classList.remove("govuk-form-group--error");

  this.$error.classList.add("govuk-!-display-none");
};

AddressFinder.prototype.handleKeydown = function (e) {
  if (e.key === "Enter") {
    e.preventDefault();
    this.handleSearch();
  }
};

AddressFinder.prototype.handleSearch = function () {
  this.$button.disabled = true;
  this.$dropdownContainer.classList.add("govuk-!-display-none");

  this.resetError();

  fetch(`${this.baseUrl}/search-postcode?postcode=${this.$input.value}`)
    .then((r) => r.json())
    .then((results) => {
      if (!results || !Array.isArray(results)) {
        throw new Error("No results found");
      }

      this.$button.disabled = false;
      this.results = results;

      this.$dropdown.innerHTML = "";

      const $initialOption = document.createElement("option");
      $initialOption.value = "";
      $initialOption.innerHTML = "Select address";
      $initialOption.selected = true;
      $initialOption.disabled = true;
      $initialOption.hidden = true;
      this.$dropdown.appendChild($initialOption);

      this.results.forEach((result, i) => {
        const $option = document.createElement("option");
        $option.value = i.toString();
        $option.innerHTML = result.description;

        this.$dropdown.appendChild($option);
      });

      this.$dropdownContainer.classList.remove("govuk-!-display-none");
      this.$dropdown.focus();
    })
    .catch((err) => {
      this.$button.disabled = false;

      this.showError();
    });
};

AddressFinder.prototype.underwriteValue = function (field, value) {
  let $input = this.$module.querySelector(`[name="${field}"]`);

  if (!$input) {
    $input = this.$module.querySelector(
      `[data-app-address-finder-map="${field}"]`,
    );
  }

  if (
    $input instanceof HTMLInputElement ||
    $input instanceof HTMLSelectElement
  ) {
    $input.value = value;
  }
};

AddressFinder.prototype.handleSelect = function () {
  const i = parseInt(this.$dropdown.value, 10);
  const result = this.results[i];

  Object.entries(result).forEach(([field, value]) =>
    this.underwriteValue(field, value),
  );

  if (this.fillCountry !== "false") {
    this.underwriteValue("country", "GB");
  }
};

export default function init(prefix, $scope) {
  const $addressFinders = ($scope || document).querySelectorAll(
    '[data-module="app-address-finder"]',
  );

  nodeListForEach($addressFinders, ($addressFinder) => {
    new AddressFinder($addressFinder, { prefix });
  });
}
