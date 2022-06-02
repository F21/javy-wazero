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
goos: linux
goarch: amd64
pkg: github.com/F21/javy-wazero/shopify-javy
cpu: Intel(R) Core(TM) i7-6700 CPU @ 3.40GHz
BenchmarkShopifyInstantiateModule/input-0-8                  404           3135726 ns/op
BenchmarkShopifyInstantiateModule/input-1-8                  378           2874987 ns/op
BenchmarkShopifyInstantiateModule/input-2-8                  421           2989359 ns/op
BenchmarkShopifyInstantiateModule/input-3-8                  430           2902166 ns/op
BenchmarkShopifyInstantiateModule/input-4-8                  417           2916718 ns/op
BenchmarkShopifyInstantiateModule/input-5-8                  397           2863885 ns/op
BenchmarkShopifyInstantiateModule/input-6-8                  409           2916866 ns/op
BenchmarkShopifyInstantiateModule/input-7-8                  418           2893322 ns/op
BenchmarkShopifyInstantiateModule/input-8-8                  416           2908626 ns/op
BenchmarkShopifyInstantiateModule/input-9-8                  403           2946496 ns/op
PASS
ok      github.com/F21/javy-wazero/shopify-javy 16.068s
goos: linux
goarch: amd64
pkg: github.com/F21/javy-wazero/suborbital-javy
cpu: Intel(R) Core(TM) i7-6700 CPU @ 3.40GHz
BenchmarkSuborbitalCallFunction/input-0-8               1000000000               0.0001566 ns/op
BenchmarkSuborbitalCallFunction/input-1-8               1000000000               0.0001420 ns/op
BenchmarkSuborbitalCallFunction/input-2-8               1000000000               0.0002191 ns/op
BenchmarkSuborbitalCallFunction/input-3-8               1000000000               0.0001599 ns/op
BenchmarkSuborbitalCallFunction/input-4-8               1000000000               0.0002473 ns/op
BenchmarkSuborbitalCallFunction/input-5-8               1000000000               0.0002135 ns/op
BenchmarkSuborbitalCallFunction/input-6-8               1000000000               0.0002033 ns/op
BenchmarkSuborbitalCallFunction/input-7-8               1000000000               0.0002502 ns/op
BenchmarkSuborbitalCallFunction/input-8-8               1000000000               0.0001736 ns/op
BenchmarkSuborbitalCallFunction/input-9-8               1000000000               0.0001849 ns/op
PASS
ok      github.com/F21/javy-wazero/suborbital-javy      1.019s
```