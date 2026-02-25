package memory

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/msrovani/myclaw/internal/core"
	"github.com/msrovani/myclaw/internal/db"
	"github.com/msrovani/myclaw/internal/providers"
)

// TestEngine_OllamaIntegration tests the memory engine with a real local Ollama instance.
// To run this, ensure Ollama is running and has 'nomic-embed-text' and a chat model (e.g. 'llama3' or 'mistral') installed.
// Run with: go test -v ./internal/memory -run TestEngine_OllamaIntegration
func TestEngine_OllamaIntegration(t *testing.T) {
	ollamaURL := os.Getenv("XXXCLAW_OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	// 1. Setup Provider
	provider := providers.NewOllamaProvider(ollamaURL)

	// Quick check if Ollama is alive
	ctx := context.Background()
	_, err := provider.Embed(ctx, "health check")
	if err != nil {
		t.Skipf("Ollama not available or model 'nomic-embed-text' missing: %v", err)
	}

	// 2. Setup DB and Engine
	mgr := db.NewManager(db.Config{BaseDataDir: t.TempDir()})
	defer mgr.CloseAll()

	engine := NewEngine(mgr, provider)

	tc := core.TenantContext{UID: "test_user", WorkspaceID: "test_ws"}
	tctx := core.WithTenant(ctx, tc)

	// 3. Add memories with automatic embeddings
	fmt.Println("Adding memories...")
	_, err = engine.AddMemory(tctx, "O usuário gosta de café preto e tecnologia Go.", AddOptions{
		Layer: LayerShortTerm,
		SessionID: "sess_1",
	})
	if err != nil {
		t.Fatalf("AddMemory 1: %v", err)
	}

	_, err = engine.AddMemory(tctx, "O usuário mora em Florianópolis.", AddOptions{
		Layer: LayerShortTerm,
		SessionID: "sess_1",
	})
	if err != nil {
		t.Fatalf("AddMemory 2: %v", err)
	}

	// 4. Search (Hybrid Search Test)
	fmt.Println("Testing Hybrid Search...")
	results, err := engine.SearchMemories(tctx, "O que o usuário gosta de beber?", SearchOptions{Limit: 5})
	if err != nil {
		t.Fatalf("SearchMemories: %v", err)
	}

	if len(results) == 0 {
		t.Errorf("No results found in hybrid search")
	} else {
		fmt.Printf("Top Match: %s (Score: %f)\n", results[0].Content, results[0].Score)
	}

	// 5. Compaction (Summarization Test)
	fmt.Println("Testing Recursive Summarization (Compaction)...")
	// Note: Ollama needs a chat model. We assume 'llama3' or similar is default or available.
	// You might need to adjust the Generate request in engine.go if you want to specify a model.
	// For now, it uses what's in the request. Let's ensure a model is specified.

	reduced, err := engine.CompactMemories(tctx, "sess_1")
	if err != nil {
		t.Fatalf("CompactMemories: %v. Make sure you have a default chat model in Ollama.", err)
	}

	fmt.Printf("Compacted %d short-term memories into 1 long-term summary.\n", reduced)

	// 6. Verify Search after compaction
	fmt.Println("Testing Search after compaction...")
	results2, err := engine.SearchMemories(tctx, "Onde o usuário mora?", SearchOptions{Limit: 5})
	if err != nil {
		t.Fatalf("Search after compaction failed: %v", err)
	}

	for _, r := range results2 {
		fmt.Printf("- Found: %s\n", r.Content)
	}
}
