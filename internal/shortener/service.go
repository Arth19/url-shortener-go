package shortener

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func GenerateShortCode() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	code := base64.URLEncoding.EncodeToString(b)
	code = strings.TrimRight(code, "=")
	code = strings.TrimRight(code, "-_")
	code = strings.ReplaceAll(code, "-", "")
	code = strings.ReplaceAll(code, "_", "")

	return code, nil
}
