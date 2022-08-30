package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)

func BenchmarkShopifyInstantiateModule(b *testing.B) {
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate WASI, which implements system I/O such as console output.
	if _, err := wasi_snapshot_preview1.Instantiate(ctx, r); err != nil {
		b.Fatal(err)
	}

	// Compile the module to reduce performance impact of instantiating it multiple times.
	compiled, err := r.CompileModule(ctx, greetWasm, wazero.NewCompileConfig())
	if err != nil {
		b.Fatal(err)
	}

	inputJSON := fmt.Sprintf(`{"name": "Person"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input := strings.NewReader(inputJSON)

		config := wazero.NewModuleConfig().WithName(strconv.Itoa(i)).WithStdin(input)
		if err := startModule(ctx, r, compiled, config); err != nil {
			b.Fatal(err)
		}
	}
}
