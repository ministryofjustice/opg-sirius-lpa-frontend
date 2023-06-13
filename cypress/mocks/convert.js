const fs = require("fs");
const path = require("path");

const pact = JSON.parse(
  fs
    .readFileSync(
      path.join(__dirname, "../../pacts/sirius-lpa-frontend-sirius.json")
    )
    .toString()
);

const mappings = pact.interactions.map((interaction) => ({
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

fs.writeFileSync(
  path.join(__dirname, "sirius.json"),
  JSON.stringify({
    mappings,
  })
);
