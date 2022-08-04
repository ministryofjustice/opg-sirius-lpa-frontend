const nunjucks = require('nunjucks');
const path = require('path');
const fs = require('fs');

var views = [
    path.join(__dirname, '/node_modules/govuk-frontend/'),
    path.join(__dirname, '/node_modules/govuk-frontend/components'),
    path.join(__dirname, '/node_modules/govuk_template_jinja/views/layouts'),
    path.join(__dirname, '/web/metatemplate'),
];

const nunjucksEnv = nunjucks.configure(views, {});

const s = nunjucksEnv.render('add_complaint.njk', {})
    .replace(/\{\!/g, '{{')
    .replace(/\!\}/g, '}}')
    .replace(/<div class="govuk-form-group govuk-form-group--error">([\s\S]+?)<p id="[a-z\- ]+" class="govuk-error-message">\s+<span class="govuk-visually-hidden">Error:<\/span>\s+([A-Za-z.]+)\s+<\/p>/gm, '<div class="govuk-form-group {{ if $2 }}govuk-form-group--error{{ end }}">$1{{ template "errors" $2 }}')

fs.writeFileSync(path.join(__dirname, '/web/template/add_complaint.gohtml'), `{{ define "page" }}${s}{{ end }}`);
