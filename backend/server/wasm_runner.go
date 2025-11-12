package main

import (
	"context"
	"io"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

type WasmInstance struct {
	ctx context.Context
	// hold references to module and instance
	runtime wazero.Runtime
	module  api.Module
}

func NewWasmInstance(ctx context.Context, wasmBytes []byte) (*WasmInstance, error) {
	r := wazero.NewRuntime(ctx)

	_, err := r.NewHostModuleBuilder("host").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, m api.Module, stack []uint64) {}).
		Export("noop").
		Instantiate(ctx)

	if err != nil {
		return nil, err
	}

	mod, err := r.Instantiate(ctx, wasmBytes)
	if err != nil {
		return nil, err
	}

	return &WasmInstance{
		ctx:     ctx,
		runtime: r,
		module:  mod,
	}, nil
}

func (wi *WasmInstance) CallExport(name string, params ...uint64) (uint64, error) {
	fn := wi.module.ExportedFunction(name)
	if fn == nil {
		return 0, io.EOF
	}
	results, err := fn.Call(wi.ctx, params...)
	if err != nil {
		return 0, err
	}
	if len(results) > 0 {
		return results[0], nil
	}
	return 0, nil
}

func (wi *WasmInstance) Close() {
	if wi.module != nil {
		_ = wi.module.Close(wi.ctx)
	}
	if wi.runtime != nil {
		_ = wi.runtime.Close(wi.ctx)
	}
}
