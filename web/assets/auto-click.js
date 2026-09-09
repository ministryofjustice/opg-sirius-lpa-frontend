const autoClick = ($scope) => {
  /** @type HTMLAnchorElement|null autoClickLink */
  const autoClickLinks = ($scope || document).querySelectorAll(
    '[data-module~="app-auto-click"]',
  );

  if (autoClickLinks) {
    autoClickLinks.forEach(($link) => {
      $link.click();
      $link.hidden = true;
    });
  }
};

export default autoClick;
