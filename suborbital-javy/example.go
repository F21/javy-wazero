package main

import (
	"context"
	"crypto/rand"
	_ "embed"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

//go:embed wasm/greet.wasm
var greetWasm []byte

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate WASI, which implements system I/O such as console output.
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		log.Panicln(err)
	}

	results := &sync.Map{} // Stores channels that we can use to send back the result

	// Register the host functions that Suborbital Javy needs
	if err := registerHostFunctions(ctx, r, results); err != nil {
		log.Panicln(err)
	}

	// Compile and instantiate the module
	module, err := r.InstantiateModuleFromBinary(ctx, greetWasm)
	if err != nil {
		log.Panicln(err)
	}

	for i := 0; i < 10; i++ {
		start := time.Now()

		input := fmt.Sprintf(`{"name": "Person %d"}`, i)

		result, err := callFunc(ctx, module, input, results)
		if err != nil {
			log.Panicln(err)
		}

		fmt.Printf("%s", result)

		fmt.Println(" Time taken:", time.Since(start))
	}
}

func callFunc(ctx context.Context, module api.Module, input string, results *sync.Map) ([]byte, error) {
	inputLength := len(input)

	allocation, err := module.ExportedFunction("allocate").Call(ctx, uint64(inputLength)) // Allocate the memory
	if err != nil {
		return nil, err
	}

	inputPtr := allocation[0]

	defer module.ExportedFunction("deallocate").Call(ctx, inputPtr) // Remember to deallocate the memory after using it, otherwise it will leak!

	if !module.Memory().Write(ctx, uint32(inputPtr), []byte(input)) { // Write the input to memory
		return nil, err
	}

	ident, err := randomIdentifier() // Generate a random id for this call (Suborbital Javy expects this to correlate the input's caller with the output)
	if err != nil {
		return nil, err
	}

	resultCh := make(chan []byte, 1) // Must have a buffer of 1, otherwise it will block since we're reading and writing from the same goroutine

	results.Store(ident, resultCh)

	defer results.Delete(ident)

	_, err = module.ExportedFunction("run_e").Call(ctx, inputPtr, uint64(inputLength), uint64(ident)) // Call run_e with the location of the input and the identifier

	if err != nil {
		return nil, err
	}

	result := <-resultCh // Wait for the result

	return result, nil
}

// The Suborbital version of Javy expects these host functions to be implemented, but we only need "return_result"
func registerHostFunctions(ctx context.Context, r wazero.Runtime, results *sync.Map) error {
	_, err := r.NewHostModuleBuilder("env"). // They must be registered under "env"
							NewFunctionBuilder().
							WithFunc(func(ctx context.Context, m api.Module, ptr uint32, len uint32, ident uint32) {
			if ch, ok := results.Load(int32(ident)); ok {

				resultCh, ok := ch.(chan []byte)

				if !ok {
					log.Panicln("Channel is not the right type")
				}

				result, ok := m.Memory().Read(ctx, ptr, len) // Read the result written by the WebAssembly module

				if ok {
					resultCh <- result // Send it
				}
			}
		}).
		Export("return_result").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32) uint32 {
			panic("get_static_file is unimplemented")
		}).
		Export("get_static_file").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("request_set_field is unimplemented")
		}).
		Export("request_set_field").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32) uint32 {
			panic("cache_get is unimplemented")
		}).
		Export("cache_get").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("add_ffi_var is unimplemented")
		}).
		Export("add_ffi_var").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32) uint32 {
			panic("get_ffi_result is unimplemented")
		}).
		Export("get_ffi_result").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32, _ uint32) {
			panic("return_error is unimplemented")
		}).
		Export("return_error").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("fetch_url is unimplemented")
		}).
		Export("fetch_url").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("graphql_query is unimplemented")
		}).
		Export("graphql_query").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("db_exec is unimplemented")
		}).
		Export("db_exec").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("cache_set is unimplemented")
		}).
		Export("cache_set").
		NewFunctionBuilder().
		WithFunc(func(_ uint32, _ uint32, _ uint32, _ uint32) uint32 {
			panic("request_get_field is unimplemented")
		}).
		Export("request_get_field").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, ptr uint32, size uint32, level uint32, ident uint32) {
			panic("log_msg is unimplemented")
		}).
		Export("log_msg").
		Instantiate(ctx, r)

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
