package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"testing"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/wasi_snapshot_preview1"
)

func BenchmarkShopifyInstantiateModule(b *testing.B) {
	ctx := context.Background()

	// Create a new WebAssembly Runtime.
	r := wazero.NewRuntime()
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

	for i := 0; i < 10; i++ {

		inputJSON := fmt.Sprintf(`{"name": "Person %d"}`, i)

		b.Run(fmt.Sprintf("input-%d", i), func(b *testing.B) {

			for j := 0; j < b.N; j++ {
				id, err := randomIdentifier()

				if err != nil {
					b.Fatal(err)
				}

				input := strings.NewReader(inputJSON)

				config := wazero.NewModuleConfig().WithName(strconv.Itoa(int(id))).WithStdin(input)

				startModule(ctx, r, compiled, config)
			}
		})
	}
}

func randomIdentifier() (int32, error) {
	// generate a random number between 0 and the largest possible int32
	num, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		return -1, fmt.Errorf("failed to generate random int: %w", err)
	}

	return int32(num.Int64()), nil
}
