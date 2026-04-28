export default function showHideActions() {
  const panel = document.querySelector(".actions-panel");
  const button = document.querySelector("#actions-toggle");
  const content = document.querySelector("#actions-content");
  const inner = document.querySelector(".actions-panel__inner");

  if (!panel || !button || !content || !inner) {
    return;
  }

  const measureHeight = () => `${Math.max(inner.scrollHeight, 150)}px`;

  const syncHeightIfOpen = () => {
    if (isOpen()) {
      content.style.height = measureHeight();
    }
  };

  const setOpen = (open) => {
    panel.classList.toggle("actions-panel--open", open);
    button.setAttribute("aria-expanded", String(open));

    if (open) {
      content.hidden = false;
      content.style.height = "0px";
      // eslint-disable-next-line no-unused-expressions
      content.offsetHeight;
      content.style.height = measureHeight();
    } else {
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

  if (typeof ResizeObserver !== "undefined") {
    const ro = new ResizeObserver(() => syncHeightIfOpen());
    ro.observe(inner);
  }

  document.body.addEventListener("htmx:afterSwap", () => syncHeightIfOpen());
}
