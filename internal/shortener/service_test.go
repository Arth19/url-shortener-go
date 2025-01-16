package shortener

import (
	"testing"
)

func TestGenerateShortCode(t *testing.T) {
	code, err := GenerateShortCode()
	if err != nil {
		t.Fatalf("Erro ao gerar short code: %v", err)
	}

	if len(code) == 0 {
		t.Error("Short code retornou vazio, esperado algo não vazio")
	}
}
