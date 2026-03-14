package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielcaso/mini-blockchain/internal/blockchain"
)

func setupTestRouter() *httptest.Server {
	bc, err := blockchain.NewBlockchain(2, "")
	if err != nil {
		panic(err)
	}
	r := SetupRouter(bc)
	return httptest.NewServer(r)
}

func TestHealthCheck(t *testing.T) {
	ts := setupTestRouter()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var body HealthResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", body.Status)
	}
}

func TestGetChain(t *testing.T) {
	ts := setupTestRouter()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/chain")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body ChainResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if body.Length != 1 {
		t.Errorf("expected length 1 (genesis), got %d", body.Length)
	}
}

func TestMineBlock(t *testing.T) {
	ts := setupTestRouter()
	defer ts.Close()

	payload := `{"data":"test block"}`
	resp, err := http.Post(ts.URL+"/mine", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}

	var block blockchain.Block
	json.NewDecoder(resp.Body).Decode(&block)
	if block.Data != "test block" {
		t.Errorf("expected data 'test block', got '%s'", block.Data)
	}
	if block.Index != 1 {
		t.Errorf("expected index 1, got %d", block.Index)
	}
}

func TestMineBlockBadRequest(t *testing.T) {
	ts := setupTestRouter()
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/mine", "application/json", bytes.NewBufferString(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestValidateChain(t *testing.T) {
	ts := setupTestRouter()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/validate")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body ValidateResponse
	json.NewDecoder(resp.Body).Decode(&body)
	if !body.Valid {
		t.Error("chain should be valid")
	}
}

func TestGetBlockByID(t *testing.T) {
	ts := setupTestRouter()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/block/0")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var block blockchain.Block
	json.NewDecoder(resp.Body).Decode(&block)
	if block.Index != 0 {
		t.Errorf("expected genesis (index 0), got %d", block.Index)
	}
}

func TestGetBlockNotFound(t *testing.T) {
	ts := setupTestRouter()
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/block/999")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}
