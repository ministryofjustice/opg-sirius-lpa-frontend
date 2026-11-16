export default function dateInputClipboard(scope) {
  scope = scope || document;

  const dateInputs = scope.querySelectorAll(
    'input[type="date"][data-module="app-date-input-clipboard"]',
  );

  if (!dateInputs.length) return;

  dateInputs.forEach((dateInput) => {
    dateInput.addEventListener("dblclick", () => handleDoubleClick(dateInput));

    dateInput.addEventListener("focusout", () => handleFocusOut(dateInput));
  });

  function handleDoubleClick(dateInput) {
    console.log("Double click on date input: ", dateInput);
    dateInput.type = "text";
    setTimeout(() => {
      dateInput.select();
    });
  }

  function handleFocusOut(dateInput) {
    dateInput.type = "date";
  }
}
