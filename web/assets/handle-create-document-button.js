const handleCreateDocumentButton = () => {
  /*  After creating a non case contact recipient, if a user clicks the create document button without selecting a
        recipient, the recipient they just created will be lost (as non case contacts are not retrieved as part of the
        call to get recipients). So JS is required to disable the button and enable it if at least one recipient has
        been selected */

  /** @type HTMLButtonElement|null createDocumentButton */
  const createDocumentButton = document.querySelector(
    '[data-module="create-document-button"]',
  );

  /** @type NodeList|null checkboxes */
  const checkboxes = document.querySelectorAll(
    '[data-module="recipient-checkbox"]',
  );

  if (createDocumentButton !== null) {
    createDocumentButton.setAttribute("disabled", "");
  }

  if (createDocumentButton && checkboxes && checkboxes.length > 0) {
    checkboxes.forEach((el, i) => {
      el.addEventListener("change", (event) => {
        let isOneRecipientSelected = Array.from(checkboxes).some(
          (x) => x.checked,
        );
        if (isOneRecipientSelected) {
          createDocumentButton.removeAttribute("disabled");
        } else {
          createDocumentButton.setAttribute("disabled", "");
        }
      });
    });
  }
};

export default handleCreateDocumentButton;
