document.body.className +=
  " js-enabled" +
  ("noModule" in HTMLScriptElement.prototype
    ? " govuk-frontend-supported"
    : "");

if (window.self !== window.parent) {
  document.documentElement.className += " app-!-html-class--embedded";
  document.body.className += " app-!-embedded";
}

if (
  document.cookie.indexOf("siriusTheme=dark") > -1 ||
  document.cookie.indexOf("siriusTheme=accessible-dark") > -1
) {
  document.documentElement.className += " app-!-html-class--dark";
}
