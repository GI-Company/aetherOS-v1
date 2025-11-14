package server

import (
	"github.com/google/uuid"
)

func genUUID() string {
	return uuid.NewString()
}
