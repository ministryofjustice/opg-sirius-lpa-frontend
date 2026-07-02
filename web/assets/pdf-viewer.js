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
    this.rotation = this.loadRotation();
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
    globalThis.addEventListener("beforeunload", () => this.saveState());
  }

  getRotationStorageKey() {
    return `pdf-viewer-rotation-${this.url}`;
  }

  loadRotation() {
    try {
      const stored = sessionStorage.getItem(this.getRotationStorageKey());
      const rotation = stored ? Number.parseInt(stored, 10) : 0;
      return Number.isNaN(rotation) ? 0 : rotation;
    } catch {
      return 0;
    }
  }

  saveRotation() {
    try {
      sessionStorage.setItem(
        this.getRotationStorageKey(),
        this.rotation.toString(),
      );
    } catch {
      // Ignore storage errors
    }
  }

  async init() {
    try {
      this.createControls();
      this.createMainArea();
      this.setupStatePreservation();

      this.container.setAttribute(
        "aria-label",
        `Select Document ${this.paneId}`,
      );
      this.container.setAttribute("role", "region");

      if (this.paneId === "1") {
        this.container.classList.add("pdf-viewer-pane--active");
        this.container.setAttribute("aria-current", "true");
        this.pagesWrapper.focus();
      } else {
        this.container.setAttribute("aria-current", "false");
      }

      // Load saved state before loading PDF
      const savedState = this.loadState();

      // Restore state if available
      // Each pane maintains its own state per document (independent of other panes)
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

      const pageInput = this.controls.querySelector(".pdf-viewer-page-input");
      if (pageInput) pageInput.max = this.totalPages.toString();

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
          Page <input type="number" class="pdf-viewer-page-input" aria-label="Current page number" value="1" min="1"> of <span class="pdf-viewer-total-pages">-</span>
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
      <div class="pdf-viewer-controls-group">
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="rotate-cw">
          Rotate Clockwise
        </button>
        <button type="button" class="govuk-button govuk-button--secondary pdf-viewer-btn" data-action="rotate-ccw">
          Rotate Counterclockwise
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

    globalThis.addEventListener("keydown", async (e) => {
      // Handle pane selection keys (1 and 2) - these work regardless of focus
      if (e.key === "1" && this.paneId === "1") {
        this.highlightPane();
        return;
      }

      if (e.key === "2" && this.paneId === "2") {
        this.highlightPane();
        return;
      }

      // For arrow keys, only respond if this viewer has focus
      if (
        this.pagesWrapper !== e.target &&
        !this.pagesWrapper.contains(e.target)
      ) {
        return; // Exit if focus is not in this viewer
      }

      switch (e.key) {
        case "ArrowRight":
          e.preventDefault();
          await this.nextPage();
          break;
        case "ArrowLeft":
          e.preventDefault();
          await this.prevPage();
          break;
      }
    });

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
    pagesWrapper.setAttribute("tabindex", "0");
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

    // Highlight this pane on click
    canvasContainer.addEventListener("click", () => {
      this.highlightPane();
    });
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
      case "rotate-cw":
        this.rotateCW();
        break;
      case "rotate-ccw":
        this.rotateCCW();
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
    this.pageContainers = [];

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

        const outputScale = globalThis.devicePixelRatio || 1;

        canvas.width = Math.floor(viewport.width * outputScale);
        canvas.height = Math.floor(viewport.height * outputScale);
        canvas.style.width = Math.floor(viewport.width) + "px";
        canvas.style.height = Math.floor(viewport.height) + "px";
        pageContainer.style.width = canvas.style.width;
        pageContainer.style.height = canvas.style.height;

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
        this.pageContainers.push(pageContainer);
        this.pageCanvases.push(canvas);
      }

      // Re-apply rotation to newly rendered canvases
      this.applyRotation();

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

  highlightPane() {
    this.container.classList.add("pdf-viewer-pane--active");
    this.container.setAttribute("aria-current", "true");
    this.pagesWrapper.focus();

    document.querySelectorAll("[data-pdf-viewer]").forEach((viewer) => {
      if (viewer !== this.container) {
        viewer.classList.remove("pdf-viewer-pane--active");
        viewer.setAttribute("aria-current", "false");
      }
    });
  }

  async preserveScrollPosition(callback) {
    // Calculate scroll ratios before action
    const scrollLeft = this.canvasContainer.scrollLeft;
    const scrollTop = this.canvasContainer.scrollTop;
    const scrollWidth = this.canvasContainer.scrollWidth;
    const scrollHeight = this.canvasContainer.scrollHeight;

    const scrollLeftRatio = scrollWidth > 0 ? scrollLeft / scrollWidth : 0;
    const scrollTopRatio = scrollHeight > 0 ? scrollTop / scrollHeight : 0;

    // Execute the action (handles both sync and async callbacks)
    await callback();

    // Restore scroll position using ratios
    this.canvasContainer.scrollLeft =
      scrollLeftRatio * this.canvasContainer.scrollWidth;
    this.canvasContainer.scrollTop =
      scrollTopRatio * this.canvasContainer.scrollHeight;
  }

  async goToPage(pageNum) {
    if (pageNum < 1 || pageNum > this.totalPages) return;
    if (pageNum === this.currentPage) return;

    // Save current scroll offset within the page
    const currentPageContainer = this.pagesWrapper.querySelector(
      `[data-page="${this.currentPage}"]`,
    );
    let scrollOffset = 0;
    if (currentPageContainer) {
      scrollOffset =
        this.canvasContainer.scrollTop - currentPageContainer.offsetTop;
    }

    this.currentPage = pageNum;
    this.updatePageInfo();
    this.updateThumbnailSelection();

    // Scroll to the new page with same offset
    const pageContainer = this.pagesWrapper.querySelector(
      `[data-page="${pageNum}"]`,
    );
    if (pageContainer) {
      this.canvasContainer.scrollTop = pageContainer.offsetTop + scrollOffset;
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
    console.log("this.currentPage, this.totalPages");
    console.log(this.currentPage, this.totalPages);
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
    await this.preserveScrollPosition(async () => {
      this.updatePageInfo();
      await this.renderAllPages();
    });
  }

  async fitToWidth() {
    if (!this.pdfDoc) return;

    const page = await this.pdfDoc.getPage(this.currentPage);
    const viewport = page.getViewport({ scale: 1 });
    const containerWidth = this.canvasContainer.clientWidth - 40; // Account for padding
    if ([90, 270].includes(this.rotation)) {
      this.scale = containerWidth / viewport.height;
    } else {
      this.scale = containerWidth / viewport.width;
    }
    await this.applyZoom();
    // applyZoom already saves compare state
  }

  rotateCW() {
    this.rotation = (this.rotation + 90) % 360;
    this.saveRotation();
    this.applyRotation();
  }

  rotateCCW() {
    this.rotation = (this.rotation - 90 + 360) % 360;
    this.saveRotation();
    this.applyRotation();
  }

  applyRotation() {
    // Apply rotation transform to all page canvases.
    // Apply translation to the canvases and flip the dimensions of the page containers
    // for canvases that have been rotated 90 degrees clockwise or counterclockwise
    // to remove the extra space caused by the rotation.
    this.pageCanvases.forEach((canvas, index) => {
      const pageContainer = this.pageContainers[index];
      if ([90, 270].includes(this.rotation)) {
        const width = Number.parseInt(canvas.style.width, 10);
        const height = Number.parseInt(canvas.style.height, 10);
        canvas.style.transform = `rotate(${this.rotation}deg) translate(${(this.rotation === 90 ? -1 : 1) * ((height - width) / 2)}px, 0)`;
        pageContainer.style.width = canvas.style.height;
        pageContainer.style.height = canvas.style.width;
      } else {
        canvas.style.transform = `rotate(${this.rotation}deg)`;
        pageContainer.style.width = canvas.style.width;
        pageContainer.style.height = canvas.style.height;
      }
    });
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
    } else if (["-", "+", "e", "."].includes(e.key)) {
      e.preventDefault();
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
