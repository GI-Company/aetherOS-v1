package main

import "github.com/google/uuid"

func uuidSafe() string {
	return uuid.NewString()
}
