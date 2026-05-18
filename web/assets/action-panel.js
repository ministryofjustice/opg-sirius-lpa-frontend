export default function showHideActions() {
  const panel = document.querySelector(".actions-panel");
  const button = document.querySelector("#actions-toggle");
  const content = document.querySelector("#actions-content");
  const inner = document.querySelector(".actions-panel__inner");

  if (!panel || !button || !content || !inner) {
    return;
  }

  // Keep the animated height in sync with the content.
  // Default/min height ensures the panel has a usable size even when content
  // is very small (or still loading).
  // The CSS max-height handles overflow via scrolling.
  const measureHeight = () => {
    const desiredHeight = Math.max(inner.scrollHeight, 150);
    return `${desiredHeight}px`;
  };

  const syncHeightIfOpen = () => {
    if (isOpen()) {
      content.style.height = measureHeight();
    }
  };

  const setOpen = (open) => {
    panel.classList.toggle("actions-panel--open", open);
    button.setAttribute("aria-expanded", String(open));

    if (open) {
      // Make it measurable before we set the target height.
      content.hidden = false;
      // Ensure we always transition from 0 to the measured value.
      content.style.height = "0px";
      // Force a reflow so the browser picks up the starting height.
      // eslint-disable-next-line no-unused-expressions
      content.offsetHeight;
      content.style.height = measureHeight();
    } else {
      // Animate to 0, then hide once the animation completes.
      content.style.height = "0px";
      content.addEventListener(
        "transitionend",
        (e) => {
          if (e.propertyName === "height") {
            content.hidden = true;
          }
        },
        { once: true },
      );
    }
  };

  const isOpen = () => button.getAttribute("aria-expanded") === "true";

  button.addEventListener("click", () => {
    setOpen(!isOpen());
  });

  // If content changes (e.g. HTMX swaps), resize while open.
  if (typeof ResizeObserver !== "undefined") {
    const ro = new ResizeObserver(() => syncHeightIfOpen());
    ro.observe(inner);
  }

  // HTMX doesn't necessarily trigger a resize event; hook into swaps.
  document.body.addEventListener("htmx:afterSwap", () => syncHeightIfOpen());
}
