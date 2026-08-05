import { AddressFinder } from "./address-finder";

describe('capitalised letters test', () => {
    test("the address will be capitalised", async () => {
    let addressHtml = document.body
    let input = document.createElement('input');
    input.setAttribute('data-app-address-finder-map', 'addressLine2');
    input.setAttribute('name', 'addressLine2');
    addressHtml.appendChild(input);


    const addressFinder = new AddressFinder(addressHtml, { prefix: 'test' });
    addressFinder.underwriteValue('addressLine2', 'Address line tWo');
    expect(addressFinder.$module.querySelector('[name="addressLine2"]').value).toBe('ADDRESS LINE TWO');
    })
});
