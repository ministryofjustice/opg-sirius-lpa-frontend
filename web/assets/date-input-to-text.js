// Converts the date input to a text input on double click or on pressing the
// control/meta key to allow for copying and pasting the date values.
// Converts the text input back to a date input on loss of focus.
// Restores the original value if the date is invalid.
export default function dateInputToText(scope) {
  scope = scope || document;

  let currentValue = "";

  const dateInputs = scope.querySelectorAll(
    'input[type="date"][data-module="app-date-input-to-text"]',
  );

  if (!dateInputs.length) return;

  dateInputs.forEach((dateInput) => {
    dateInput.addEventListener("dblclick", () =>
      switchToTextInputAndSelect(dateInput),
    );
    dateInput.addEventListener("focusout", () => handleFocusOut(dateInput));
    dateInput.addEventListener("keydown", (e) => handleKeydown(e, dateInput));
  });

  function switchToTextInputAndSelect(dateInput) {
    if (dateInput.type === "text") return;

    dateInput.type = "text";
    currentValue = dateInput.value;
    dateInput.value = dateInput.value.split("-").reverse().join("/");

    setTimeout(() => {
      dateInput.select();
    });
  }

  function handleFocusOut(dateInput) {
    dateInput.value = dateInput.value.split("/").reverse().join("-");
    dateInput.type = "date";
    if (dateInput.value === "") {
      dateInput.value = currentValue;
      currentValue = "";
    }
  }

  function handleKeydown(e, dateInput) {
    if (e.key === "Enter" && dateInput.type === "text") {
      e.preventDefault();
      dateInput.dispatchEvent(new Event("focusout"));
      return;
    }

    if (e.key === "Control" || e.key === "Meta") {
      switchToTextInputAndSelect(dateInput);
    }
  }
}
