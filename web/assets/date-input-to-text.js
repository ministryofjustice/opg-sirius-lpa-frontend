// Converts the date input to a text input on double click to allow for copying
// and pasting the date values. Converts the text input back to a date input on focus out.
// Restores the original value if the date is invalid.
export default function dateInputToText(scope) {
  scope = scope || document;

  const regExp = /^\d{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$/;

  let currentValue = "";

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
    currentValue = dateInput.value;
    dateInput.value = dateInput.value.split("-").reverse().join("/");

    setTimeout(() => {
      dateInput.select();
    });
  }

  function handleFocusOut(dateInput) {
    dateInput.value = dateInput.value.split("/").reverse().join("-");
    if (!regExp.test(dateInput.value)) {
      dateInput.value = currentValue;
      currentValue = "";
    }

    dateInput.type = "date";
  }

  function handlePreventSubmit(e, dateInput) {
    if (e.key === "Enter" && dateInput.type === "text") {
      e.preventDefault();
      dateInput.dispatchEvent(new Event("focusout"));
    }
  }
}
