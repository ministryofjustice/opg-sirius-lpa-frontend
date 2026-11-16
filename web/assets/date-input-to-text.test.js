import dateInputToText from "./date-input-to-text";

describe("Convert date input to text", () => {
  let dateInput;
  let form;
  const submitMock = jest.fn();

  beforeEach(async () => {
    document.body.innerHTML = `
    <form>
      <input type="date" data-module="app-date-input-to-text" value="2026-11-16">
    </form>
  `;

    dateInput = document.querySelector(
      'input[type="date"][data-module="app-date-input-to-text"]',
    );

    form = document.querySelector("form");
    form.onsubmit = submitMock;

    dateInputToText();
  });

  afterEach(() => {
    dateInput = null;
    form = null;
    document.body.innerHTML = "";
  });

  test("it should convert the input type to text and format the value correctly on double click", async () => {
    expect(dateInput.type).toBe("date");
    dateInput.dispatchEvent(new Event("dblclick"));
    expect(dateInput.type).toBe("text");
    expect(dateInput.value).toBe("16/11/2026");
  });

  test("it should convert the input type back to text and format the value correctly when losing focus", async () => {
    dateInput.dispatchEvent(new Event("dblclick"));
    dateInput.dispatchEvent(new Event("focusout"));
    expect(dateInput.type).toBe("date");
    expect(dateInput.value).toBe("2026-11-16");
  });

  test("it should not submit the form when pressing enter while the input is in text mode", async () => {
    dateInput.dispatchEvent(new Event("dblclick"));
    dateInput.dispatchEvent(new KeyboardEvent("keydown", { key: "Enter" }));
    expect(submitMock).not.toHaveBeenCalled();
    expect(dateInput.type).toBe("date");
    expect(dateInput.value).toBe("2026-11-16");
  });
});
