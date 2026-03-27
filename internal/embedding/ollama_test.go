package embedding_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kazakumo/magic_wand/internal/embedding"
)

func newTestEmbedder(baseURL string) *embedding.OllamaEmbedder {
	return embedding.NewOllamaEmbedder(baseURL, "nomic-embed-text", 5)
}

func TestEmbed_Success(t *testing.T) {
	want := []float32{0.1, 0.2, 0.3}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/embeddings" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{"embedding": want})
	}))
	defer srv.Close()

	e := newTestEmbedder(srv.URL)
	got, err := e.Embed(context.Background(), "hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("want %d dims, got %d", len(want), len(got))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("dim[%d]: want %v, got %v", i, want[i], got[i])
		}
	}
}

func TestEmbed_ServerError_Retries(t *testing.T) {
	attempts := 0
	want := []float32{0.5, 0.6}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"embedding": want})
	}))
	defer srv.Close()

	e := newTestEmbedder(srv.URL)
	got, err := e.Embed(context.Background(), "retry test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 3 {
		t.Errorf("want 3 attempts, got %d", attempts)
	}
	if len(got) != len(want) {
		t.Fatalf("want %d dims, got %d", len(want), len(got))
	}
}

func TestEmbed_AllRetriesFail(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	e := newTestEmbedder(srv.URL)
	_, err := e.Embed(context.Background(), "fail")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestEmbed_ContextCanceled(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	e := newTestEmbedder(srv.URL)
	_, err := e.Embed(ctx, "canceled")
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestEmbed_EmptyEmbedding(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"embedding": []float32{}})
	}))
	defer srv.Close()

	e := newTestEmbedder(srv.URL)
	_, err := e.Embed(context.Background(), "empty")
	if err == nil {
		t.Fatal("expected error for empty embedding")
	}
}
