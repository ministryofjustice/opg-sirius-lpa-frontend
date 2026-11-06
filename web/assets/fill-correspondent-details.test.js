import fillCorrespondentDetails from "./fill-correspondent-details";

describe("Fill correspondent details", () => {
  let radio;
  let firstname;
  let lastname;

  beforeEach(async () => {
    document.body.innerHTML = `
    <form>
      <div data-module="app-fill-correspondent-details">
        <input type="radio" data-fields='{"firstname": "test", "lastname": "smith"}'/>
      </div>
      <input name="firstname"/>
      <input name="lastname"/>
    </form>
  `;

    radio = document.querySelector(
      '[data-module="app-fill-correspondent-details"] input[type="radio"]',
    );
    firstname = document.querySelector('input[name="firstname"]');
    lastname = document.querySelector('input[name="lastname"]');

    fillCorrespondentDetails();
  });

  afterEach(() => {
    radio = null;
    firstname = null;
    lastname = null;
    document.body.innerHTML = "";
  });

  test("it should fill in the form inputs using the fields data attribute", async () => {
    expect(firstname.value).toBe("");
    expect(lastname.value).toBe("");

    radio.dispatchEvent(new Event("change"));

    expect(firstname.value).toBe("test");
    expect(lastname.value).toBe("smith");
  });
});
