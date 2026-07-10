import { PDFViewer } from "./pdf-viewer.js";

// Mock canvas.getContext before importing pdf-viewer
HTMLCanvasElement.prototype.getContext = jest.fn(() => ({
  fillRect: jest.fn(),
  clearRect: jest.fn(),
  getImageData: jest.fn(),
  putImageData: jest.fn(),
  createImageData: jest.fn(),
  setTransform: jest.fn(),
  drawImage: jest.fn(),
  save: jest.fn(),
  fillText: jest.fn(),
  restore: jest.fn(),
  beginPath: jest.fn(),
  moveTo: jest.fn(),
  lineTo: jest.fn(),
  closePath: jest.fn(),
  stroke: jest.fn(),
  translate: jest.fn(),
  scale: jest.fn(),
  rotate: jest.fn(),
  arc: jest.fn(),
  fill: jest.fn(),
  measureText: jest.fn(() => ({ width: 0 })),
  transform: jest.fn(),
  rect: jest.fn(),
  clip: jest.fn(),
}));

Element.prototype.scrollIntoView = jest.fn();

// Mock pdfjs-dist before importing pdf-viewer
jest.mock("pdfjs-dist", () => ({
  GlobalWorkerOptions: {
    workerSrc: "",
  },
  getDocument: jest.fn(() => ({
    promise: Promise.resolve({
      numPages: 3,
      getPage: jest.fn(() =>
        Promise.resolve({
          getViewport: jest.fn(({ scale = 1 }) => ({
            width: 612 * scale,
            height: 792 * scale,
          })),
          render: jest.fn(() => ({
            promise: Promise.resolve(),
          })),
        }),
      ),
    }),
  })),
}));

