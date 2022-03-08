import MOJFrontend from "@ministryofjustice/frontend/moj/all.js";
import GOVUKFrontend from "govuk-frontend/govuk/all.js";
import $ from "jquery";
import "./main.scss";

document.body.className = document.body.className
  ? document.body.className + " js-enabled"
  : "js-enabled";

// Expose jQuery on window so MOJFrontend can use it
window.$ = $;

// we aren't using the JS tabs, but they try to initialise this will stop them breaking
GOVUKFrontend.Tabs.prototype.setup = () => {};

GOVUKFrontend.initAll();
MOJFrontend.initAll();
