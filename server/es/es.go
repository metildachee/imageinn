package es

import (
	"context"
	"encoding/json"
	"github.com/metildachee/imageinn/server/config"
	"github.com/olivere/elastic/v7"
	"log"
)

// TODO: Add interface

type DocumentStructure struct {
	URL         string  `json:"url"`
	Caption     string  `json:"caption"`
	ID          string  `json:"id"`
	CategoryIDs []int64 `json:"category_ids"`
}

type Searcher struct {
	client *elastic.Client
	index  string
}

func NewSearcher(config config.Config) (*Searcher, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(config.Elasticsearch.Url),
		elastic.SetSniff(config.Elasticsearch.Sniff),
		elastic.SetHealthcheck(config.Elasticsearch.HealthCheck),
		//elastic.SetBasicAuth("elastic", "rTur2xp5QmFAvEAPJjfT"),
	)
	if err != nil {
		return nil, err
	}

	return &Searcher{
		client: client,
		index:  config.Elasticsearch.Index,
	}, nil
}

func unmarshalResults(hits []*elastic.SearchHit) ([]DocumentStructure, error) {
	documents := make([]DocumentStructure, len(hits))
	for i, hit := range hits {
		var doc DocumentStructure
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			return nil, err
		}
		doc.ID = hit.Id // Set the document ID from the search hit
		documents[i] = doc
	}
	return documents, nil
}

func (s *Searcher) doSearch(ctx context.Context, query elastic.Query) ([]DocumentStructure, int64, error) {
	src, err := query.Source()
	if err != nil {
		log.Println("error getting query source:", err)
		return nil, 0, err
	}
	srcBytes, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		log.Println("error marshaling query source to JSON:", err)
		return nil, 0, err
	}
	log.Println("query source:", string(srcBytes))

	searchResult, err := s.client.Search().
		Index(s.index).
		Query(query).
		From(0).Size(10).
		Do(ctx)
	if err != nil {
		log.Println("search got error", err)
		return nil, 0, err
	}

	totalHits := searchResult.TotalHits()
	hits, unmarshalErr := unmarshalResults(searchResult.Hits.Hits)
	log.Println("hits", hits)
	return hits, totalHits, unmarshalErr
}

func (s *Searcher) SearchByCategoryID(ctx context.Context, categoryID int64) ([]DocumentStructure, int64, error) {
	termQuery := elastic.NewTermQuery("category_ids", categoryID)
	return s.doSearch(ctx, termQuery)
}

func (s *Searcher) SearchByCategoryIDsAndOr(ctx context.Context, categoryIDs []int64) ([]DocumentStructure, int64, error) {
	andCategoryQuery := elastic.NewBoolQuery()
	for _, category := range categoryIDs {
		andCategoryQuery.Must(elastic.NewTermQuery("category_ids", category))
	}
	orCategoryQuery := elastic.NewBoolQuery()
	for _, category := range categoryIDs {
		orCategoryQuery.Should(elastic.NewTermQuery("category_ids", category))
	}
	query := elastic.NewBoolQuery().Should(andCategoryQuery, orCategoryQuery)
	return s.doSearch(ctx, query)
}

func (s *Searcher) SearchByKeywordsAndOr(ctx context.Context, keywords []string) ([]DocumentStructure, int64, error) {
	andQuery := elastic.NewBoolQuery()
	for _, keyword := range keywords {
		andQuery.Must(elastic.NewMatchQuery("caption", keyword))
	}
	orQuery := elastic.NewBoolQuery()
	for _, keyword := range keywords {
		orQuery.Should(elastic.NewMatchQuery("caption", keyword))
	}
	query := elastic.NewBoolQuery().Should(andQuery, orQuery)
	return s.doSearch(ctx, query)
}

func (s *Searcher) SearchByKeywordsAndCategoryIDsStrictAnd(ctx context.Context, categoryIDs []int64, keywords []string) ([]DocumentStructure, int64, error) {
	// All AND logic
	categoryQuery := elastic.NewBoolQuery()
	for _, category := range categoryIDs {
		categoryQuery.Must(elastic.NewTermQuery("category_ids", category))
	}
	keywordQuery := elastic.NewBoolQuery()
	for _, keyword := range keywords {
		keywordQuery.Must(elastic.NewMatchQuery("caption", keyword))
	}
	combinedQuery := elastic.NewBoolQuery().Must(categoryQuery, keywordQuery)

	return s.doSearch(ctx, combinedQuery)
}
