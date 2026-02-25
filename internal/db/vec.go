package db

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"

	"github.com/msrovani/myclaw/internal/core"
)

// SearchResult representa um resultado de busca (Vetor ou FTS).
type SearchResult struct {
	ID       string
	Content  string
	Metadata string
	Distance float32 // Para busca vetorial (menor é melhor)
	Score    float32 // Para busca FTS
}

// Float32ToBytes converte float32 para slice de bytes (Little Endian).
func Float32ToBytes(floats []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, floats)
	return buf.Bytes(), err
}

// BytesToFloat32 converte slice de bytes para float32.
func BytesToFloat32(b []byte) ([]float32, error) {
	if len(b)%4 != 0 {
		return nil, fmt.Errorf("invalid byte length for float32 array: %d", len(b))
	}
	floats := make([]float32, len(b)/4)
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &floats)
	return floats, err
}

// CosineDistance calcula a distância de cosseno.
func CosineDistance(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 1.0
	}

	var dot, normA, normB float32
	n := len(a)
	for i := 0; i < n; i++ {
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

// SearchVector usa a extensão sqlite-vec ou fallback em Go.
func (d *DB) SearchVector(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("search vector: %w", err)
	}

	embBytes, err := Float32ToBytes(queryEmbedding)
	if err != nil {
		return nil, err
	}

	q := `SELECT m.id, m.content, m.metadata, v.distance
	      FROM memories_vec0 v
	      JOIN memories m ON v.rowid = m.rowid
	      WHERE v.embedding MATCH ?1 AND m.uid = ?2 AND m.workspace_id = ?3
	      ORDER BY v.distance LIMIT ?4`

	rows, err := d.ReadRows(ctx, q, embBytes, tc.UID, tc.WorkspaceID, limit)
	if err != nil {
		return d.SearchVectorFallback(ctx, queryEmbedding, limit)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var res SearchResult
		if err := rows.Scan(&res.ID, &res.Content, &res.Metadata, &res.Distance); err != nil {
			return nil, err
		}
		results = append(results, res)
	}
	return results, nil
}

// SearchVectorFallback busca KNN em memória Go.
func (d *DB) SearchVectorFallback(ctx context.Context, queryEmbedding []float32, limit int) ([]SearchResult, error) {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("search vector fallback: %w", err)
	}

	q := "SELECT id, content, metadata, embedding FROM memories WHERE uid = ? AND workspace_id = ? AND embedding IS NOT NULL"
	rows, err := d.ReadRows(ctx, q, tc.UID, tc.WorkspaceID)
	if err != nil {
		return nil, fmt.Errorf("search vector fallback query: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var res SearchResult
		var embBytes []byte
		if err := rows.Scan(&res.ID, &res.Content, &res.Metadata, &embBytes); err != nil {
			continue
		}

		dbEmb, err := BytesToFloat32(embBytes)
		if err != nil {
			continue
		}

		res.Distance = CosineDistance(queryEmbedding, dbEmb)
		results = append(results, res)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Distance < results[j].Distance
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}

// SearchFTS realiza busca textual usando FTS5 com isolamento por subquery e parâmetros numerados.
func (d *DB) SearchFTS(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	tc, err := core.TenantFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("search fts: %w", err)
	}

	query = sanitizeFTSQuery(query)
	if query == "" {
		return nil, nil
	}

	// Estratégia Final: Isolar o MATCH em uma subquery CTE e usar parâmetros numerados (?1).
	// Isso evita conflitos de parsing entre o operador MATCH e o JOIN no driver modernc.org/sqlite.
	q := `WITH matched_fts AS (
	          SELECT rowid, rank
	          FROM memories_fts
	          WHERE memories_fts MATCH ?1
	      )
	      SELECT m.id, m.content, m.metadata, f.rank
	      FROM matched_fts f
	      JOIN memories m ON f.rowid = m.rowid
	      WHERE m.uid = ?2 AND m.workspace_id = ?3
	      ORDER BY f.rank LIMIT ?4`

	rows, err := d.ReadRows(ctx, q, query, tc.UID, tc.WorkspaceID, limit)
	if err != nil {
		return nil, fmt.Errorf("search fts query: %w", err)
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var res SearchResult
		var rank float64
		if err := rows.Scan(&res.ID, &res.Content, &res.Metadata, &rank); err != nil {
			return nil, err
		}
		res.Score = float32(rank)
		results = append(results, res)
	}

	return results, nil
}

func sanitizeFTSQuery(query string) string {
	var b strings.Builder
	b.Grow(len(query))
	needSpace := false

	for _, r := range query {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if needSpace && b.Len() > 0 {
				b.WriteRune(' ')
			}
			b.WriteRune(r)
			needSpace = false
			continue
		}
		if unicode.IsSpace(r) {
			needSpace = true
			continue
		}
		needSpace = true
	}

	return strings.TrimSpace(b.String())
}

// ReciprocalRankFusion combina os resultados.
func ReciprocalRankFusion(vectorResults []SearchResult, ftsResults []SearchResult, k int) []SearchResult {
	if k <= 0 {
		k = 60
	}
	scoreMap := make(map[string]float32)
	contentMap := make(map[string]SearchResult)

	for i, res := range vectorResults {
		rank := i + 1
		scoreMap[res.ID] += 1.0 / float32(k+rank)
		contentMap[res.ID] = res
	}

	for i, res := range ftsResults {
		rank := i + 1
		scoreMap[res.ID] += 1.0 / float32(k+rank)
		contentMap[res.ID] = res
	}

	var combined []SearchResult
	for id, score := range scoreMap {
		res := contentMap[id]
		res.Score = score
		combined = append(combined, res)
	}

	sort.Slice(combined, func(i, j int) bool {
		return combined[i].Score > combined[j].Score
	})

	return combined
}
