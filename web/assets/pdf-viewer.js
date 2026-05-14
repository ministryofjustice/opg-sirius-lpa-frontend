import * as pdfjsLib from "pdfjs-dist";

// Set the worker source path - worker file is copied to static/javascript during build
const prefix = document.body.getAttribute("data-prefix") || "";
pdfjsLib.GlobalWorkerOptions.workerSrc = `${prefix}/javascript/pdf.worker.min.mjs`;

class PDFViewer {
  constructor(container, url) {
    this.container = container;
    this.url = url;
    this.pdfDoc = null;
    this.currentPage = 1;
    this.totalPages = 0;
    this.scale = 1.0;
    this.rendering = false;
    this.thumbnailsVisible = false;
    this.thumbnailsRendered = false;
    this.pageCanvases = [];
    this.isScrolling = false;

    this.init();
  }

  async init() {
    try {
      this.createControls();
      this.createMainArea();

      const loadingTask = pdfjsLib.getDocument(this.url);
      this.pdfDoc = await loadingTask.promise;
      this.totalPages = this.pdfDoc.numPages;

      this.updatePageInfo();
      await this.renderAllPages();

      // Add keyboard navigation
      document.addEventListener("keydown", (e) =>
        this.handleKeyboardNavigation(e),
      );
    } catch (error) {
      console.error("Error loading PDF:", error);
      this.showError("Unable to load PDF document");
    }
  }

  createControls() {
    const controls = document.createElement("div");
    controls.className = "pdf-viewer-controls";
    controls.innerHTML = `
      <div class="pdf-viewer-controls-group">
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="toggle-thumbnails" aria-label="Toggle thumbnails" aria-expanded="false">
          Thumbnails
        </button>
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="prev" aria-label="Previous page">
          <span aria-hidden="true">←</span> Previous
        </button>
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="next" aria-label="Next page">
          Next <span aria-hidden="true">→</span>
        </button>
        <span class="pdf-viewer-page-info">
          Page <input type="number" class="pdf-viewer-page-input" aria-label="Current page number" value="1" min="1" max="1"> of <span class="pdf-viewer-total-pages">-</span>
        </span>
      </div>
      <div class="pdf-viewer-controls-group">
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="zoom-out" aria-label="Zoom out">
          <span aria-hidden="true">−</span>
        </button>
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="zoom-in" aria-label="Zoom in">
          <span aria-hidden="true">+</span>
        </button>
        <span class="pdf-viewer-scale-select-container">
          <select class="pdf-viewer-scale-select" aria-label="Zoom level">
            <option value="auto" data-scale-type="auto">Automatic Zoom</option>
            <option value="page-actual" data-scale-type="page-actual">Actual Size</option>
            <option value="page-fit" data-scale-type="page-fit">Page Fit</option>
            <option value="page-width" data-scale-type="page-width" selected>Page Width</option>
            <option value="0.5">50%</option>
            <option value="0.75">75%</option>
            <option value="1">100%</option>
            <option value="1.25">125%</option>
            <option value="1.5">150%</option>
            <option value="2">200%</option>
            <option value="3">300%</option>
            <option value="4">400%</option>
          </select>
        </span>
      </div>
    `;

    controls.addEventListener("click", (e) => this.handleControlClick(e));
    this.container.appendChild(controls);
    this.controls = controls;

    // Add event listener for page input
    const pageInput = controls.querySelector(".pdf-viewer-page-input");
    if (pageInput) {
      pageInput.addEventListener("keydown", (e) =>
        this.handlePageInputKeydown(e),
      );
      pageInput.addEventListener("blur", (e) => this.handlePageInputBlur(e));
    }

    // Add event listener for scale select
    const scaleSelect = controls.querySelector(".pdf-viewer-scale-select");
    if (scaleSelect) {
      scaleSelect.addEventListener("change", (e) =>
        this.handleScaleSelectChange(e),
      );
    }
  }

