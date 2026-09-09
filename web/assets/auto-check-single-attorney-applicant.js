export default function autoCheckSingleAttorneyApplicant(scope) {
    scope = scope || document;

    const attorneyRadio = scope.querySelector(
        '[data-module="applicant-attorney-radio"] input',
    );
    const checkboxesContainer = scope.querySelector(
        '[data-module="applicant-attorney-checkboxes"]',
    );

    if (!attorneyRadio || !checkboxesContainer) return;

    attorneyRadio.addEventListener("change", handleAutoCheck);

    function handleAutoCheck() {
        if (!attorneyRadio.checked) return;

        const checkboxes = checkboxesContainer.querySelectorAll(
            'input[type="checkbox"]',
        );

        if (checkboxes.length === 1) {
            checkboxes[0].checked = true;
        }
    }

    handleAutoCheck();
}