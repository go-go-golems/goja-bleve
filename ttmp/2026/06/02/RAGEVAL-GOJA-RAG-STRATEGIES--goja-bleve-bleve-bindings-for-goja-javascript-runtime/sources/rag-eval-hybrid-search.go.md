---
Title: Source: rag-eval-hybrid-search.go.md
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
	"sort"
)

const DefaultRRFK = 60

type HybridQueryRequest struct {
	BM25        QueryRequest
	Vector      VectorQueryRequest
	Limit       int
	RRFK        int
	BM25Limit   int
	VectorLimit int
}

// QueryHybrid runs BM25 and vector retrieval independently, then merges their
// ranked lists with reciprocal-rank fusion.
func (s *Service) QueryHybrid(ctx context.Context, req HybridQueryRequest) (*QueryResult, error) {
	if req.BM25.Query == "" && req.Vector.Query == "" {
		return nil, fmt.Errorf("query is required")
	}
	query := req.BM25.Query
	if query == "" {
		query = req.Vector.Query
	}
	if req.BM25.Query == "" {
		req.BM25.Query = query
	}
	if req.Vector.Query == "" {
		req.Vector.Query = query
	}
	if req.RRFK <= 0 {
		req.RRFK = DefaultRRFK
	}
	limit := normalizeLimit(req.Limit)
	if req.BM25Limit > 0 {
		req.BM25.Limit = req.BM25Limit
	} else if req.BM25.Limit <= 0 {
		req.BM25.Limit = 50
	}
	if req.VectorLimit > 0 {
		req.Vector.Limit = req.VectorLimit
	} else if req.Vector.Limit <= 0 {
		req.Vector.Limit = 50
	}

	bm25Result, err := s.QueryBM25(ctx, req.BM25)
	if err != nil {
		return nil, fmt.Errorf("bm25 search: %w", err)
	}
	vectorResult, err := s.QueryVector(ctx, req.Vector)
	if err != nil {
		return nil, fmt.Errorf("vector search: %w", err)
	}

	merged := map[string]*RetrievalResult{}
	add := func(name string, items []RetrievalResult) {
		for _, item := range items {
			candidate, ok := merged[item.ChunkID]
			if !ok {
				copy := item
				copy.Score = 0
				copy.Retriever = "hybrid"
				copy.Components = map[string]RetrievalComponent{}
				merged[item.ChunkID] = &copy
				candidate = &copy
			}
			candidate.Score += 1.0 / float64(req.RRFK+item.Rank)
			candidate.Components[name] = RetrievalComponent{Rank: item.Rank, Score: item.Score}
		}
	}
	add("bm25", bm25Result.Items)
	add("vector", vectorResult.Items)

	items := make([]RetrievalResult, 0, len(merged))
	for _, item := range merged {
		items = append(items, *item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score == items[j].Score {
			return items[i].ChunkID < items[j].ChunkID
		}
		return items[i].Score > items[j].Score
	})
	if len(items) > limit {
		items = items[:limit]
	}
	for i := range items {
		items[i].Rank = i + 1
	}
	return &QueryResult{Query: query, IndexID: req.BM25.IndexID, Retriever: "hybrid", Items: items}, nil
}
