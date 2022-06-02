# Shopify Javy and Wazero

## Compile the WebAssembly module:
Shopify publishes pre-built binaries for Javy here: https://github.com/Shopify/javy/tags

We will use a docker container to compile the WebAssembly binary to avoid polluting the local environment.
1. From the root of the repository, run `docker compose run compile-shopify-javy-wasm`

## Run the example:
1. cd `shopify-javy`
2. `go run example.go`