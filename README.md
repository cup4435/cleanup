# Stock Analytics Cleanup

A set of stock analytics microservices using Go and gRPC.

## Services

- **API Gateway**: HTTP REST API on port 8080
- **Stock Proxy**: gRPC service on port 50051

## Dev Environment

### VSCode

```bash
# Open dev environment in VSCode
code .

# Reopen in devcontainer
```

### Neovim

```bash
# Start dev environment
docker-compose -f docker-compose.dev.yml up -d --build

# Connect to dev environment (e.g. api-gateway)
docker-compose -f docker-compose.dev.yml exec api-gateway bash

# Open Neovim
nvim .
```

## Prod Environment

```bash
# Start production environment
docker-compose -f docker-compose.prod.yml up -d --build
```
