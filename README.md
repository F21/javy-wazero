# Javy - Wazero example

There are currently 2 flavors of Javy:
- [Shopify's original version](https://github.com/Shopify/javy): Only accepts input via stdin and only produces output via stdout
- [Suborbital's fork](https://github.com/suborbital/javy): Reads and writes to memory (exports `allocate` and `deallocate`)

This repository provides an example written for both versions
- [Shopify Javy](shopify-javy)
- [Suborbital Javy](suborbital-javy)

## Why Javy?
On the surface, it might seem quite strange to have a compiler that compiles JavaScript to WebAssembly, given that we
can execute JavaScript directly in the browser and Node.js. However, Javy unlocks interesting use-cases for running 
JavaScript libraries in any WebAssembly environment.

As an example, to compile [MJML](https://mjml.io/) into HTML, one would need to use the [MJML library](https://github.com/mjmlio/mjml),
which is written in JavaScript. By writing a thin wrapper around the library and compiling it into WebAssembly using
Javy, we can now embed the WebAssembly module in any application that has a WebAssembly runtime. We recently released
[mjml-go](https://github.com/Boostport/mjml-go) which does exactly this. Please check it out to see how to use Wazero
beyond the simple examples in this repository.

## Why Wazero?
[Wazero](https://github.com/tetratelabs/wazero) is a pure-Go WebAssembly runtime. It does not require CGO nor does it
need any external dependencies. Another benefit is that the API is well-though-out and easy to grok. 

## Benchmarks
To run the benchmarks, run `go test -bench=. ./...` from the root of the repository.

Here are the results on my machine:
```
goos: darwin
goarch: amd64
pkg: github.com/F21/javy-wazero/shopify-javy
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkShopifyInstantiateModule-16    	     597	   2099444 ns/op
PASS
ok  	github.com/F21/javy-wazero/shopify-javy	4.342s
goos: darwin
goarch: amd64
pkg: github.com/F21/javy-wazero/suborbital-javy
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkSuborbitalCallFunction-16    	   13616	     85742 ns/op
PASS
ok  	github.com/F21/javy-wazero/suborbital-javy	5.490s
```