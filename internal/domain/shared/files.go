package shared

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

const (
	MAX_FILE_SIZE_500_KB = 500 * 1024
	MAX_FILE_SIZE        = MAX_FILE_SIZE_500_KB
)

var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".svg":  true,
}

func SaveBase64File(baseDir, base64Data string) (string, error) {
	if base64Data == "" {
		return "", nil
	}

	decoded, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", ErrInvalidGoldenPosterData
	}

	if len(decoded) > MAX_FILE_SIZE {
		return "", ErrInvalidGoldenPosterData
	}

	extension := detectFileExtension(decoded)
	if !AllowedExtensions[extension] {
		return "", ErrInvalidGoldenPosterData
	}

	filename := UUIDv4() + extension
	cleanBaseDir := strings.TrimSuffix(baseDir, "/")
	fullPath := filepath.Join(cleanBaseDir, filename)

	if err := os.MkdirAll(cleanBaseDir, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(fullPath, decoded, 0600); err != nil {
		return "", err
	}

	return filename, nil
}

func detectFileExtension(data []byte) string {
	if len(data) < 8 {
		return ""
	}

	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return ".jpg"
	}

	if len(data) >= 8 && data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47 {
		return ".png"
	}

	if len(data) >= 5 && string(data[0:5]) == "<?xml" {
		return ".svg"
	}

	if len(data) >= 4 && string(data[0:4]) == "<svg" {
		return ".svg"
	}

	return ""
}
