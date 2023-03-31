package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func BenchmarkSuborbitalCallFunction(b *testing.B) {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
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
	module, err := r.Instantiate(ctx, greetWasm)
	if err != nil {
		log.Panicln(err)
	}

	input := fmt.Sprintf(`{"name": "Person"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err = callFunc(ctx, module, input, results); err != nil {
			b.Fatal(err)
		}
	}
}
