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
	"github.com/tetratelabs/wazero/wasi"
)

//go:embed js/greet.wasm
var greetWasm []byte

func main() {
	// Choose the context to use for function calls.
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime()
	defer r.Close(ctx) // This closes everything this Runtime created.

	// Instantiate WASI, which implements system I/O such as console output.
	if _, err := wasi.InstantiateSnapshotPreview1(ctx, r); err != nil {
		log.Panicln(err)
	}

	// Compile the module to reduce performance impact of instantiating it multiple times.
	compiled, err := r.CompileModule(ctx, greetWasm, wazero.NewCompileConfig())
	if err != nil {
		log.Panicln(err)
	}

	config := wazero.NewModuleConfig().WithStdout(os.Stdout)
	for i := 0; i < 100; i++ {
		start := time.Now()

		// The only way to pass data to Javy is via Stdin, arguments or env variables.
		input := strings.NewReader(fmt.Sprintf(`{"name": "Person %d"}`, i))
		startModule(ctx, r, compiled, config.WithStdin(input))

		fmt.Println("time taken", time.Since(start))
	}
}

// startModule implicitly invokes the "_start" function defined by Javy.
// Note: If you want to run modules in parallel, config must override the name.
func startModule(ctx context.Context, r wazero.Runtime, compiled wazero.CompiledModule, config wazero.ModuleConfig) {
	module, err := r.InstantiateModule(ctx, compiled, config)
	if err != nil {
		log.Panicln(err)
	}
	defer module.Close(ctx)
}
