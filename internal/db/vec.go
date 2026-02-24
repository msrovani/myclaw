package db

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"sort"

	"github.com/msrovani/myclaw/internal/core"
)

// SearchResult represents a matched memory from vector or FTS search.
type SearchResult struct {
	ID       string
	Content  string
	Metadata string
	Distance float32 // For vector search (lower is closer)
	Score    float32 // For FTS search
}

// Float32ToBytes converts a slice of float32 to a byte slice for BLOB storage.
// This matches sqlite-vec's float32 format (Little Endian).
func Float32ToBytes(floats []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, floats)
	return buf.Bytes(), err
}

// BytesToFloat32 converts a BLOB byte slice back into a float32 slice.
func BytesToFloat32(b []byte) ([]float32, error) {
	if len(b)%4 != 0 {
		return nil, fmt.Errorf("invalid byte length for float32 array: %d", len(b))
	}
	floats := make([]float32, len(b)/4)
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &floats)
	return floats, err
}

// CosineDistance calculates 1 - CosineSimilarity. Lower means closer.
func CosineDistance(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 1.0 // Max distance if lengths mismatch or are empty
	}
	var dot, normA, normB float32
	// Optional loop unrolling for performance could be done here.
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 1.0
	}
	similarity := dot / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
	return 1.0 - similarity
}

// SearchVectorFallback performs an exact KNN search in Go memory.
// It retrieves ALL embeddings for the tenant from the physical DB,
// calculates Cosine Distance against the query, and returns the top K.
// This is used for graceful degradation when sqlite-vec (CGo) is unavailable.
func (d *DB) SearchVectorFallback(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("search vector: %w", err)
	}

	// Read all embeddings. Defense in depth: filter by uid/workspace_id even in isolated DB.
	q := "SELECT id, content, metadata, embedding FROM memories WHERE uid = ? AND workspace_id = ? AND embedding IS NOT NULL"
	rows, err := d.ReadRows(ctx, q, tc.UID, tc.WorkspaceID)
	if err != nil {
		return nil, fmt.Errorf("search vector query: %w", err)
	}
	defer rows.Close()

	var results []SearchResult

	for rows.Next() {
		var res SearchResult
		var embBytes []byte
		if err := rows.Scan(&res.ID, &res.Content, &res.Metadata, &embBytes); err != nil {
			return nil, fmt.Errorf("search vector scan: %w", err)
		}

		dbEmb, err := BytesToFloat32(embBytes)
		if err != nil {
			// Skip malformed blobs naturally
			continue
		}

		res.Distance = CosineDistance(queryEmbedding, dbEmb)
		results = append(results, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("search vector iteration: %w", err)
	}

	// Sort by distance (lower is better)
	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}
