export default function toggleHeaderMenu() {
    const headerContainer = document.querySelector(
        '[data-module="app-header-container"]',
    );

    if (headerContainer) {
        const menuButton = headerContainer.querySelector('[data-module="app-header-menu-toggle"]');
        const menuItemList = headerContainer.querySelector('[data-id="app-header-menu-item-list"]')

        if (menuButton && menuItemList) {
            const tabletMQry = window.matchMedia('(max-width: 768px)');

            if (tabletMQry.matches) {
                menuItemList.hidden = true;
            } else {
                menuButton.hidden = true;
            }

            menuButton.addEventListener("click", () => {
                if (menuItemList.hidden) {
                    menuItemList.hidden = false;
                    menuButton.ariaExpanded = "true";
                } else {
                    menuItemList.hidden = true;
                    menuButton.ariaExpanded = "false";
                }
            });

            tabletMQry.onchange = (e) => {
                if (e.matches) { // screen is 0 - 768px
                    menuItemList.hidden = true;
                    menuButton.ariaExpanded = "false";
                    menuButton.hidden = false;
                } else { // screen is 768 -> bigger
                    menuItemList.hidden = false;
                    menuButton.ariaExpanded = "true";
                    menuButton.hidden = true;
                }
            };
        }
    }
}
