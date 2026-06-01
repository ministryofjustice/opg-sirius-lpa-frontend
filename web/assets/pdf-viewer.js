import * as pdfjsLib from "pdfjs-dist";

// Set the worker source path - worker file is copied to static/javascript during build
const prefix = document.body.getAttribute("data-prefix") || "";
pdfjsLib.GlobalWorkerOptions.workerSrc = `${prefix}/javascript/pdf.worker.min.mjs`;

// Storage key for persisting viewer state across page navigations
const STORAGE_KEY = "pdfViewerState";

class PDFViewer {
  constructor(container, url, paneId) {
    this.container = container;
    this.url = url;
    this.paneId = paneId || "1";
    this.pdfDoc = null;
    this.currentPage = 1;
    this.totalPages = 0;
    this.scale = 1;
    this.rendering = false;
    this.thumbnailsVisible = false;
    this.thumbnailsRendered = false;
    this.pageCanvases = [];
    this.isScrolling = false;
    this.storageKey = this.getStorageKey();
  }

  getStorageKey() {
    // Create a unique key based on the PDF URL and pane ID
    // This ensures each pane maintains its own state for the same document
    return `${STORAGE_KEY}_${this.paneId}_${this.url}`;
  }

  saveState() {
    // Save state when navigating away, regardless of mode
    // This allows state to transfer from single-pane view to compare mode
    if (!this.canvasContainer) return;

    const state = {
      scale: this.scale,
      scrollLeftRatio:
        this.canvasContainer.scrollWidth > 0
          ? this.canvasContainer.scrollLeft / this.canvasContainer.scrollWidth
          : 0,
      scrollTopRatio:
        this.canvasContainer.scrollHeight > 0
          ? this.canvasContainer.scrollTop / this.canvasContainer.scrollHeight
          : 0,
      currentPage: this.currentPage,
      timestamp: Date.now(),
    };

    try {
      sessionStorage.setItem(this.storageKey, JSON.stringify(state));
    } catch (e) {
      // sessionStorage might be unavailable or full
      console.warn("Unable to save PDF viewer state:", e);
    }
  }

  loadState() {
    try {
      const stored = sessionStorage.getItem(this.storageKey);
      if (!stored) return null;

      const state = JSON.parse(stored);

      // Only restore state if it's less than 15 minutes old
      const resetTime = 15 * 60 * 1000;
      if (Date.now() - state.timestamp > resetTime) {
        sessionStorage.removeItem(this.storageKey);
        return null;
      }

      return state;
    } catch (e) {
      console.warn("Unable to load PDF viewer state:", e);
      return null;
    }
  }

  setupStatePreservation() {
    // Save state before page unload
    window.addEventListener("beforeunload", () => this.saveState());
  }

