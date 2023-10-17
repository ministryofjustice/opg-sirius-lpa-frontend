const init = () => {
  /** @type NodeListOf<HTMLAnchorElement> */
  const modalLinks = document.querySelectorAll(
    '[data-module~="app-modal-link"]',
  );

  if (modalLinks.length === 0) {
    return;
  }

  const $dialog = document.createElement("dialog");
  $dialog.style.padding = "0";
  $dialog.style.width = "50%";
  $dialog.style.height = "50%";
  $dialog.style.overflow = "hidden";

  document.body.appendChild($dialog);

  window.addEventListener("message", (event) => {
    if (
      event.origin !== `${window.location.protocol}//${window.location.host}`
    ) {
      return;
    }

    if (event.data === "form-done") {
      window.location.reload();
    } else if (event.data === "form-cancel") {
      $dialog.close();
    }
  });

  $dialog.addEventListener("click", () => {
    $dialog.close();
  });

  modalLinks.forEach((modalLink) => {
    modalLink.addEventListener("click", (e) => {
      e.preventDefault();

      $dialog.childNodes.forEach(($child) => $dialog.removeChild($child));

      const $iframe = document.createElement("iframe");
      $iframe.style.border = "0";
      $iframe.style.width = "100%";
      $iframe.style.height = "100%";
      $iframe.src = modalLink.href;
      $dialog.appendChild($iframe);
      $dialog.showModal();
    });
  });
};

export default init;
