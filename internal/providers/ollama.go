package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OllamaProvider implements the Provider interface for local Ollama instances.
type OllamaProvider struct {
	baseURL      string
	defaultModel string
	httpClient   *http.Client
}

// NewOllamaProvider creates a new Ollama provider.
func NewOllamaProvider(baseURL string) *OllamaProvider {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	return &OllamaProvider{
		baseURL:      baseURL,
		defaultModel: "qwen3-vl:235b-cloud", // Modelo especificado pelo usuário
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Aumentado para modelos maiores
		},
	}
}

func (p *OllamaProvider) ID() string {
	return "ollama"
}

// SetDefaultModel allows changing the model used for Generate.
func (p *OllamaProvider) SetDefaultModel(model string) {
	p.defaultModel = model
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Options  map[string]any  `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaChatResponse struct {
	Message      ollamaMessage `json:"message"`
	PromptEval   int           `json:"prompt_eval_count"`
	EvalCount    int           `json:"eval_count"`
	TotalDuration int64        `json:"total_duration"`
}

func (p *OllamaProvider) Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error) {
	url := fmt.Sprintf("%s/api/chat", p.baseURL)

	model := req.Model
	if model == "" {
		model = p.defaultModel
	}

	messages := make([]ollamaMessage, 0, len(req.Messages))
	for _, m := range req.Messages {
		messages = append(messages, ollamaMessage{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}

	ollamaReq := ollamaChatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Options: map[string]any{
			"temperature": req.Temperature,
		},
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return GenerateResponse{}, err
	}

	hreq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return GenerateResponse{}, err
	}
	hreq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(hreq)
	if err != nil {
		return GenerateResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return GenerateResponse{}, fmt.Errorf("ollama error (%d): %s", resp.StatusCode, string(respBody))
	}

	var ollamaResp ollamaChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return GenerateResponse{}, err
	}

	return GenerateResponse{
		Content:      ollamaResp.Message.Content,
		InputTokens:  ollamaResp.PromptEval,
		OutputTokens: ollamaResp.EvalCount,
	}, nil
}

type ollamaEmbedRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
}

type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (p *OllamaProvider) Embed(ctx context.Context, text string) ([]float32, error) {
	url := fmt.Sprintf("%s/api/embeddings", p.baseURL)

	ollamaReq := ollamaEmbedRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, err
	}

	hreq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	hreq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(hreq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama embed error (%d): %s", resp.StatusCode, string(respBody))
	}

	var ollamaResp ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	return ollamaResp.Embedding, nil
}
