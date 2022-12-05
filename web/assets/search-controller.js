const searchController = () => {
    const searchForm = document.querySelector('[data-module="search"]');
    const searchInput = document.querySelector('[data-module="search"] input');

    if (searchForm && searchInput) {
        searchForm.addEventListener('submit', function(e) {
            if (searchInput.value.length !== 0) {
                window.location.href = window.location.origin + "/search?term=" + searchInput.value;
            }
        });

        window.addEventListener("load", (event) => {
            if (window.location.pathname.includes("/search")) {
                let params = new URLSearchParams(window.location.search);
                searchInput.value = params.get('term');
            }
        });
    }
};
export default searchController;
