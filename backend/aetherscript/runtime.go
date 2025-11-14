package aetherscript

import (
	"fmt"
)

func Execute(code []byte) error {
	stack := make([]interface{}, 0)
	variables := make(map[string]interface{})

	ip := 0
	for ip < len(code) {
		op := code[ip]
		ip++
		switch op {
		case OpPush:
			val, n := readNullTerminated(code, ip)
			ip += n
			stack = append(stack, val)
		case OpCall:
			funcName, n := readNullTerminated(code, ip)
			ip += n
			// In a real implementation, you would call the function
			// For now, we'll just print the function name and arguments
			fmt.Printf("Function call: %s, args: %v\n", funcName, stack)
			stack = make([]interface{}, 0) // Clear the stack after function call
		case OpStore:
			varName, n := readNullTerminated(code, ip)
			ip += n
			if len(stack) > 0 {
				variables[varName] = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			} else {
				return fmt.Errorf("stack is empty, cannot store variable")
			}
		case OpLoad:
			varName, n := readNullTerminated(code, ip)
			ip += n
			if val, ok := variables[varName]; ok {
				stack = append(stack, val)
			} else {
				return fmt.Errorf("undefined variable: %s", varName)
			}
		default:
			return fmt.Errorf("unknown opcode: %d", op)
		}
	}
	return nil
}

func readNullTerminated(data []byte, offset int) (string, int) {
	end := offset
	for end < len(data) && data[end] != 0 {
		end++
	}
	return string(data[offset:end]), (end - offset) + 1
}
