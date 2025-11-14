package server

import "log"

// wasm_runner.go - placeholder; integrate your wasm engine here
type WasmRunner struct {
	bus *BusServer
}

func NewWasmRunner(bus *BusServer) *WasmRunner {
	return &WasmRunner{bus: bus}
}

func (wr *WasmRunner) LaunchWasm(appId string, wasmBytes []byte) error {
	// Implement sandboxing/loading using your preferred wasm runtime (wasmtime, wazero, etc.)
	log.Println("WASM launch requested for", appId)
	// on success, you might publish events to the bus
	return nil
}
