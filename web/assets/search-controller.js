const searchController = () => {
  const searchInput = document.querySelector('[data-module="search"] input');

  window.addEventListener("load", (event) => {
    if (window.location.pathname.includes("/search")) {
      let params = new URLSearchParams(window.location.search);
      searchInput.value = params.get("term");
    }
  });
};
export default searchController;
