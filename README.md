# OPA Management Application

This project aims to provide a simple and easy way to host your own OPA playground,
manage your policies and data, and test them against your own services.

## Features

- **OPA Playground**: A simple playground to test your policies and data.
- **Policy Management**: Manage your policies and data.
- **Service Testing**: Test your policies against your own services.
- **Decision Logs**: View the decision logs of your policies.
- **Decision Tree**: View the decision tree of your policies.

## Run OPA locally

```bash
opa run --server --config-file=./config.yaml --addr=localhost:8181 --diagnostic-addr=localhost:8282
```
