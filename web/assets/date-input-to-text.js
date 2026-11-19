export default function dateInputToText(scope) {
  scope = scope || document;

  const dateInputs = scope.querySelectorAll(
    'input[type="date"][data-module="app-date-input-to-text"]',
  );

  dateInputs.forEach((dateInput) => {
    if (dateInput.dataset.dateInputToTextInit) return;
    dateInput.dataset.dateInputToTextInit = "true";

    dateInput.addEventListener("dblclick", () =>
      switchToTextInput(dateInput, { deferSelect: true }),
    );
    dateInput.addEventListener("focusout", () => handleFocusOut(dateInput));
    dateInput.addEventListener("keydown", (e) => handleKeydown(e, dateInput));
  });

  function toDisplayFormat(isoValue) {
    return isoValue.split("-").reverse().join("/");
  }

  function toIsoFormat(displayValue) {
    return displayValue.split("/").reverse().join("-");
  }

  function switchToTextInput(dateInput, { deferSelect = false } = {}) {
    if (dateInput.type === "text") return;

    dateInput.dataset.originalValue = dateInput.value;
    dateInput.type = "text";
    dateInput.value = toDisplayFormat(dateInput.dataset.originalValue);

    if (deferSelect) {
      setTimeout(() => dateInput.select());
    } else {
      dateInput.select();
    }
  }

  function handleFocusOut(dateInput) {
    if (dateInput.type !== "text") return;

    const originalValue = dateInput.dataset.originalValue ?? "";
    const enteredValue = dateInput.value.trim();
    delete dateInput.dataset.originalValue;

    dateInput.value = toIsoFormat(enteredValue);
    dateInput.type = "date";

    if (enteredValue !== "" && dateInput.value === "") {
      dateInput.value = originalValue;
    }
  }

  function copyToClipboard(e, dateInput) {
    if (!dateInput.value || !navigator.clipboard) return;

    e.preventDefault();
    navigator.clipboard
      .writeText(toDisplayFormat(dateInput.value))
      .catch(() => {});
  }

  function handleKeydown(e, dateInput) {
    if (e.key === "Enter" && dateInput.type === "text") {
      e.preventDefault();
      dateInput.dispatchEvent(new Event("focusout"));
      return;
    }

    if (!(e.ctrlKey || e.metaKey) || dateInput.type === "text") return;

    const key = e.key.toLowerCase();

    if (key === "c") {
      copyToClipboard(e, dateInput);
    } else if (key === "v") {
      switchToTextInput(dateInput);
    }
  }
}
