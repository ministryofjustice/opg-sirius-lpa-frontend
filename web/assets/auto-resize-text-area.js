const autoResizeTextArea = () => {
  const textareas = document.querySelectorAll(
    '[data-module="app-auto-resize"]',
  );
  for (let textarea of textareas) {
    textarea.style.height = "auto";
    textarea.style.height = textarea.scrollHeight + 4 + "px";

    textarea.addEventListener(
      "input",
      function () {
        textarea.style.height = "auto";
        textarea.style.height = textarea.scrollHeight + 4 + "px";
      },
      false,
    );
  }
};
export default autoResizeTextArea;
