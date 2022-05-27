# Javy - Wazero example

To run: `go run main.go`

Note: Javy doesn't export any functions: the only way to invoke functions is via `_start`.
For this reason, a module needs to be instantiated each time you want to invoke anything.

See https://github.com/Shopify/javy/issues/111
