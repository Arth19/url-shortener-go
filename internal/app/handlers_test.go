package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Arth19/url-shortener-go/internal/storage"
	"github.com/gorilla/mux"
)

func clearTable(t *testing.T) {
	if err := storage.DB.Exec("DELETE FROM urls").Error; err != nil {
		t.Fatalf("Erro ao limpar tabela urls: %v", err)
	}
}

func createTestServer(t *testing.T) *httptest.Server {
	err := storage.InitPostgres()
	if err != nil {
		t.Fatalf("Erro ao inicializar Postgres nos testes: %v", err)
	}

	clearTable(t)

	r := mux.NewRouter()
	SetupRoutes(r)
	ts := httptest.NewServer(r)
	return ts
}

func TestEncurtarURLHandler(t *testing.T) {
	ts := createTestServer(t)
	defer ts.Close()

	body := []byte(`{"url":"https://www.google.com"}`)

	resp, err := http.Post(ts.URL+"/encurtar", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Erro ao fazer POST /encurtar: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status esperado 200, obteve %d", resp.StatusCode)
	}

	var result struct {
		ShortCode string `json:"short_code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Erro ao decodificar JSON: %v", err)
	}

	if result.ShortCode == "" {
		t.Error("Esperado um short_code não vazio, mas veio vazio")
	}
}

func TestRedirecionarHandler(t *testing.T) {
	ts := createTestServer(t)
	defer ts.Close()

	if _, err := storage.SaveURL("teste", "https://www.google.com"); err != nil {
		t.Fatalf("Erro ao salvar shortCode: %v", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(ts.URL + "/teste")
	if err != nil {
		t.Fatalf("Erro ao fazer GET: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("Esperado 302 Found, obteve %d", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	if loc != "https://www.google.com" {
		t.Errorf("Esperado Location=https://www.google.com, obtido %s", loc)
	}
}

func TestStatsHandler(t *testing.T) {
	ts := createTestServer(t)
	defer ts.Close()

	body := []byte(`{"url":"https://www.google.com"}`)
	resp, err := http.Post(ts.URL+"/encurtar", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Erro ao POST /encurtar: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Esperado status 200 ao encurtar, obtido %d", resp.StatusCode)
	}

	var encurtarResp struct {
		ShortCode string `json:"short_code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&encurtarResp); err != nil {
		t.Fatalf("Erro ao decodificar JSON: %v", err)
	}

	_, _ = http.Get(ts.URL + "/" + encurtarResp.ShortCode)
	_, _ = http.Get(ts.URL + "/" + encurtarResp.ShortCode)

	statsURL := ts.URL + "/stats/" + encurtarResp.ShortCode
	statsResp, err := http.Get(statsURL)
	if err != nil {
		t.Fatalf("Erro ao GET /stats/%s: %v", encurtarResp.ShortCode, err)
	}
	defer statsResp.Body.Close()

	if statsResp.StatusCode != http.StatusOK {
		t.Fatalf("Esperado status 200 em /stats, obtido %d", statsResp.StatusCode)
	}

	var statsBody struct {
		ShortCode  string `json:"short_code"`
		Original   string `json:"original"`
		ClickCount uint   `json:"click_count"`
	}
	if err := json.NewDecoder(statsResp.Body).Decode(&statsBody); err != nil {
		t.Fatalf("Erro ao decodificar JSON de /stats: %v", err)
	}

	if statsBody.ShortCode != encurtarResp.ShortCode {
		t.Errorf("Esperado short_code=%s, obteve %s", encurtarResp.ShortCode, statsBody.ShortCode)
	}
	if statsBody.Original != "https://www.google.com" {
		t.Errorf("Esperado original=https://www.google.com, obteve %s", statsBody.Original)
	}
	if statsBody.ClickCount < 2 {
		t.Errorf("Esperado click_count >= 2, obteve %d", statsBody.ClickCount)
	}
}
