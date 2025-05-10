package storage

import (
	"path/filepath"

	"github.com/google/uuid"
)

func GenerateTeamLogoFileName(ext string) string {
	return uuid.New().String() + ext
}

func BuildTeamLogoPath(fileName string) string {
	return filepath.Join("uploads/logos", fileName)
}
