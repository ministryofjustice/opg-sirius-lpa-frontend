// Converts the date input to a text input on double click to allow for copying
// and pasting the date values. Converts the text input back to a date input on focus out.
export default function dateInputToText(scope) {
  scope = scope || document;

  const dateInputs = scope.querySelectorAll(
    'input[type="date"][data-module="app-date-input-to-text"]',
  );

  if (!dateInputs.length) return;

  dateInputs.forEach((dateInput) => {
    dateInput.addEventListener("dblclick", () => handleDoubleClick(dateInput));
    dateInput.addEventListener("focusout", () => handleFocusOut(dateInput));
    dateInput.addEventListener("keydown", (e) =>
      handlePreventSubmit(e, dateInput),
    );
  });

  function handleDoubleClick(dateInput) {
    dateInput.type = "text";
    dateInput.value = dateInput.value.split("-").reverse().join("/");

    setTimeout(() => {
      dateInput.select();
    });
  }

  function handleFocusOut(dateInput) {
    dateInput.value = dateInput.value.split("/").reverse().join("-");
    dateInput.type = "date";
  }

  function handlePreventSubmit(e, dateInput) {
    if (e.key === "Enter" && dateInput.type === "text") {
      e.preventDefault();
      dateInput.dispatchEvent(new Event("focusout"));
    }
  }
}
