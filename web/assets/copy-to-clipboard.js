export default async function copyToClipboard() {
  const copyButtons = document.querySelectorAll("button[data-copy-to-clipboard]");

  copyButtons.forEach((copyButton) => {
    copyButton.addEventListener("click", (e) => {
      e.preventDefault();

      navigator.clipboard.writeText(copyButton.dataset.copyToClipboard);

      const originalButtonText = copyButton.innerText;
      copyButton.classList.add("disable-click");
      copyButton.textContent = "Copied";

      const screenReaderAlert = document.createElement("span");
      screenReaderAlert.classList.add("govuk-visually-hidden");
      screenReaderAlert.ariaLive = "polite";
      screenReaderAlert.textContent = "Copied to clipboard";
      copyButton.parentElement.appendChild(screenReaderAlert);

      setTimeout(() => {
        copyButton.classList.remove("disable-click");
        copyButton.textContent = originalButtonText;
        copyButton.parentElement.removeChild(screenReaderAlert);
      }, 4000);

      copyButton.blur();
    });
  });
}