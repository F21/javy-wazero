package main

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

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

	config := wazero.NewModuleConfig().WithStdout(os.Stdout)

	// Instantiate WASI, which implements system I/O such as console output.
	if _, err := wasi.InstantiateSnapshotPreview1(ctx, r); err != nil {
		log.Panicln(err)
	}

	compiled, err := r.CompileModule(ctx, greetWasm, wazero.NewCompileConfig())

	if err != nil {
		log.Panicln(err)
	}

	for i := 0; i < 100; i++ {
		input := bytes.NewBufferString(fmt.Sprintf(`{"name": "Person %d"}`, i))

		module, err := r.InstantiateModule(ctx, compiled, config.WithName(fmt.Sprintf("module-%d", i)).WithStdin(input))

		if err != nil {
			log.Panicln(err)
		}

		defer module.Close(nil)
	}
}
