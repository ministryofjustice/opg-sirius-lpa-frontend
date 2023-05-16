export default function autoApplyFilter() {
    const filters = document.querySelectorAll("input.govuk-checkboxes__input");
    let timeout = null;

    filters.forEach((filter) => {
        filter.addEventListener('click', () => {
            if (timeout !== null) clearTimeout(timeout);

            timeout = setTimeout(() => {
                document.forms["search-filters"].submit();
            }, 500)
        })
    })
}
