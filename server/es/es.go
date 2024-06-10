package es

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/metildachee/imageinn/server/config"
	"io"
	"net/http"
)

type DocumentStructure struct {
	ImgBase64     string    `json:"img"`
	Title         string    `json:"title"`
	ID            string    `json:"id"`
	CategoryNames []string  `json:"category_names"`
	Embedding     []float64 `json:"embedding"`
	Score         float64   `json:"score"`
}

const (
	TitleField        = "title"
	CategoryNameField = "category_names"
)

type SearchClient struct {
	client   *elasticsearch.Client
	index    string
	endpoint string
}

func decodeResponseData(data []byte) ([]float64, error) {
	var respData interface{}
	err := json.Unmarshal(data, &respData)
	if err != nil {
		return nil, err
	}

	type ResponseData struct {
		TextFeatures []float64 `json:"text_features"`
	}
	var responseData ResponseData
	err = json.Unmarshal(data, &responseData)

	return responseData.TextFeatures, nil
}

func getEmbedding(txt string, endpoint string) ([]float64, error) {
	requestBody := map[string]interface{}{
		"text": txt,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("unexpected status code: %d", resp.StatusCode))
	}

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return decodeResponseData(bts)
}

func NewSearcher(config config.Config) (*SearchClient, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{config.Elasticsearch.Url},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	if config.Elasticsearch.ApiKey != "" {
		cfg.APIKey = config.Elasticsearch.ApiKey
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &SearchClient{
		client:   es,
		index:    config.Elasticsearch.Index,
		endpoint: config.ModelEndpoint,
	}, nil
}

func (s *SearchClient) doSearch(ctx context.Context, query map[string]interface{}) ([]DocumentStructure, int64, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, 0, err
	}

	res, err := s.client.Search(
		s.client.Search.WithContext(ctx),
		s.client.Search.WithIndex(s.index),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	var result map[string]interface{}
	if decodeErr := json.NewDecoder(res.Body).Decode(&result); decodeErr != nil {
		return nil, 0, decodeErr
	}

	hitsMap, ok := result["hits"].(map[string]interface{})
	if !ok {
		return nil, 0, errors.New("failed to extract hits map from response")
	}

	hits, ok := hitsMap["hits"].([]interface{})
	if !ok {
		return nil, 0, errors.New("failed to extract hits array from hits map")
	}

	docs := make([]DocumentStructure, 0)
	for _, hit := range hits {
		hitMap, ok := hit.(map[string]interface{})
		if !ok {
			return nil, 0, errors.New("error extracting hit from response")
		}

		source, ok := hitMap["_source"].(map[string]interface{})
		if !ok {
			return nil, 0, errors.New("error extracting _source from hit")
		}

		score, ok := hitMap["_score"].(float64)
		if !ok {
			return nil, 0, errors.New("error extracting _score from hit")
		}

		id, ok := hitMap["_id"].(string)
		if !ok {
			return nil, 0, errors.New("error extracting _score from hit")
		}

		var doc DocumentStructure
		docBytes, marshalErr := json.Marshal(source)
		if marshalErr != nil {
			fmt.Println("Error marshaling _source to bytes:", marshalErr)
			return nil, 0, marshalErr
		}
		if unmarshalErr := json.Unmarshal(docBytes, &doc); unmarshalErr != nil {
			fmt.Println("Error unmarshaling _source into DocumentStructure:", unmarshalErr)
			return nil, 0, unmarshalErr
		}

		doc.Score = score
		doc.ID = id
		docs = append(docs, doc)
	}

	totalDocs, ok := result["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	if !ok {
		return nil, 0, errors.New("could not get total number of hits")
	}
	return docs, int64(totalDocs), nil
}

func (s *SearchClient) SearchText(ctx context.Context, query string) ([]DocumentStructure, int64, error) {
	esQuery := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"img", "embedding"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"title": query,
						},
					},
					{
						"match": map[string]interface{}{
							"category_names": query,
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
	}

	return s.doSearch(ctx, esQuery)
}

