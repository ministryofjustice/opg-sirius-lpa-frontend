const autoClick = () => {
  /** @type HTMLAnchorElement|null autoClickLink */
  const autoClickLinks = document.querySelectorAll(
    '[data-module~="app-auto-click"]'
  );

  if (autoClickLinks) {
    autoClickLinks.forEach(($link) => {
      $link.click();
      $link.hidden = true;
    });
  }
};

export default autoClick;