  createMainArea() {
    const mainArea = document.createElement("div");
    mainArea.className = "pdf-viewer-main-area";

    // Create thumbnail panel
    const thumbnailPanel = document.createElement("div");
    thumbnailPanel.className = "pdf-viewer-thumbnail-panel";
    thumbnailPanel.setAttribute("aria-label", "Page thumbnails");

    const thumbnailList = document.createElement("div");
    thumbnailList.className = "pdf-viewer-thumbnail-list";
    thumbnailPanel.appendChild(thumbnailList);

    // Create canvas container
    const canvasContainer = document.createElement("div");
    canvasContainer.className = "pdf-viewer-canvas-container";

    // Create pages wrapper for all pages
    const pagesWrapper = document.createElement("div");
    pagesWrapper.className = "pdf-viewer-pages-wrapper";
    canvasContainer.appendChild(pagesWrapper);

    mainArea.appendChild(thumbnailPanel);
    mainArea.appendChild(canvasContainer);
    this.container.appendChild(mainArea);

    this.mainArea = mainArea;
    this.thumbnailPanel = thumbnailPanel;
    this.thumbnailList = thumbnailList;
    this.canvasContainer = canvasContainer;
    this.pagesWrapper = pagesWrapper;

    // Add scroll listener to track current page
    canvasContainer.addEventListener("scroll", () => this.handleScroll());
  }

  handleControlClick(e) {
    const btn = e.target.closest("[data-action]");
    if (!btn) return;

    const action = btn.dataset.action;
    switch (action) {
      case "prev":
        this.prevPage();
        break;
      case "next":
        this.nextPage();
        break;
      case "zoom-in":
        this.zoomIn();
        break;
      case "zoom-out":
        this.zoomOut();
        break;
      case "toggle-thumbnails":
        this.toggleThumbnails();
        break;
    }
  }

  handleScroll() {
    if (this.isScrolling) return;

    const containerRect = this.canvasContainer.getBoundingClientRect();
    const containerMiddle = containerRect.top + containerRect.height / 2;

    let closestPage = 1;
    let closestDistance = Infinity;

    this.pageCanvases.forEach((canvas, index) => {
      const rect = canvas.getBoundingClientRect();
      const pageMiddle = rect.top + rect.height / 2;
      const distance = Math.abs(pageMiddle - containerMiddle);

      if (distance < closestDistance) {
        closestDistance = distance;
        closestPage = index + 1;
      }
    });

    if (closestPage !== this.currentPage) {
      this.currentPage = closestPage;
      this.updatePageInfo();
      this.updateThumbnailSelection();
    }
  }

  async renderAllPages() {
    if (this.rendering) return;
    this.rendering = true;

    // Clear existing pages
    this.pagesWrapper.innerHTML = "";
    this.pageCanvases = [];

    try {
      for (let pageNum = 1; pageNum <= this.totalPages; pageNum++) {
        const page = await this.pdfDoc.getPage(pageNum);
        const viewport = page.getViewport({ scale: this.scale });

        // Create page container
        const pageContainer = document.createElement("div");
        pageContainer.className = "pdf-viewer-page";
        pageContainer.dataset.page = pageNum;

        const canvas = document.createElement("canvas");
        canvas.className = "pdf-viewer-canvas";

        // Support high-DPI displays
        const outputScale = window.devicePixelRatio || 1;

        canvas.width = Math.floor(viewport.width * outputScale);
        canvas.height = Math.floor(viewport.height * outputScale);
        canvas.style.width = Math.floor(viewport.width) + "px";
        canvas.style.height = Math.floor(viewport.height) + "px";

        const ctx = canvas.getContext("2d");
        const transform =
          outputScale !== 1 ? [outputScale, 0, 0, outputScale, 0, 0] : null;

        const renderContext = {
          canvasContext: ctx,
          transform: transform,
          viewport: viewport,
        };

        await page.render(renderContext).promise;

        pageContainer.appendChild(canvas);
        this.pagesWrapper.appendChild(pageContainer);
        this.pageCanvases.push(canvas);
      }

      this.rendering = false;
    } catch (error) {
      console.error("Error rendering pages:", error);
      this.rendering = false;
    }
  }

