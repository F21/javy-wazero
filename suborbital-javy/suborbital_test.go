package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)

func BenchmarkSuborbitalCallFunction(b *testing.B) {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime()
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate WASI, which implements system I/O such as console output.
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		log.Panicln(err)
	}

	results := &sync.Map{}

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

		input := fmt.Sprintf(`{"name": "Person %d"}`, i)

		b.Run(fmt.Sprintf("input-%d", i), func(b *testing.B) {
			_ = callFunc(ctx, module, input, results)
		})
	}
}
