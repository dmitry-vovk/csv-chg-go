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
