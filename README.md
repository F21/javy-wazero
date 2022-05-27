# Javy - Wazero example
## Compiling javascript to wasm
1. Install javy: `wget https://github.com/Shopify/javy/releases/download/v0.3.0/javy-x86_64-linux-v0.3.0.gz`
2. Unzip: `gzip -dc javy-x86_64-linux-v0.3.0.gz > /usr/local/bin/javy`
3. Compile: `javy -o js/greet.wasm js/greet.js`

## Running the example
To run: `go run main.go`

Note: Javy doesn't export any functions: the only way to invoke functions is via `_start`.
For this reason, a module needs to be instantiated each time you want to invoke anything.

See https://github.com/Shopify/javy/issues/111
