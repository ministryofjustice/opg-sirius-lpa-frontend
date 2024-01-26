import tinymce from "tinymce/tinymce.min.js";
import "tinymce/icons/default";
import "tinymce/themes/silver";
import "tinymce/plugins/paste";
import "tinymce/plugins/lists";

const textEditor = () => {
  const prefix = document.body.getAttribute("data-prefix");

  tinymce.init({
    selector: "#documentTextEditor",
    menubar: false,
    statusbar: false,
    toolbar: "bold italic | bullist numlist",
    plugins: "paste lists",
    paste_as_text: true,
    paste_word_valid_elements: "h1,h2,h3,strong,em,b,i",
    browser_spellcheck: true,
    gecko_spellcheck: true,
    height: 300,
    cache_suffix: "?v=5.10.7",
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