describe("Pdf Viewer", () => {
  let container;
  let viewer;
  let canvas;

  beforeEach(async () => {
    document.body.innerHTML = `
    <div id="test-container"></div>
  `;

    container = document.getElementById("test-container");
    viewer = new PDFViewer(container, "http://localhost/test.pdf", "1");
    await viewer.init();

    // Wait for rendering to complete
    await new Promise((resolve) => setTimeout(resolve));

    // Get the canvas element
    canvas = document.querySelector("canvas.pdf-viewer-canvas");
    Object.defineProperty(viewer.canvasContainer, "clientWidth", {
      value: 960,
      configurable: true,
    });
  });

  afterEach(() => {
    // Clean up canvas and viewer
    canvas = null;
    viewer = null;
    container = null;
    document.body.innerHTML = "";
    // Clear sessionStorage to reset rotation state
    sessionStorage.clear();
  });

  describe("Rotation button", () => {
    describe("Given the counter clockwise rotation button is pressed", () => {
      test("there will be the counter clockwise style", async () => {
        expect(canvas.style.width).toBe("612px");
        expect(canvas.style.height).toBe("792px");
        expect(canvas.style.transform).toContain("rotate(0deg)");

        viewer.rotateCCW();

        expect(canvas.style.transform).toContain("rotate(270deg)");
        expect(canvas.style.transform).toContain("translate");
      });
    });

    describe("Given the clockwise rotation button is pressed", () => {
      test("there will be the clockwise style", async () => {
        expect(canvas.style.width).toBe("612px");
        expect(canvas.style.height).toBe("792px");
        expect(canvas.style.transform).toContain("rotate(0deg)");

        viewer.rotateCW();

        expect(canvas.style.transform).toContain("rotate(90deg)");
        expect(canvas.style.transform).toContain("translate");
      });
    });
  });

  describe("Given the fit width button is pressed", () => {
    test("the width and height style will change", async () => {
      const initialWidth = canvas.style.width;
      const initialHeight = canvas.style.height;
      expect(initialWidth).toBe("612px");
      expect(initialHeight).toBe("792px");

      await viewer.fitToWidth();

      await new Promise((resolve) => setTimeout(resolve));

      const newCanvas = viewer.pageCanvases[0];
      const newWidth = newCanvas.style.width;
      const newHeight = newCanvas.style.height;

      // The dimensions should change based on the new scale
      expect(newWidth).not.toBe(initialWidth);
      expect(newWidth).toBe("920px");
      expect(newHeight).not.toBe(initialHeight);
      expect(newHeight).toBe("1190px");
    });
  });

  describe("Zoom button", () => {
    describe("Given the zoom in button is pressed", () => {
      test("the canvas will zoom in", async () => {
        const initialWidth = canvas.style.width;
        const initialHeight = canvas.style.height;
        expect(initialWidth).toBe("612px");
        expect(initialHeight).toBe("792px");

        viewer.zoomIn();

        await new Promise((resolve) => setTimeout(resolve));

        const newCanvas = viewer.pageCanvases[0];
        const newWidth = newCanvas.style.width;
        const newHeight = newCanvas.style.height;

        expect(newWidth).not.toBe(initialWidth);
        expect(newWidth).toBe("765px");
        expect(newHeight).not.toBe(initialHeight);
        expect(newHeight).toBe("990px");
      });
    });

    describe("Given the zoom out button is pressed", () => {
      test("the canvas will zoom out", async () => {
        const initialWidth = canvas.style.width;
        const initialHeight = canvas.style.height;
        expect(initialWidth).toBe("612px");
        expect(initialHeight).toBe("792px");

        viewer.zoomOut();

        await new Promise((resolve) => setTimeout(resolve));

        const newCanvas = viewer.pageCanvases[0];
        const newWidth = newCanvas.style.width;
        const newHeight = newCanvas.style.height;

        expect(newWidth).not.toBe(initialWidth);
        expect(newWidth).toBe("489px");
        expect(newHeight).not.toBe(initialHeight);
        expect(newHeight).toBe("633px");
      });
    });
  });

  describe("Page nav button", () => {
    describe("Given the next page button is pressed", () => {
      test("the current page will change to the next one", async () => {
        expect(viewer.currentPage).toBe(1);

        viewer.nextPage();

        await new Promise((resolve) => setTimeout(resolve));
        expect(viewer.currentPage).toBe(2);
      });
    });

    describe("Given the previous page button is pressed on the second page", () => {
      test("the current page will change to the previous page", async () => {
        viewer.nextPage();

        expect(viewer.currentPage).toBe(2);

        viewer.prevPage();

        await new Promise((resolve) => setTimeout(resolve));
        expect(viewer.currentPage).toBe(1);
      });
    });

    describe("Given the previous page button is pressed on the first page", () => {
      test("the current page will not change", async () => {
        expect(viewer.currentPage).toBe(1);

        viewer.prevPage();

        await new Promise((resolve) => setTimeout(resolve));
        expect(viewer.currentPage).toBe(1);
      });
    });

    describe("Given the next page button is pressed on the last page", () => {
      test("the current page will not change to the next one", async () => {
        viewer.currentPage = 3;
        expect(viewer.currentPage).toBe(3);

        viewer.nextPage();

        await new Promise((resolve) => setTimeout(resolve));
        expect(viewer.currentPage).toBe(3);
      });
    });
  });

  describe("Given the thumbnail button is pressed", () => {
    test("the thumbnail panel will be visible", async () => {
      await viewer.toggleThumbnails();

      await new Promise((resolve) => setTimeout(resolve));

      expect(
        viewer.thumbnailPanel.classList.contains(
          "pdf-viewer-thumbnail-panel--visible",
        ),
      ).toBe(true);
    });
  });

  describe("Given the thumbnail button is pressed twice", () => {
    test("the thumbnail panel will not be visible", async () => {
      await viewer.toggleThumbnails();

      await new Promise((resolve) => setTimeout(resolve));

      await viewer.toggleThumbnails();

      await new Promise((resolve) => setTimeout(resolve));

      expect(
        viewer.thumbnailPanel.classList.contains(
          "pdf-viewer-thumbnail-panel--visible",
        ),
      ).toBe(false);
    });
  });
});