func (s *SearchClient) SearchTextFuzzy(ctx context.Context, query string) ([]DocumentStructure, int64, error) {
	esQuery := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"img", "embedding"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":     query,
							"fields":    []string{"title", "category_names"},
							"fuzziness": "AUTO",
						},
					},
				},
			},
		},
	}
	return s.doSearch(ctx, esQuery)
}

func (s *SearchClient) SearchTextWithExclusions(ctx context.Context, query string, excludes []string) ([]DocumentStructure, int64, error) {
	esQuery := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"img", "embedding"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"title": query,
						},
					},
					{
						"match": map[string]interface{}{
							"category_names": query,
						},
					},
				},
				"must_not": []map[string]interface{}{
					{
						"terms": map[string]interface{}{
							"title": excludes,
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
	}
	return s.doSearch(ctx, esQuery)
}

func (s *SearchClient) SearchTextInImage(ctx context.Context, query string) ([]DocumentStructure, int64, error) {
	embedding, err := getEmbedding(query, s.endpoint)
	if err != nil {
		return nil, 0, err
	}

	return s.doKNN(ctx, embedding)
}

func (s *SearchClient) SearchTextWithFuzzy(ctx context.Context, query string, isAnd bool, excludes []string) ([]DocumentStructure, int64, error) {
	logic := "should"
	if isAnd {
		logic = "must"
	}
	esQuery := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"embedding"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				logic: []map[string]interface{}{
					{
						"multi_match": map[string]interface{}{
							"query":     query,
							"fields":    []string{"title", "category_names"},
							"fuzziness": "AUTO",
						},
					},
					{
						"match": map[string]interface{}{
							"title": query,
						},
					},
					{
						"match": map[string]interface{}{
							"category_names": query,
						},
					},
				},
				"must_not": []map[string]interface{}{
					{
						"terms": map[string]interface{}{
							"title": excludes,
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
	}
	return s.doSearch(ctx, esQuery)
}

func (s *SearchClient) SearchTextNoFuzzy(ctx context.Context, query string, isAnd bool, excludes []string) ([]DocumentStructure, int64, error) {
	logic := "should"
	if isAnd {
		logic = "must"
	}
	esQuery := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"embedding"},
		},
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				logic: []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"title": query,
						},
					},
					{
						"match": map[string]interface{}{
							"category_names": query,
						},
					},
				},
				"must_not": []map[string]interface{}{
					{
						"terms": map[string]interface{}{
							"title": excludes,
						},
					},
					{
						"terms": map[string]interface{}{
							"category_names": excludes,
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
	}
	return s.doSearch(ctx, esQuery)
}

func (s *SearchClient) searchByID(ctx context.Context, id string) ([]DocumentStructure, int64, error) {
	esQuery := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"img"},
		},
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": []string{id},
			},
		},
	}

	return s.doSearch(ctx, esQuery)
}

func (s *SearchClient) SearchByIDs(ctx context.Context, ids []string) ([]DocumentStructure, int64, error) {
	esQuery := map[string]interface{}{
		"_source": map[string]interface{}{
			"excludes": []string{"embedding"},
		},
		"query": map[string]interface{}{
			"ids": map[string]interface{}{
				"values": ids,
			},
		},
	}

	return s.doSearch(ctx, esQuery)
}

func (s *SearchClient) SearchSimilarByID(ctx context.Context, id string) ([]DocumentStructure, int64, error) {
	docs, _, err := s.searchByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}

	return s.doKNN(ctx, docs[0].Embedding)
}

func (s *SearchClient) doKNN(ctx context.Context, target []float64) ([]DocumentStructure, int64, error) {
	knnQuery := map[string]interface{}{
		"knn": map[string]interface{}{
			"field":          "embedding",
			"query_vector":   target,
			"k":              10,
			"num_candidates": 100,
		},
		"fields": []string{TitleField, CategoryNameField},
	}

	return s.doSearch(ctx, knnQuery)
}
