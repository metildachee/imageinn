package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	config "github.com/metildachee/imageinn/server/config"
	"github.com/metildachee/imageinn/server/utils"
	"github.com/olivere/elastic/v7" // imports as package "elastic"
	"log"
	"net/http"
	"strings"
)

type DocumentStructure struct {
	URL         string  `json:"url"`
	Caption     string  `json:"caption"`
	ID          string  `json:"id"`
	CategoryIDs []int64 `json:"category_ids"`
}

type SearchRequest struct {
	from        int64
	to          int64
	keywords    []string
	categoryIDs []int64
	imageID     int64
	requestID   string
}

type SearchResponse struct {
	Images     []DocumentStructure `json:"images"`
	TotalCount int64               `json:"total_count"`
}

type Searcher struct {
	client *elastic.Client
	index  string
}

func unmarshalResults(hits []*elastic.SearchHit) ([]DocumentStructure, error) {
	documents := make([]DocumentStructure, 0)

	for _, hit := range hits {
		doc := DocumentStructure{}
		if unmarshalErr := json.Unmarshal(hit.Source, &doc); unmarshalErr != nil {
			return nil, unmarshalErr
		}
		documents = append(documents, doc)
		log.Printf("document by URL: %s, Caption: %s\n", doc.URL, doc.Caption)
	}

	return documents, nil
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

type WebHandler struct {
	searcher *Searcher
}

func NewWebHandler(searcher *Searcher) *WebHandler {
	return &WebHandler{searcher: searcher}
}

func (s *Searcher) doSearch(ctx context.Context, query elastic.Query) ([]DocumentStructure, int64, error) {
	searchResult, err := s.client.Search().
		Index(s.index).
		Query(query).
		From(0).Size(10).
		Do(ctx)
	if err != nil {
		return nil, 0, err
	}
	src, _ := query.Source()
	log.Println("source", src)

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

func validateAndProcessRequest(ctx context.Context, r *http.Request) (*SearchRequest, error) {
	queryParameters := r.URL.Query()
	keywordInputs := queryParameters.Get("q")
	categoryInputs := queryParameters.Get("category_ids")
	imageIDInput := queryParameters.Get("id")

	keywordInputsTrimmed := strings.Trim(keywordInputs, "")
	categoryInputsTrimmed := strings.Trim(categoryInputs, "")
	imageIDInputsTrimmed := strings.Trim(imageIDInput, "")

	if keywordInputsTrimmed == "" && categoryInputsTrimmed == "" && imageIDInputsTrimmed == "" {
		return nil, fmt.Errorf("empty request")
	}

	searchRequest := &SearchRequest{}

	var keywords []string
	if keywordInputsTrimmed != "" {
		keywords = strings.Split(keywordInputs, " ")
		searchRequest.keywords = keywords
	}

	if categoryInputsTrimmed != "" {
		categories, err := utils.StringToInt64Array(categoryInputs)
		if err != nil {
			return nil, err
		}
		searchRequest.categoryIDs = categories
	}

	if imageIDInputsTrimmed != "" {
		imageID, err := utils.StrToInt64(imageIDInput)
		if err != nil {
			return nil, err
		}
		searchRequest.imageID = imageID
	}

	if searchRequest.imageID != 0 && (len(searchRequest.categoryIDs) > 0 || len(searchRequest.keywords) > 0) {
		return nil, fmt.Errorf("image ID should be independent results but got keywords || categories")
	}

	return searchRequest, nil
}

func (h *WebHandler) searchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	searchReq, err := validateAndProcessRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var (
		searchResults []DocumentStructure
		count         int64
	)
	if len(searchReq.keywords) != 0 && len(searchReq.categoryIDs) == 0 && searchReq.imageID == 0 { // keyword search
		searchResults, count, err = h.searcher.SearchByKeywordsAndOr(ctx, searchReq.keywords)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	if len(searchReq.keywords) == 0 && len(searchReq.categoryIDs) != 0 && searchReq.imageID == 0 { // category search
		searchResults, count, err = h.searcher.SearchByCategoryIDsAndOr(ctx, searchReq.categoryIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// document search

	if len(searchReq.keywords) != 0 && len(searchReq.categoryIDs) != 0 && searchReq.imageID == 0 { // keyword and category search
		searchResults, count, err = h.searcher.SearchByKeywordsAndCategoryIDsStrictAnd(ctx, searchReq.categoryIDs, searchReq.keywords)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	searchResp := &SearchResponse{
		Images:     searchResults,
		TotalCount: count,
	}

	jsonResponse, err := json.Marshal(searchResp)
	if err != nil {
		http.Error(w, "failed to serialize search results", http.StatusInternalServerError)
		return
	}

	log.Println("request", searchReq, "response", searchResp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	configPath := "config/config.yml" // change into flag
	serverConfig := config.LoadConfig(configPath)
	searcher, err := NewSearcher(*serverConfig)
	if err != nil {
		log.Fatalln("load searcher failed", err)
	}
	webHandler := NewWebHandler(searcher)

	r := mux.NewRouter()
	r.HandleFunc("/search", webHandler.searchHandler).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
