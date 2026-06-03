---
Title: Source: rag-eval-bm25-search.go.md
Ticket: RAGEVAL-GOJA-RAG-STRATEGIES
Status: active
Topics: [bleve]
DocType: reference
Intent: short-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Source material fetched for goja-bleve design investigation"
LastUpdated: 2026-06-02T21:30:00.000000-04:00
WhatFor: "Reference source material"
WhenToUse: "When investigating bleve documentation and APIs"
---

package search

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

type indexedChunk struct {
	ChunkID     string `json:"chunk_id"`
	DocumentID  string `json:"document_id"`
	SourceID    string `json:"source_id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	StrategyID  string `json:"strategy_id"`
	ChunkIndex  int    `json:"chunk_index"`
	Text        string `json:"text"`
	TokenCount  int    `json:"token_count"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
}

// BuildBM25 rebuilds a disposable Bleve BM25 index from canonical SQLite chunks.
func (s *Service) BuildBM25(ctx context.Context, req BuildIndexRequest) (*BuildIndexResult, error) {
	if req.StrategyID == "" {
		return nil, fmt.Errorf("strategy id is required")
	}
	sourceIDs := normalizeSourceIDs(req.SourceIDs)
	indexID := req.IndexID
	if indexID == "" {
		indexID = deriveIndexID(req.StrategyID, sourceIDs)
	}
	indexPath := s.indexPath(indexID)
	if _, err := os.Stat(indexPath); err == nil && !req.Force {
		return nil, fmt.Errorf("index %s already exists at %s; use --force to rebuild", indexID, indexPath)
	} else if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("stat index path: %w", err)
	}

	chunks, err := s.queries.ListChunksWithDocumentContextFiltered(req.StrategyID, sourceIDs, req.DocumentIDs, req.Limit)
	if err != nil {
		return nil, err
	}

	tmpPath := indexPath + ".tmp"
	if err := os.RemoveAll(tmpPath); err != nil {
		return nil, fmt.Errorf("remove temporary index path: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(indexPath), 0755); err != nil {
		return nil, fmt.Errorf("create index parent directory: %w", err)
	}

	mapping := bleve.NewIndexMapping()
	idx, err := bleve.New(tmpPath, mapping)
	if err != nil {
		return nil, fmt.Errorf("create bleve index: %w", err)
	}
	closed := false
	defer func() {
		if !closed {
			_ = idx.Close()
		}
	}()

	documentIDs := map[string]bool{}
	batch := idx.NewBatch()
	for i, ch := range chunks {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		documentIDs[ch.DocumentID] = true
		doc := indexedChunk{
			ChunkID:     ch.ChunkID,
			DocumentID:  ch.DocumentID,
			SourceID:    ch.SourceID,
			Title:       ch.Title,
			URL:         ch.URL,
			StrategyID:  ch.StrategyID,
			ChunkIndex:  ch.ChunkIndex,
			Text:        ch.Text,
			TokenCount:  ch.TokenCount,
			StartOffset: ch.StartOffset,
			EndOffset:   ch.EndOffset,
		}
		if err := batch.Index(ch.ChunkID, doc); err != nil {
			return nil, fmt.Errorf("index chunk %s: %w", ch.ChunkID, err)
		}
		if (i+1)%500 == 0 {
			if err := idx.Batch(batch); err != nil {
				return nil, fmt.Errorf("flush index batch: %w", err)
			}
			batch = idx.NewBatch()
		}
	}
	if batch.Size() > 0 {
		if err := idx.Batch(batch); err != nil {
			return nil, fmt.Errorf("flush final index batch: %w", err)
		}
	}
	if err := idx.Close(); err != nil {
		return nil, fmt.Errorf("close bleve index: %w", err)
	}
	closed = true

	if req.Force {
		if err := os.RemoveAll(indexPath); err != nil {
			return nil, fmt.Errorf("remove old index path: %w", err)
		}
	}
	if err := os.Rename(tmpPath, indexPath); err != nil {
		return nil, fmt.Errorf("move index into place: %w", err)
	}

	if err := s.queries.UpsertSearchIndex(db.SearchIndex{
		ID:            indexID,
		Name:          indexID,
		StrategyID:    req.StrategyID,
		IndexType:     "bm25",
		IndexPath:     indexPath,
		DocumentCount: len(documentIDs),
		ChunkCount:    len(chunks),
		Status:        "active",
	}); err != nil {
		return nil, err
	}

	return &BuildIndexResult{
		IndexID:       indexID,
		StrategyID:    req.StrategyID,
		SourceIDs:     sourceIDs,
		DocumentIDs:   req.DocumentIDs,
		IndexPath:     indexPath,
		ChunkCount:    len(chunks),
		DocumentCount: len(documentIDs),
		Rebuilt:       true,
	}, nil
}

