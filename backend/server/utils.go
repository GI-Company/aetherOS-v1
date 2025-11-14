
// ===========================
// backend/server/utils.go
// ===========================
package server

import (
	"github.com/google/uuid"
)

func genUUID() string {
	return uuid.New().String()
}
