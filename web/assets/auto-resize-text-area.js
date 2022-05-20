const autoResizeTextArea = () => {
  const txtareas = document.querySelectorAll('[data-module="app-auto-resize"]');
  for (let i = 0; i < txtareas.length; i++) {
    txtareas[i].style.height = "auto";
    txtareas[i].style.height = txtareas[i].scrollHeight + 4 + "px";

    txtareas[i].addEventListener(
      "input",
      function () {
        txtareas[i].style.height = "auto";
        txtareas[i].style.height = txtareas[i].scrollHeight + 4 + "px";
      },
      false
    );
  }
};
export default autoResizeTextArea;