  async toggleThumbnails() {
    this.thumbnailsVisible = !this.thumbnailsVisible;
    this.thumbnailPanel.classList.toggle(
      "pdf-viewer-thumbnail-panel--visible",
      this.thumbnailsVisible,
    );

    const toggleBtn = this.controls.querySelector(
      '[data-action="toggle-thumbnails"]',
    );
    if (toggleBtn) {
      toggleBtn.setAttribute(
        "aria-expanded",
        this.thumbnailsVisible.toString(),
      );
      toggleBtn.classList.toggle(
        "pdf-viewer-btn--active",
        this.thumbnailsVisible,
      );
    }

    if (this.thumbnailsVisible && !this.thumbnailsRendered) {
      await this.renderThumbnails();
      this.thumbnailsRendered = true;
    }
  }

  async renderThumbnails() {
    const thumbnailScale = 0.2;

    for (let pageNum = 1; pageNum <= this.totalPages; pageNum++) {
      const page = await this.pdfDoc.getPage(pageNum);
      const viewport = page.getViewport({ scale: thumbnailScale });

      const thumbnailItem = document.createElement("button");
      thumbnailItem.type = "button";
      thumbnailItem.className = "pdf-viewer-thumbnail-item";
      if (pageNum === this.currentPage) {
        thumbnailItem.classList.add("pdf-viewer-thumbnail-item--active");
      }
      thumbnailItem.setAttribute("aria-label", `Go to page ${pageNum}`);
      thumbnailItem.dataset.page = pageNum;

      const canvas = document.createElement("canvas");
      canvas.className = "pdf-viewer-thumbnail-canvas";
      canvas.width = viewport.width;
      canvas.height = viewport.height;

      const ctx = canvas.getContext("2d");
      await page.render({
        canvasContext: ctx,
        viewport: viewport,
      }).promise;

      const pageLabel = document.createElement("span");
      pageLabel.className = "pdf-viewer-thumbnail-label";
      pageLabel.textContent = pageNum;

      thumbnailItem.appendChild(canvas);
      thumbnailItem.appendChild(pageLabel);
      thumbnailItem.addEventListener("click", () => this.goToPage(pageNum));

      this.thumbnailList.appendChild(thumbnailItem);
    }
  }

  updateThumbnailSelection() {
    const thumbnails = this.thumbnailList.querySelectorAll(
      ".pdf-viewer-thumbnail-item",
    );
    thumbnails.forEach((thumb) => {
      const pageNum = parseInt(thumb.dataset.page, 10);
      thumb.classList.toggle(
        "pdf-viewer-thumbnail-item--active",
        pageNum === this.currentPage,
      );
    });
  }

  async goToPage(pageNum) {
    if (pageNum < 1 || pageNum > this.totalPages) return;
    if (pageNum === this.currentPage) return;

    this.currentPage = pageNum;
    this.updatePageInfo();
    this.updateThumbnailSelection();

    // Scroll to the page
    const pageContainer = this.pagesWrapper.querySelector(
      `[data-page="${pageNum}"]`,
    );
    if (pageContainer) {
      this.isScrolling = true;
      pageContainer.scrollIntoView({ behavior: "auto", block: "start" });
      // Reset scrolling flag after animation
      setTimeout(() => {
        this.isScrolling = false;
      }, 500);
    }
  }

  updatePageInfo() {
    const pageInput = this.controls.querySelector(".pdf-viewer-page-input");
    const totalPagesEl = this.controls.querySelector(".pdf-viewer-total-pages");

    if (pageInput) {
      pageInput.value = this.currentPage;
      pageInput.max = this.totalPages;
    }
    if (totalPagesEl) totalPagesEl.textContent = this.totalPages;

    // Update button states
    const prevBtn = this.controls.querySelector('[data-action="prev"]');
    const nextBtn = this.controls.querySelector('[data-action="next"]');

    if (prevBtn) prevBtn.disabled = this.currentPage <= 1;
    if (nextBtn) nextBtn.disabled = this.currentPage >= this.totalPages;
  }

