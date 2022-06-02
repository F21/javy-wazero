# Suborbital Javy and Wazero

## Compile the WebAssembly module:
Suborbital does not publish any prebuilt binaries, however, they publish a docker image here: https://hub.docker.com/r/suborbital/builder-js

Use docker to compile the WebAssembly binary to avoid polluting the local environment.
1. From the root of the repository, run `docker compose run compile-suborbital-javy-wasm`

Notes:
As `greet.js` imports various modules such as the Suborbital runtime, we need to bundle everything into one JavaScript
file using webpack. Javy expects one JavaScript file, so this step is crucial. 

## Run the example:
1. cd `suborbital-javy`
2. `go run example.go`