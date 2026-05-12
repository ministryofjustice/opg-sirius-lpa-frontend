export default async function copyToClipboard() {
  document.addEventListener("click", (e) => {
    if (e.target.matches("button[data-copy-to-clipboard]")) {
      e.preventDefault();

      const copyButton = e.target;

      navigator.clipboard.writeText(copyButton.dataset.copyToClipboard);

      const originalButtonText = copyButton.innerText;
      copyButton.classList.add("disable-click");
      copyButton.inert = true;
      copyButton.textContent = "Copied";

      const screenReaderAlert = document.createElement("span");
      screenReaderAlert.classList.add("govuk-visually-hidden");
      screenReaderAlert.ariaLive = "polite";
      screenReaderAlert.textContent = "Copied to clipboard";
      copyButton.parentElement.appendChild(screenReaderAlert);

      setTimeout(() => {
        copyButton.classList.remove("disable-click");
        copyButton.inert = false;
        copyButton.textContent = originalButtonText;
        copyButton.parentElement.removeChild(screenReaderAlert);
      }, 4000);

      copyButton.blur();
    }
  })
}
