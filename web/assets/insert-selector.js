import { nodeListForEach } from "@ministryofjustice/frontend";
import { createElement as el } from "./lib/createElement";
import { initAll as initGOVUKFrontend } from "govuk-frontend";
import handleInsertCheckboxes from "./handle-insert-checkboxes";

function InsertSelector($module) {
  this.$module = $module;

  const $dataSource = $module.querySelector('[data-id="insert-selector-data"]');
  this.data = JSON.parse($dataSource.innerHTML);

  const selectorAttribute = $module.getAttribute("data-initiator-selector");
  this.$initiator = document.querySelector(`${selectorAttribute}-select`);

  this.$containerTemplate = $module.querySelector(
    '[data-id="insert-selector-template-container"]'
  );
  this.$panelTemplate = $module.querySelector(
    '[data-id="insert-selector-template-panel"]'
  );

  this.templateId = "";
}

InsertSelector.prototype.init = function () {
  this.$initiator.addEventListener("confirm", this.onSelectTemplate.bind(this));
};

InsertSelector.prototype.onSelectTemplate = function (e) {
  if (this.$initiator.value === this.templateId) {
    return;
  }

  this.templateId = this.$initiator.value;
  const template = this.data.templates.find((x) => x.id === this.templateId);

  if (template) {
    this.buildSelector(template.inserts);
  }
};

InsertSelector.prototype.buildSelector = function (inserts) {
  this.resetSelector();

  if (Object.keys(inserts).length) {
    this.populateSelector(inserts);
  }
};

InsertSelector.prototype.resetSelector = function () {
  if (this.$container) {
    this.$module.removeChild(this.$container);
    this.$container = null;
  }
};

InsertSelector.prototype.populateSelector = function (insertLists) {
  this.$container = this.$containerTemplate.content.children[0].cloneNode(true);

  this.$tabContainer = this.$container.querySelector(".govuk-tabs__list");
  this.$panelContainer = this.$container.querySelector(
    '[data-module="govuk-tabs"]'
  );

  if (!insertLists.all) {
    const all = [];

    Object.values(insertLists).forEach((inserts) => {
      inserts.forEach((insert) => {
        if (!all.includes(insert)) all.push(insert);
      });
    });

    insertLists = { all, ...insertLists };
  }

  Object.entries(insertLists).forEach(([key, inserts]) => {
    this.$tabContainer.appendChild(
      el(
        "li",
        {
          class: "govuk-tabs__list-item",
        },
        [
          el("a", { class: "govuk-tabs__tab", href: `#panel-${key}` }, [
            key.charAt(0).toUpperCase() + key.slice(1),
          ]),
        ]
      )
    );

    const $rows = inserts.map((insert) => {
      return el(
        "tr",
        {
          class: "govuk-table__row app-!-table-row__no-border",
        },
        [
          el("td", { class: "govuk-table__cell" }, [
            el("div", { class: "govuk-checkboxes__item" }, [
              el("input", {
                class: "govuk-checkboxes__input",
                id: `f-${insert.id}-${key}`,
                name: "insert",
                type: "checkbox",
                value: insert.id,
                "data-module": "insert-checkbox",
              }),
              el(
                "label",
                {
                  class: "govuk-label govuk-checkboxes__label",
                  for: `f-${insert.id}-${key}`,
                },
                [`${insert.id}: ${insert.label}`]
              ),
            ]),
          ]),
        ]
      );
    });

    const $panel = this.$panelTemplate.content.children[0].cloneNode(true);

    $panel.setAttribute("id", `panel-${key}`);
    $rows.forEach(($child) =>
      $panel.querySelector("tbody").appendChild($child)
    );

    this.$panelContainer.appendChild($panel);
  });

  this.$module.appendChild(this.$container);

  initGOVUKFrontend({ scope: this.$module });
  handleInsertCheckboxes({ scope: this.$module });
};

export default function init($scope) {
  const $insertSelectors = ($scope || document).querySelectorAll(
    '[data-module="app-insert-selector"]'
  );

  nodeListForEach($insertSelectors, ($insertSelector) => {
    new InsertSelector($insertSelector).init();
  });
}
