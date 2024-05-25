package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/metildachee/imageinn/server/es"
	"github.com/metildachee/imageinn/server/memcache"
	"github.com/metildachee/imageinn/server/utils"
	"log"
	"net/http"
	"strings"
)

type SearchRequest struct {
	from        int64
	to          int64
	keywords    []string
	categoryIDs []int64
	imageID     int64
	requestID   string
}

type SearchResponse struct {
	Images     []es.DocumentStructure `json:"images"`
	TotalCount int64                  `json:"total_count"`
}

type Bucket struct {
	Key   string `json:"key_string"`
	Count int64  `json:"count"`
}

type CategoryRequest struct {
	Count int64 `json:"count"`
}

type CategoryResponse struct {
	Categories []Bucket `json:"categories"`
}

type WebHandler struct {
	searcher *es.Searcher
	memcache *memcache.Memcache
}

func NewWebHandler(searcher *es.Searcher, memcache *memcache.Memcache) *WebHandler {
	return &WebHandler{searcher: searcher, memcache: memcache}
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

func validateCategoryRequest(ctx context.Context, r *http.Request) (*CategoryRequest, error) {
	queryParameters := r.URL.Query()
	countQuery := queryParameters.Get("count")

	count, err := utils.StrToInt64(countQuery)
	if err != nil {
		return nil, err
	}

	if count <= 0 {
		count = 10
	}

	return &CategoryRequest{Count: count}, nil
}

func (h *WebHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	searchReq, err := validateAndProcessRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var (
		searchResults []es.DocumentStructure
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
		log.Println("SearchByKeywordsAndCategoryIDsStrictAnd", searchResults)
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

func (h *WebHandler) CategoryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	catReq, err := validateCategoryRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	catResults, err := h.searcher.SearchCategoryInformation(ctx, int(catReq.Count))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	categoryResponse := CategoryResponse{Categories: make([]Bucket, 0)}
	for _, category := range catResults {
		id, conversionErr := utils.StrToInt64(utils.InterfaceToString(category.Key))
		if conversionErr != nil {
			http.Error(w, conversionErr.Error(), http.StatusInternalServerError)
		}

		catName, lookupErr := h.memcache.Lookup(id)
		if lookupErr != nil {
			log.Println("look up got error but its ok", lookupErr)
		}
		categoryResponse.Categories = append(categoryResponse.Categories, Bucket{
			Key:   catName,
			Count: category.DocCount,
		})
	}

	jsonResponse, err := json.Marshal(categoryResponse)
	if err != nil {
		http.Error(w, "failed to serialize search results", http.StatusInternalServerError)
		return
	}

	log.Println("request", catReq, "response", catResults)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