  async init() {
    try {
      this.createControls();
      this.createMainArea();
      this.setupStatePreservation();

      // Load saved state before loading PDF
      const savedState = this.loadState();

      // Restore state if available
      // Each pane maintains its own state per document (independent from other panes)
      const shouldRestoreState = savedState;

      if (shouldRestoreState && savedState) {
        this.scale = savedState.scale;
        this.currentPage = savedState.currentPage;
      }

      const loadingTask = pdfjsLib.getDocument(this.url);
      this.pdfDoc = await loadingTask.promise;
      this.totalPages = this.pdfDoc.numPages;

      // Ensure current page is within bounds
      if (this.currentPage > this.totalPages) {
        this.currentPage = this.totalPages;
      }

      this.updatePageInfo();
      await this.renderAllPages();

      // Restore scroll position after rendering
      if (shouldRestoreState && savedState) {
        this.canvasContainer.scrollLeft =
          savedState.scrollLeftRatio * this.canvasContainer.scrollWidth;
        this.canvasContainer.scrollTop =
          savedState.scrollTopRatio * this.canvasContainer.scrollHeight;
      }
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
        <span class="pdf-viewer-page-info">
          Page <input type="text" class="pdf-viewer-page-input" aria-label="Current page number" value="1"> of <span class="pdf-viewer-total-pages">-</span>
        </span>
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="next" aria-label="Next page">
          Next <span aria-hidden="true">→</span>
        </button>
      </div>
      <div class="pdf-viewer-controls-group">
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="zoom-out" aria-label="Zoom out">
          <span aria-hidden="true">−</span>
        </button>
        <input type="text" class="pdf-viewer-zoom-input" aria-label="Zoom level" value="100%">
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="zoom-in" aria-label="Zoom in">
          <span aria-hidden="true">+</span>
        </button>
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="fit-width" aria-label="Fit to width">
          Fit Width
        </button>
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

    // Add event listener for zoom input
    const zoomInput = controls.querySelector(".pdf-viewer-zoom-input");
    if (zoomInput) {
      zoomInput.addEventListener("keydown", (e) =>
        this.handleZoomInputKeydown(e),
      );
      zoomInput.addEventListener("blur", (e) => this.handleZoomInputBlur(e));
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
      case "fit-width":
        this.fitToWidth();
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

        let transform;
        if (outputScale === 1) {
          transform = null;
        } else {
          transform = [outputScale, 0, 0, outputScale, 0, 0];
        }

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
      const pageNum = Number.parseInt(thumb.dataset.page, 10);
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
      pageContainer.scrollIntoView({ behavior: "smooth", block: "start" });
      // Reset scrolling flag after animation
      setTimeout(() => {
        this.isScrolling = false;
      }, 500);
    }
  }

  updatePageInfo() {
    const pageInput = this.controls.querySelector(".pdf-viewer-page-input");
    const totalPagesEl = this.controls.querySelector(".pdf-viewer-total-pages");
    const zoomInput = this.controls.querySelector(".pdf-viewer-zoom-input");

    if (pageInput) pageInput.value = this.currentPage;
    if (totalPagesEl) totalPagesEl.textContent = this.totalPages;
    if (zoomInput) zoomInput.value = Math.round(this.scale * 100) + "%";

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
    this.scale = Math.min(this.scale * 1.25, 5);
    await this.applyZoom();
  }

  async zoomOut() {
    this.scale = Math.max(this.scale / 1.25, 0.25);
    await this.applyZoom();
  }

  async applyZoom() {
    // Calculate scroll ratios before re-rendering
    const scrollLeft = this.canvasContainer.scrollLeft;
    const scrollTop = this.canvasContainer.scrollTop;
    const scrollWidth = this.canvasContainer.scrollWidth;
    const scrollHeight = this.canvasContainer.scrollHeight;

    // Calculate the ratio of scroll position to total scrollable area
    const scrollLeftRatio = scrollWidth > 0 ? scrollLeft / scrollWidth : 0;
    const scrollTopRatio = scrollHeight > 0 ? scrollTop / scrollHeight : 0;

    this.updatePageInfo();
    await this.renderAllPages();

    // Apply the same ratio to the new scrollable area
    this.canvasContainer.scrollLeft =
      scrollLeftRatio * this.canvasContainer.scrollWidth;
    this.canvasContainer.scrollTop =
      scrollTopRatio * this.canvasContainer.scrollHeight;
  }

  async fitToWidth() {
    if (!this.pdfDoc) return;

    const page = await this.pdfDoc.getPage(this.currentPage);
    const viewport = page.getViewport({ scale: 1 });
    const containerWidth = this.canvasContainer.clientWidth - 40; // Account for padding
    this.scale = containerWidth / viewport.width;
    await this.applyZoom();
    // applyZoom already saves compare state
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
      const pageNum = Number.parseInt(e.target.value, 10);
      if (!Number.isNaN(pageNum)) {
        this.goToPage(pageNum);
      }
    }
  }

  handlePageInputBlur(e) {
    const pageNum = Number.parseInt(e.target.value, 10);
    if (Number.isNaN(pageNum) || pageNum < 1 || pageNum > this.totalPages) {
      e.target.value = this.currentPage;
    } else {
      this.goToPage(pageNum);
    }
  }

  handleZoomInputKeydown(e) {
    if (e.key === "Enter") {
      e.target.blur();
    }
  }

  async handleZoomInputBlur(e) {
    const value = e.target.value.replace("%", "").trim();
    const zoomPercent = Number.parseInt(value, 10);

    if (Number.isNaN(zoomPercent) || zoomPercent < 25 || zoomPercent > 500) {
      e.target.value = Math.round(this.scale * 100) + "%";
    } else {
      this.scale = zoomPercent / 100;
      await this.applyZoom();
    }
  }
}

export default function initPdfViewer() {
  const viewers = document.querySelectorAll("[data-pdf-viewer]");
  viewers.forEach((container) => {
    const url = container.dataset.pdfUrl;
    const paneId = container.dataset.pdfPane;
    if (url) {
      const pdfViewer = new PDFViewer(container, url, paneId);
      pdfViewer.init();
    }
  });
}
