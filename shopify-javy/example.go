package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tetratelabs/wazero"
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

	// Compile the module to reduce performance impact of instantiating it multiple times.
	compiled, err := r.CompileModule(ctx, greetWasm, wazero.NewCompileConfig())
	if err != nil {
		log.Panicln(err)
	}

	// Set stdout
	config := wazero.NewModuleConfig().WithStdout(os.Stdout)

	for i := 0; i < 10; i++ {
		start := time.Now()

		// The only way to pass data to Javy is via Stdin, arguments or env variables as JSON.
		input := strings.NewReader(fmt.Sprintf(`{"name": "Person %d"}`, i))
		err = startModule(ctx, r, compiled, config.WithStdin(input)) // Set stdin
		if err != nil {
			log.Panicln(err)
		}

		fmt.Println(" Time taken:", time.Since(start))
	}
}

// startModule implicitly invokes the "_start" function defined by Javy.
// Note: If you want to run modules in parallel, config must override the name.
func startModule(ctx context.Context, r wazero.Runtime, compiled wazero.CompiledModule, config wazero.ModuleConfig) error {
	if mod, err := r.InstantiateModule(ctx, compiled, config); err != nil {
		return err
	} else {
		return mod.Close(ctx)
	}
}
