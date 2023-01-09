const previewDraft = () => {
  /** @type HTMLMetaElement|null previewDraft */
  const previewDraft = document.querySelector("[data-app-reload~=\"previewDraft\"]");

  if (previewDraft) {
      window.open(`${previewDraft.id}`,
            'Correspondence preview',
            'height=100,width=100');
  }
};

export default previewDraft;
