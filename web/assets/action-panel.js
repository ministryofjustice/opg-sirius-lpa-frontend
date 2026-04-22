export default function showHideActions() {
  const panel = document.querySelector(".actions-panel");
  const button = document.querySelector("#actions-toggle");
  const content = document.querySelector("#actions-content");

  if (!panel || !button || !content) {
    return;
  }

  const setOpen = (open) => {
    panel.classList.toggle("actions-panel--open", open);
    button.setAttribute("aria-expanded", String(open));
    content.hidden = !open;
  };

  const isOpen = () => button.getAttribute("aria-expanded") === "true";

  button.addEventListener("click", () => {
    setOpen(!isOpen());
  });
}
