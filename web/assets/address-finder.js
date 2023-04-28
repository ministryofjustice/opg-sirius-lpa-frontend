import { nodeListForEach } from "@ministryofjustice/frontend";

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

  this.$inputs = this.$module.querySelectorAll(".govuk-form-group");
  const id = this.$module.id || Math.random().toString(36).substring(2);

  const $container = document.createElement("div");
  $container.innerHTML = AddressFinder.template(id);

  this.$module.appendChild($container);

  /** @type {HTMLInputElement} */
  this.$input = $container.querySelector("input");
  /** @type {HTMLDivElement} */
  this.$dropdownContainer = $container.querySelector(".govuk-details__text");
  /** @type {HTMLSelectElement} */
  this.$dropdown = $container.querySelector("select");
  /** @type {HTMLDivElement} */
  this.$error = $container.querySelector(".govuk-error-message");
  /** @type {HTMLButtonElement} */
  this.$button = $container.querySelector("button");

  this.$input.addEventListener("keydown", this.handleKeydown.bind(this));
  this.$button.addEventListener("click", this.handleSearch.bind(this));
  this.$dropdown.addEventListener("change", this.handleSelect.bind(this));

  const $label = $container.querySelector(`[for="f-${id}-input"]`);
  $label.innerHTML = this.$module.getAttribute("data-app-address-finder-label");

  const $link = $container.querySelector(".govuk-link");
  $link?.addEventListener("click", this.toggleInputs.bind(this));

  this.$inputContainer = document.createElement("div");
  this.$inputContainer.classList.add("govuk-details__text");
  nodeListForEach(this.$inputs, ($input) => {
    this.$inputContainer.appendChild($input);
  });
  this.$module.appendChild(this.$inputContainer);
  this.hideInputs();
}

AddressFinder.template = (id) => `
  <div class="govuk-form-group">
    <label class="govuk-label" for="f-${id}-input"></label>
    <div class="govuk-hint" id="f-${id}-hint">
      Enter a UK postcode, or enter the address manually
    </div>
    <p id="f-${id}-error" class="govuk-error-message govuk-!-display-none">
      <span class="govuk-visually-hidden">Error:</span>
      No matching address found. Please try again using a UK postcode, or enter the address manually
    </p>
    <input
      class="govuk-input govuk-input--width-10"
      id="f-${id}-input"
      aria-describedby="f-${id}-hint"
    />
    <button
      class="govuk-button govuk-button--secondary govuk-!-margin-left-2 govuk-!-margin-bottom-0"
      type="button"
    >
      Find address
    </button>
  </div>
  <div class="govuk-form-group govuk-details__text govuk-!-display-none">
    <label class="govuk-label" for="f-${id}-select">
      Select an address
    </label>
    <select class="govuk-select" id="f-${id}-select"></select>
  </div>
  <div class="govuk-body">
    <a href="#" class="govuk-link govuk-link--no-visited-state">
      Enter address manually
    </a>
  </div>
`;

AddressFinder.prototype.hideInputs = function () {
  this.$inputContainer.classList.add("govuk-!-display-none");
};

AddressFinder.prototype.toggleInputs = function (e) {
  e.preventDefault();

  this.$inputContainer.classList.toggle("govuk-!-display-none");
};

AddressFinder.prototype.showError = function () {
  this.$input.classList.add("govuk-input--error");
  this.$input.setAttribute(
    "aria-describedby",
    this.$input.getAttribute("aria-describedby") + " " + this.$error.id
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
      describedBy.replace(this.$error.id, "").trim()
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
  this.hideInputs();
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
      this.results.forEach((result, i) => {
        const $option = document.createElement("option");
        $option.value = i.toString();
        $option.innerHTML = result.description;

        this.$dropdown.appendChild($option);
      });

      this.$dropdownContainer.classList.remove("govuk-!-display-none");
      this.$dropdown.focus();
      this.$dropdown.dispatchEvent(new Event("change"));
    })
    .catch((err) => {
      this.$button.disabled = false;

      this.showError();
    });
};

AddressFinder.prototype.handleSelect = function () {
  const i = parseInt(this.$dropdown.value, 10);
  const result = this.results[i];

  Object.entries(result).forEach(([field, value]) => {
    const $input = this.$module.querySelector(`[name="${field}"]`);
    if ($input instanceof HTMLInputElement) {
      $input.value = value;
    }
  });
};

export default function init(prefix, $scope) {
  const $addressFinders = ($scope || document).querySelectorAll(
    '[data-module="app-address-finder"]'
  );

  nodeListForEach($addressFinders, ($addressFinder) => {
    new AddressFinder($addressFinder, { prefix });
  });
}
