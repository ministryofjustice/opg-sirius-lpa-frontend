# Playwright tests

These are intended to ultimately replate the Cypress tests, but this is not a quick process, so for the duration of the migration both will need to be maintained.

## Instructions

To run locally ensure Sirius is started and then run headless:

```sh
make run-playwright
```

To run with the UI:

```sh
make test-ui
```