const autoResizeTextArea = () => {
    const txtarea = document.getElementById("auto-resize-textarea");
    for (let i = 0; i < txtarea.length; i++) {
        txtarea[i].setAttribute("style", "height:" + (txtarea[i].scrollHeight) + "px;overflow-y:hidden;");
        txtarea[i].addEventListener("input", OnInput, false);
    }

    function OnInput() {
        this.style.height = "auto";
        this.style.height = (this.scrollHeight) + "px";
    }
};
export default autoResizeTextArea;
