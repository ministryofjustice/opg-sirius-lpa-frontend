export default function closeDatePickerOnZoom(scope) {
  scope = scope || document;

  let lastDPR = window.devicePixelRatio;
  const dateInputs = scope.querySelectorAll('input[type="date"]');

  function closeNativePicker(input) {
    if (scope.activeElement !== input) return;

    input.blur();

    //Fallback as the blur event can be inconsistent on some browsers
    const originalType = input.type;
    input.type = "text";
    input.type = originalType;
  }

  window.addEventListener("resize", () => {
    if (window.devicePixelRatio !== lastDPR) {
      lastDPR = window.devicePixelRatio;
      dateInputs.forEach((input) => {
        closeNativePicker(input);
      });
    }
  });
}
