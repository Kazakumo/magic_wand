package embedding

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultTimeoutSec = 30
	maxRetries        = 3
)

type OllamaEmbedder struct {
	baseURL    string
	model      string
	timeoutSec int
	httpClient *http.Client
}

func NewOllamaEmbedder(baseURL, model string, timeoutSec int) *OllamaEmbedder {
	if timeoutSec <= 0 {
		timeoutSec = defaultTimeoutSec
	}
	return &OllamaEmbedder{
		baseURL:    baseURL,
		model:      model,
		timeoutSec: timeoutSec,
		httpClient: &http.Client{},
	}
}

type embedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (o *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	var lastErr error
	for attempt := range maxRetries {
		if attempt > 0 {
			wait := time.Duration(1<<(attempt-1)) * time.Second
			select {
			case <-ctx.Done():
				return nil, fmt.Errorf("context canceled during retry wait: %w", ctx.Err())
			case <-time.After(wait):
			}
		}

		vec, err := o.doEmbed(ctx, text)
		if err == nil {
			return vec, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("embed after %d retries: %w", maxRetries, lastErr)
}

func (o *OllamaEmbedder) doEmbed(ctx context.Context, text string) ([]float32, error) {
	reqCtx, cancel := context.WithTimeout(ctx, time.Duration(o.timeoutSec)*time.Second)
	defer cancel()

	body, err := json.Marshal(embedRequest{Model: o.model, Prompt: text})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, o.baseURL+"/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if len(result.Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding returned")
	}

	return result.Embedding, nil
}
