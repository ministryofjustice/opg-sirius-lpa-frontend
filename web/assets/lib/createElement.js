/**
 * @template {keyof HTMLElementTagNameMap} K
 * @param {K} tag
 * @param {{[name: string]: string}} attrs
 * @returns {HTMLElementTagNameMap[K]}
 */
export function createElement(tag, attrs = {}, $children = []) {
  const $element = document.createElement(tag);

  Object.entries(attrs).forEach(([key, value]) => {
    $element.setAttribute(key, value);
  });

  $children.forEach(($child) => {
    if ($child instanceof HTMLElement) {
      $element.appendChild($child);
    } else if (typeof $child === "string") {
      $element.appendChild(new Text($child));
    }
  });

  return $element;
}