// QueryBM25 runs lexical search against a previously built Bleve index.
func (s *Service) QueryBM25(ctx context.Context, req QueryRequest) (*QueryResult, error) {
	if req.IndexID == "" {
		return nil, fmt.Errorf("index id is required")
	}
	if req.Query == "" {
		return nil, fmt.Errorf("query is required")
	}
	if req.PreviewRunes == 0 {
		req.PreviewRunes = DefaultPreviewRunes
	}
	limit := normalizeLimit(req.Limit)

	meta, ok, err := s.queries.GetSearchIndex(req.IndexID)
	if err != nil {
		return nil, err
	}
	indexPath := s.indexPath(req.IndexID)
	if ok && meta.IndexPath != "" {
		indexPath = meta.IndexPath
	}

	idx, err := bleve.Open(indexPath)
	if err != nil {
		return nil, fmt.Errorf("open bm25 index %s: %w", req.IndexID, err)
	}
	defer idx.Close()

	textQuery := bleve.NewMatchQuery(req.Query)
	textQuery.SetField("text")
	titleQuery := bleve.NewMatchQuery(req.Query)
	titleQuery.SetField("title")
	titleQuery.SetBoost(2.0)
	query := bleve.NewDisjunctionQuery(textQuery, titleQuery)

	searchRequest := bleve.NewSearchRequestOptions(query, limit, 0, false)
	searchRequest.Fields = []string{"chunk_id", "document_id", "source_id", "title", "url", "strategy_id", "chunk_index", "text"}
	searchResult, err := idx.SearchInContext(ctx, searchRequest)
	if err != nil {
		return nil, fmt.Errorf("search bm25 index: %w", err)
	}

	result := &QueryResult{Query: req.Query, IndexID: req.IndexID, Retriever: "bm25"}
	for i, hit := range searchResult.Hits {
		fields := hit.Fields
		item := RetrievalResult{
			Rank:       i + 1,
			ChunkID:    stringField(fields, "chunk_id"),
			DocumentID: stringField(fields, "document_id"),
			SourceID:   stringField(fields, "source_id"),
			Title:      stringField(fields, "title"),
			URL:        stringField(fields, "url"),
			StrategyID: stringField(fields, "strategy_id"),
			ChunkIndex: intField(fields, "chunk_index"),
			Score:      hit.Score,
			Retriever:  "bm25",
			Preview:    preview(stringField(fields, "text"), req.PreviewRunes),
		}
		if item.ChunkID == "" {
			item.ChunkID = hit.ID
		}
		result.Items = append(result.Items, item)
	}
	return result, nil
}

func stringField(fields map[string]interface{}, name string) string {
	value, ok := fields[name]
	if !ok || value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}

func intField(fields map[string]interface{}, name string) int {
	value, ok := fields[name]
	if !ok || value == nil {
		return 0
	}
	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case jsonNumber:
		i, _ := v.Int64()
		return int(i)
	default:
		return 0
	}
}

type jsonNumber interface {
	Int64() (int64, error)
}
