import dateInputToText from "./date-input-to-text";

Object.assign(navigator, {
  clipboard: {
    writeText: jest.fn().mockImplementation(() => Promise.resolve()),
  },
});

describe("Convert date input to text", () => {
  let dateInput;

  beforeEach(async () => {
    document.body.innerHTML = `
    <form>
      <input type="date" data-module="app-date-input-to-text" value="2026-11-16">
    </form>
  `;

    dateInput = document.querySelector(
      'input[type="date"][data-module="app-date-input-to-text"]',
    );

    dateInputToText();
  });

  afterEach(() => {
    dateInput = null;
    document.body.innerHTML = "";
    jest.clearAllMocks();
  });

  test("it should convert the input type to text and format the value correctly on double click", async () => {
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(new Event("dblclick"));
    expect(dateInput.type).toBe("text");
    expect(dateInput.value).toBe("16/11/2026");
  });

  test("it should convert the input type to text and format the value correctly when pressing control + v", async () => {
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(
      new KeyboardEvent("keydown", { key: "V", ctrlKey: true }),
    );
    expect(dateInput.type).toBe("text");
    expect(dateInput.value).toBe("16/11/2026");
  });

  test("it should convert the input type to text and format the value correctly when pressing meta + v", async () => {
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(
      new KeyboardEvent("keydown", { key: "V", metaKey: true }),
    );
    expect(dateInput.type).toBe("text");
    expect(dateInput.value).toBe("16/11/2026");
  });

  test("it should convert the input type back to date and format the value correctly when losing focus", async () => {
    dateInput.dispatchEvent(new Event("dblclick"));
    dateInput.dispatchEvent(new Event("focusout"));
    expect(dateInput.type).toBe("date");
    expect(dateInput.value).toBe("2026-11-16");
  });

  test("it should convert the input type back to date and format the value correctly pressing enter", async () => {
    dateInput.dispatchEvent(new Event("dblclick"));
    const event = new KeyboardEvent("keydown", {
      key: "Enter",
      cancelable: true,
    });
    dateInput.dispatchEvent(event);
    expect(event.defaultPrevented).toBe(true);
    expect(dateInput.type).toBe("date");
    expect(dateInput.value).toBe("2026-11-16");
  });

  test("it should restore to the original value if the new value is not a valid date", async () => {
    dateInput.dispatchEvent(new Event("dblclick"));
    dateInput.value = "invalid";
    dateInput.dispatchEvent(new Event("focusout"));
    expect(dateInput.value).toBe("2026-11-16");
  });

  test("it should copy the date input value to the clipboard when pressing control + c", async () => {
    dateInput.dispatchEvent(
      new KeyboardEvent("keydown", { key: "C", ctrlKey: true }),
    );
    expect(navigator.clipboard.writeText).toHaveBeenCalledWith("16/11/2026");
    dateInput.type = "date";
  });

  test("it should copy the date input value to the clipboard when pressing meta + c", async () => {
    dateInput.dispatchEvent(
      new KeyboardEvent("keydown", { key: "C", metaKey: true }),
    );
    expect(navigator.clipboard.writeText).toHaveBeenCalledWith("16/11/2026");
    dateInput.type = "date";
  });

  test("it should not change the field type if the control, meta, c and v keys are pressed independently", async () => {
    dateInput.dispatchEvent(
      new KeyboardEvent("keydown", { key: "C", ctrlKey: false }),
    );
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(
      new KeyboardEvent("keydown", { key: "C", metaKey: false }),
    );
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(
      new KeyboardEvent("keydown", { key: "V", ctrlKey: false }),
    );
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(
      new KeyboardEvent("keydown", { key: "V", metaKey: false }),
    );
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(new KeyboardEvent("keydown", { key: "Control" }));
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(new KeyboardEvent("keydown", { key: "Meta" }));
    expect(dateInput.type).toBe("date");
  });
});
