package uuidGen

import (
	"github.com/google/uuid"
)

func GenerateUUID() (uuid.UUID) {
	newUUID := uuid.New()
	return newUUID
}
