package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)

//go:embed wasm/greet.wasm
var greetWasm []byte

var results map[int32]chan []byte

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime()
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate WASI, which implements system I/O such as console output.
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		log.Panicln(err)
	}

	// Register the host functions that Suborbital Javy needs
	if err := registerHostFunctions(ctx, r); err != nil {
		log.Panicln(err)
	}

	// Compile and instantiate the module
	module, err := r.InstantiateModuleFromBinary(ctx, greetWasm)

	if err != nil {
		log.Panicln(err)
	}

	results = map[int32]chan []byte{} // Stores channels that we can use to send back the result

	for i := 0; i < 10; i++ {
		start := time.Now()

		input := fmt.Sprintf(`{"name": "Person %d"}`, i)
		inputLength := len(input)

		allocation, err := module.ExportedFunction("allocate").Call(ctx, uint64(inputLength)) // Allocate the memory

		if err != nil {
			log.Panicln(err)
		}

		inputPtr := allocation[0]

		defer module.ExportedFunction("deallocate").Call(ctx, inputPtr) // Remember to deallocate the memory after using it, otherwise it will leak!

		if !module.Memory().Write(ctx, uint32(inputPtr), []byte(input)) { // Write the input to memory
			log.Panicln(err)
		}

		ident, err := randomIdentifier() // Generate a random id for this call (Suborbital Javy expects this to correlate the input's caller with the output)

		if err != nil {
			log.Panicln(err)
		}

		resultCh := make(chan []byte, 1) // Must have a buffer of 1, otherwise it will block since we're reading and writing from the same goroutine

		results[ident] = resultCh

		defer delete(results, ident)

		_, err = module.ExportedFunction("run_e").Call(ctx, inputPtr, uint64(inputLength), uint64(ident)) // Call run_e with the location of the input and the identifier

		if err != nil {
			log.Panicln(err)
		}

		result := <-resultCh // Wait for the result

		fmt.Printf("%s", result)

		fmt.Println(" Time taken:", time.Since(start))
	}
}

// The Suborbital version of Javy expects these host functions to be implemented, but we only need "return_result"
func registerHostFunctions(ctx context.Context, r wazero.Runtime) error {
	_, err := r.NewModuleBuilder("env"). // They must be registered under "env"
						ExportFunction("return_result", func(ctx context.Context, m api.Module, ptr uint32, len uint32, ident uint32) {
			if ch, ok := results[int32(ident)]; ok {
				result, ok := m.Memory().Read(ctx, ptr, len) // Read the result written by the WebAssembly module

				if ok {
					ch <- result // Send it
				}
			}
		}).
		ExportFunction("get_static_file", func(_ uint32, _ uint32, _ uint32) uint32 {
			panic("get_static_file is unimplemented")
		}).
		ExportFunction("request_set_field", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("request_set_field is unimplemented")
		}).
		ExportFunction("cache_get", func(_ uint32, _ uint32, _ uint32) uint32 {
			panic("cache_get is unimplemented")
		}).
		ExportFunction("add_ffi_var", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("add_ffi_var is unimplemented")
		}).
		ExportFunction("get_ffi_result", func(_ uint32, _ uint32) uint32 {
			panic("get_ffi_result is unimplemented")
		}).
		ExportFunction("return_error", func(_ uint32, _ uint32, _ uint32, _ uint32) {
			panic("return_error is unimplemented")
		}).
		ExportFunction("fetch_url", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("fetch_url is unimplemented")
		}).
		ExportFunction("graphql_query", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("graphql_query is unimplemented")
		}).
		ExportFunction("db_exec", func(_ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("db_exec is unimplemented")
		}).
		ExportFunction("cache_set", func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("cache_set is unimplemented")
		}).
		ExportFunction("request_get_field", func(_ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("request_get_field is unimplemented")
		}).
		ExportFunction("log_msg", func(ctx context.Context, m api.Module, ptr uint32, size uint32, level uint32, ident uint32) {
			panic("log_msg is unimplemented")
		}).Instantiate(ctx, r)

	return err
}

func randomIdentifier() (int32, error) {
	// generate a random number between 0 and the largest possible int32
	num, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		return -1, fmt.Errorf("failed to generate random int: %w", err)
	}

	return int32(num.Int64()), nil
}
