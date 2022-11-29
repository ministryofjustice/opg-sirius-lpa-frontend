import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import $ from "jquery";

export default function showHideFilter(prefix) {
    new MOJFrontend.FilterToggleButton({
        bigModeMediaQuery: "(min-width: 48.063em)",
        startHidden: false,
        toggleButton: {
            container: $(".moj-action-bar__filter"),
            showText: "Show filter",
            hideText: "Hide filter",
            classes: "govuk-button--secondary",
        },
        closeButton: {
            container: $(".moj-filter__header-action"),
            text: "Close",
        },
        filter: {
            container: $(".moj-filter"),
        },
    });
}