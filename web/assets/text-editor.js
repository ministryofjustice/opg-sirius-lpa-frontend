/* Import TinyMCE */
import tinymce from "tinymce/tinymce.min.js";

/* Default icons are required for TinyMCE 5.3 or above */
// import 'tinymce/icons/default';

/* A theme is also required */
// import 'tinymce/themes/modern/theme';

/* Initialize TinyMCE */
const textEditor = () => {
    tinymce.init({
        selector: '#documentTextEditor',
        menubar: false,
        statusbar: false,
        toolbar: 'bold italic | bullist numlist',
        plugins: 'paste lists',
        paste_as_text: true,
        browser_spellcheck: true,
        height: 300,
        cache_suffix: '?v=5.6.1',
    });
}

export default textEditor;

