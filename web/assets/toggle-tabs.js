const toggleTabs = () => {
    /** @type NodeList|null tabs */
    const tabs = document.querySelectorAll(
        '[data-module="tab-link"]'
    );

    /** @type NodeList|null tabSections */
    const tabSections = document.querySelectorAll(
        '[data-module="tab-content"]'
    );

    if (tabs && tabs.length > 0) {
        const defaultSelectedTab = tabs[0]
        defaultSelectedTab.classList.add("govuk-tabs__list-item--selected")

        tabSections.forEach((section) => {
            if(section.id !== defaultSelectedTab.id){
                section.classList.add("govuk-tabs__panel--hidden")
            }
        });

        tabs.forEach((el, i) => {
            el.addEventListener("click", (event) => {
                tabs.forEach((otherEl) => {
                    el.id === otherEl.id ? el.classList.add("govuk-tabs__list-item--selected") : otherEl.classList.remove("govuk-tabs__list-item--selected");
                })

                tabSections.forEach((section) => {
                    section.id !== el.id ? section.classList.add("govuk-tabs__panel--hidden") : section.classList.remove("govuk-tabs__panel--hidden");
                });
            });
        })
    }
};

export default toggleTabs;

