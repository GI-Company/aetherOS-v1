package aetherscript

import (
	"fmt"
)

func Compile(node Node) ([]byte, error) {
	var bytecode []byte

	switch n := node.(type) {
	case ListNode:
		if len(n.Children) == 0 {
			return nil, fmt.Errorf("empty list")
		}

		symbol, ok := n.Children[0].(SymbolNode)
		if !ok {
			return nil, fmt.Errorf("expected symbol as first element of list")
		}

		for i := 1; i < len(n.Children); i++ {
			childCode, err := Compile(n.Children[i])
			if err != nil {
				return nil, err
			}
			bytecode = append(bytecode, childCode...)
		}

		bytecode = append(bytecode, OpCall)
		bytecode = append(bytecode, []byte(symbol.Value)...)
		bytecode = append(bytecode, 0) // Null terminator for symbol

	case NumberNode:
		bytecode = append(bytecode, OpPush)
		// A more robust implementation would handle different number types
		bytecode = append(bytecode, []byte(fmt.Sprintf("%f", n.Value))...)
		bytecode = append(bytecode, 0) // Null terminator

	case StringNode:
		bytecode = append(bytecode, OpPush)
		bytecode = append(bytecode, []byte(n.Value)...)
		bytecode = append(bytecode, 0) // Null terminator

	case SymbolNode:
		bytecode = append(bytecode, OpLoad)
		bytecode = append(bytecode, []byte(n.Value)...)
		bytecode = append(bytecode, 0) // Null terminator
	default:
		return nil, fmt.Errorf("unknown node type: %T", n)
	}

	return bytecode, nil
}
