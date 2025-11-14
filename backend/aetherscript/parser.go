package aetherscript

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Node interface{}

type NumberNode struct{
	Value float64
}

type StringNode struct{
	Value string
}

type SymbolNode struct{
	Value string
}

type ListNode struct{
	Children []Node
}

func Parse(code string) (Node, error) {
	code = strings.ReplaceAll(code, "(", " ( ")
	code = strings.ReplaceAll(code, ")", " ) ")
	tokens := strings.Fields(code)
	
	if len(tokens) == 0 {
		return nil, errors.New("empty code")
	}

	pos := 0
	return parseTokens(tokens, &pos)
}

func parseTokens(tokens []string, pos *int) (Node, error) {
	if *pos >= len(tokens) {
		return nil, errors.New("unexpected end of input")
	}

	token := tokens[*pos]
	*pos++

	if token == "(" {
		var children []Node
		for *pos < len(tokens) && tokens[*pos] != ")" {
			child, err := parseTokens(tokens, pos)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		}
		if *pos >= len(tokens) || tokens[*pos] != ")" {
			return nil, errors.New("missing closing parenthesis")
		}
		*pos++ // Consume the closing parenthesis
		return ListNode{Children: children}, nil
	} else if token == ")" {
		return nil, errors.New("unexpected closing parenthesis")
	} else if val, err := strconv.ParseFloat(token, 64); err == nil {
		return NumberNode{Value: val}, nil
	} else if strings.HasPrefix(token, "\"") && strings.HasSuffix(token, "\"") {
		return StringNode{Value: token[1 : len(token)-1]}, nil
	} else {
		return SymbolNode{Value: token}, nil
	}
}
