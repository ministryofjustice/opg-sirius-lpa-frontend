const autoResizeTextArea = () => {
    const txtareas = document.querySelectorAll(".js-auto-resize-textarea");
    for (let i = 0; i < txtareas.length; i++) {
        txtareas[i].style.height = "auto";
        txtareas[i].style.height = (txtareas[i].scrollHeight) + "px";
        txtareas[i].style.overflowY = "hidden";

        txtareas[i].addEventListener("input", function () {
            txtareas[i].style.height = "auto";
            txtareas[i].style.height = (txtareas[i].scrollHeight) + "px";
        }, false);
    }
};
export default autoResizeTextArea;
