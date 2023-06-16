const fs = require("fs");
const path = require("path");

const toMappings = (filename) => {
  const pact = JSON.parse(
    fs
      .readFileSync(path.join(__dirname, `../../pacts/${filename}.json`))
      .toString()
  );

  return pact.interactions.map((interaction) => ({
    name: interaction.description,
    request: {
      method: interaction.request.method,
      url:
        interaction.request.path +
        (interaction.request.query ? `?${interaction.request.query}` : ""),
    },
    response: {
      status: interaction.response.status,
      headers: interaction.response.headers,
      body: JSON.stringify(interaction.response.body),
    },
  }));
};

const mappings = [
  ...toMappings("sirius-lpa-frontend-sirius"),
  ...toMappings("ignored-ignored"),
];

fs.writeFileSync(
  path.join(__dirname, "migrated-from-pact.json"),
  JSON.stringify({
    mappings,
  })
);
