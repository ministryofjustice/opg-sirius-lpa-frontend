const toggleInsertCheckboxes = () => {
    /*Insert checkboxes may appear twice and therefore js is needed to make sure
    that a check in one tab is displayed in the other*/

    /** @type NodeList|null checkboxes */
    const checkboxes = document.querySelectorAll(
        '[data-module="insert-checkbox"]'
    );

    if (checkboxes && checkboxes.length > 0) {
        checkboxes.forEach((el, i) => {
            el.addEventListener("change", (event) => {
                checkboxes.forEach((otherEl) => {
                    if (otherEl.id === el.id) {
                        el.checked ? otherEl.setAttribute("checked", "checked") : otherEl.removeAttribute("checked");
                    }
                })
            });
        })
    }
};

export default toggleInsertCheckboxes;
