export default function autoApplyFilter() {
    if (document.body.className.includes('js-enabled')) {
        document.querySelector('.moj-filter form button').classList.add("govuk-!-display-none");
    }

    const filters = document.querySelectorAll('[data-module="app-auto-apply-filter"]');
    let timeout = null;

    filters.forEach((filter) => {
        filter.addEventListener('click', () => {
            if (timeout !== null) clearTimeout(timeout);

            timeout = setTimeout(() => {
                filter.closest("form").submit();
            }, 1000)
        })
    })
}
