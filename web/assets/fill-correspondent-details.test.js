import fillCorrespondentDetails from "./fill-correspondent-details";

describe("Fill correspondent details", () => {
  let radios;
  let fieldsContainer;
  let firstname;
  let lastname;

  beforeEach(async () => {
    document.body.innerHTML = `
    <form>
      <div data-module="app-fill-correspondent-details">
        <input type="radio" />
        <input type="radio" data-fields='{"firstname": "test", "lastname": "smith"}'/>
      </div>
      <input name="firstname"/>
      <input name="lastname"/>
    </form>
    <div class="govuk-!-display-none" data-module="app-fill-correspondent-details--fields"></div>
  `;

    radios = document.querySelectorAll(
      '[data-module="app-fill-correspondent-details"] input[type="radio"]',
    );
    fieldsContainer = document.querySelector(
      '[data-module="app-fill-correspondent-details--fields"]',
    );
    firstname = document.querySelector('input[name="firstname"]');
    lastname = document.querySelector('input[name="lastname"]');

    fillCorrespondentDetails();
  });

  afterEach(() => {
    radios = null;
    fieldsContainer = null;
    firstname = null;
    lastname = null;
    document.body.innerHTML = "";
  });

  test("it should show the fields container if it is hidden", async () => {
    expect(fieldsContainer.classList.contains("govuk-!-display-none")).toBe(
      true,
    );
    radios[0].dispatchEvent(new Event("change"));
    expect(fieldsContainer.classList.contains("govuk-!-display-none")).toBe(
      false,
    );
  });

  test("it should fill in the form inputs using the fields data attribute", async () => {
    expect(firstname.value).toBe("");
    expect(lastname.value).toBe("");

    radios[1].dispatchEvent(new Event("change"));
    expect(firstname.value).toBe("test");
    expect(lastname.value).toBe("smith");
  });

  test("it clear the form inputs if no fields data attribute is present", async () => {
    expect(firstname.value).toBe("");
    expect(lastname.value).toBe("");

    radios[1].dispatchEvent(new Event("change"));
    expect(firstname.value).toBe("test");
    expect(lastname.value).toBe("smith");

    radios[0].dispatchEvent(new Event("change"));
    expect(firstname.value).toBe("");
    expect(lastname.value).toBe("");
  });
});
