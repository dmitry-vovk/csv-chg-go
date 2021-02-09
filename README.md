![coverage-badge-do-not-edit](https://img.shields.io/badge/Coverage-97%25-brightgreen.svg?longCache=true&style=flat)
![Test](https://github.com/dmitry-vovk/csv-chg-go/workflows/Test/badge.svg)

# Divido Code Challenge with Go

> Warehouse Stocks Checker

## Intro

We want the applicant to exclusively focus on coding at this stage. For this reason we provide a complete suite with:
- Makefile targets
- Github Actions setup

The Github Action workflow is triggered on PRs or pushes to master, runs tests and checks whether the 
coverage percentage is greater or equal than 80%, otherwise it fails. You can achieve the same locally running `make`.

## Requirements

The task for this challenge is to develop a simple worker that, given in input items UUIDs (from a CSV file), will 
retrieve the quantity of stocks from a REST endpoint in the warehouse API and will create a "low stock alert" by calling
a different endpoint in case the number is below 5.

```
for UUID in CSV file
    response = call `GET /item/{UUID}` from Warehouse API
    if response.quantity < 5
        call `POST /low-stock-alert/{UUID}`
```

The UUIDs in input are to be retrieved from a CSV file that has the following structure:

```csv
767d967f-b55b-4457-bfee-685eaa6d0583
ee88ff32-f753-4a49-abf1-2885fdfcafba
9e2cb4dd-bd6e-48aa-9c0d-696a058226ed
...
```

[Here](/.divido/warehouse-api-specs.yml) you can find the OpenAPIv3 specification for the Warehouse API. We expect these
endpoints to be mocked based on the tests you are writing.

## Nice to have

- High performance for speed and memory usage
- Support graceful shutdown for the worker

## How to run

The app accepts command line arguments:
 * `-api <address>` (required) -- base URL of warehouse API, e.g. `https://api.warehouse.tld/v1`
 * `-input <source>` (required) -- source CSV, can be either local file path, or URL. Also, can be omitted if the last command line argument is `--`, in this case the app will read input from `stdin`. 
 * `-interval 60s` (optional) -- interval between request runs. The format should be supported by `time.ParseDelay()` function.
 * `-workers 1` (optional) -- number of parallel API requests to make.

Dockerfile can be found in the repository root that will run the app.
