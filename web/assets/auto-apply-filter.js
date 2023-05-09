export default function autoApplyFilter() {
    const filters = document.querySelectorAll("input.govuk-checkboxes__input");

    filters.forEach((filter) => {
        filter.addEventListener('click', () => {
            setTimeout(() => {
                document.forms["search-filters"].submit();
            }, 500)
        })
    })
}
