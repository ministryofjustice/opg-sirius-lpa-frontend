import hugerte from "hugerte/hugerte.min.js";
import "hugerte/icons/default";
import "hugerte/themes/silver";
import "hugerte/plugins/lists";
import "hugerte/models/dom";

const textEditor = () => {
  const prefix = document.body.getAttribute("data-prefix");

  hugerte.init({
    selector: "#documentTextEditor",
    menubar: false,
    statusbar: false,
    toolbar: "bold italic | bullist numlist",
    plugins: "lists",
    paste_as_text: true,
    paste_word_valid_elements: "h1,h2,h3,strong,em,b,i",
    browser_spellcheck: true,
    gecko_spellcheck: true,
    height: 300,
    content_css: prefix + "/stylesheets/all.css",
    base_url: prefix + "/javascript",
    body_class: document.documentElement.classList.contains(
      "app-!-html-class--dark",
    )
      ? "app-!-html-class--dark"
      : "",
  });
};

export default textEditor;