  async prevPage() {
    if (this.currentPage <= 1) return;
    await this.goToPage(this.currentPage - 1);
  }

  async nextPage() {
    if (this.currentPage >= this.totalPages) return;
    await this.goToPage(this.currentPage + 1);
  }

  async zoomIn() {
    this.scale = Math.min(this.scale * 1.25, 5.0);
    this.updatePageInfo();
    await this.renderAllPages();
    // Scroll back to current page after re-render
    const pageContainer = this.pagesWrapper.querySelector(
      `[data-page="${this.currentPage}"]`,
    );
    if (pageContainer) {
      pageContainer.scrollIntoView({ block: "start" });
    }
  }

  async zoomOut() {
    this.scale = Math.max(this.scale / 1.25, 0.25);
    this.updatePageInfo();
    await this.renderAllPages();
    // Scroll back to current page after re-render
    const pageContainer = this.pagesWrapper.querySelector(
      `[data-page="${this.currentPage}"]`,
    );
    if (pageContainer) {
      pageContainer.scrollIntoView({ block: "start" });
    }
  }

  showError(message) {
    const errorEl = document.createElement("div");
    errorEl.className = "pdf-viewer-error";
    errorEl.innerHTML = `
      <p class="govuk-error-message">${message}</p>
    `;
    this.container.appendChild(errorEl);
  }

  handlePageInputKeydown(e) {
    if (e.key === "Enter") {
      const pageNum = parseInt(e.target.value, 10);
      if (!isNaN(pageNum)) {
        this.goToPage(pageNum);
      }
    }
  }

  handlePageInputBlur(e) {
    const pageNum = parseInt(e.target.value, 10);
    if (isNaN(pageNum) || pageNum < 1 || pageNum > this.totalPages) {
      e.target.value = this.currentPage;
    } else {
      this.goToPage(pageNum);
    }
  }

  async handleScaleSelectChange(e) {
    const value = e.target.value;
    const scaleType =
      e.target.options[e.target.selectedIndex].dataset.scaleType;

    if (scaleType === "auto") {
      this.scale = 1.0;
    } else if (scaleType === "page-actual") {
      this.scale = 1.0;
    } else if (scaleType === "page-fit") {
      if (!this.pdfDoc) return;
      const page = await this.pdfDoc.getPage(this.currentPage);
      const viewport = page.getViewport({ scale: 1.0 });
      const containerHeight = this.canvasContainer.clientHeight - 40;
      this.scale = containerHeight / viewport.height;
    } else if (scaleType === "page-width") {
      if (!this.pdfDoc) return;
      const page = await this.pdfDoc.getPage(this.currentPage);
      const viewport = page.getViewport({ scale: 1.0 });
      const containerWidth = this.canvasContainer.clientWidth - 40;
      this.scale = containerWidth / viewport.width;
    } else {
      // Numeric scale value
      this.scale = parseFloat(value);
    }

    this.updatePageInfo();
    await this.renderAllPages();
    // Scroll back to current page after re-render
    const pageContainer = this.pagesWrapper.querySelector(
      `[data-page="${this.currentPage}"]`,
    );
    if (pageContainer) {
      pageContainer.scrollIntoView({ block: "start" });
    }
  }

  handleKeyboardNavigation(e) {
    if (e.key === "ArrowLeft") {
      // Left arrow key
      e.preventDefault();
      this.prevPage();
    } else if (e.key === "ArrowRight") {
      // Right arrow key
      e.preventDefault();
      this.nextPage();
    }
  }
}

export default function initPdfViewer() {
  const viewers = document.querySelectorAll("[data-pdf-viewer]");
  viewers.forEach((container) => {
    const url = container.dataset.pdfUrl;
    if (url) {
      new PDFViewer(container, url);
    }
  });
}
