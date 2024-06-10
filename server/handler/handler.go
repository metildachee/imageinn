package handler

import (
	"context"
	"encoding/json"
	"github.com/metildachee/imageinn/server/es"
	"github.com/metildachee/imageinn/server/memcache"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type SearchRequestLogic struct {
	And bool `json:"and"`
	Or  bool `json:"or"`
}

type SearchRequestTextOptions struct {
	IsNLP    bool     `json:"is_nlp"`
	IsFuzzy  bool     `json:"is_fuzzy"`
	Excludes []string `json:"excludes"`
	IsAnd    bool     `json:"is_and"`
}

type SearchRequestText struct {
	Query       string                   `json:"query"`
	TextOptions SearchRequestTextOptions `json:"text-options"`
}

type SearchRequestImage struct {
	Query string `json:"query"`
}

type SearchRequestID struct {
	ID string `json:"id"`
}

type SearchRequestIDs struct {
	IDs []string `json:"ids"`
}

func (r *SearchRequestID) GetID() string {
	if r == nil {
		return ""
	}
	return r.ID
}

func (r *SearchRequestIDs) GetIDs() []string {
	if r == nil {
		return []string{}
	}
	return r.IDs
}

func (r *SearchRequestImage) GetQuery() string {
	if r == nil {
		return ""
	}
	return r.Query
}

func (rt *SearchRequestText) GetQuery() string {
	if rt == nil {
		return ""
	}
	return rt.Query
}

func (rt *SearchRequestText) GetTextOptions() SearchRequestTextOptions {
	if rt == nil {
		return SearchRequestTextOptions{}
	}
	return rt.TextOptions
}

func (rt SearchRequestTextOptions) GetIsNLP() bool {
	return rt.IsNLP
}

func (rt SearchRequestTextOptions) GetIsFuzzy() bool {
	return rt.IsFuzzy
}

func (rt SearchRequestTextOptions) GetIsAnd() bool {
	return rt.IsAnd
}

func (rt SearchRequestTextOptions) GetExcludes() []string {
	if rt.Excludes == nil {
		return []string{}
	}
	return rt.Excludes
}

type SearchRequestOption struct {
	Text    SearchRequestText `json:"text"`
	Image   []float64         `json:"image"`
	IsImage bool              `json:"is_image"`
}

type SearchRequest struct {
	from  int64
	to    int64
	query string
}

type SearchResponse struct {
	Images     []es.DocumentStructure `json:"images"`
	TotalCount int64                  `json:"total_count"`
}

type WebHandler struct {
	searcher *es.SearchClient
	memcache *memcache.Memcache
}

func NewWebHandler(searcher *es.SearchClient, memcache *memcache.Memcache) *WebHandler {
	return &WebHandler{searcher: searcher, memcache: memcache}
}

func validateAndProcessIDRequest(ctx context.Context, r *http.Request) (*SearchRequestID, error) {
	queryParameters := r.URL.Query()
	id := queryParameters.Get("id")

	searchRequest := &SearchRequestID{
		ID: id,
	}

	return searchRequest, nil
}

func validateAndProcessIDsRequest(ctx context.Context, r *http.Request) (*SearchRequestIDs, error) {
	queryParameters := r.URL.Query()
	idInputs := queryParameters.Get("ids")

	excludeInputsTrimmed := strings.Trim(idInputs, "")
	var ids []string
	if excludeInputsTrimmed != "" {
		ids = strings.Split(excludeInputsTrimmed, ",")
	}
	searchRequest := &SearchRequestIDs{
		IDs: ids,
	}

	return searchRequest, nil
}

func validateAndProcessImgRequest(ctx context.Context, r *http.Request) (*SearchRequestImage, error) {
	queryParameters := r.URL.Query()
	query := queryParameters.Get("q")

	searchRequest := &SearchRequestImage{
		Query: query,
	}

	return searchRequest, nil
}

func validateAndProcessRequest(ctx context.Context, r *http.Request) (*SearchRequestText, error) {
	queryParameters := r.URL.Query()
	query := queryParameters.Get("q")
	isFuzzyInput := queryParameters.Get("is_fuzzy")
	excludeInputs := queryParameters.Get("excludes")
	isAndInput := queryParameters.Get("is_and")

	isAnd, err := strconv.ParseBool(isAndInput)
	if err != nil {
		return nil, err
	}

	isFuzzy, err := strconv.ParseBool(isFuzzyInput)
	if err != nil {
		return nil, err
	}

	excludeInputsTrimmed := strings.Trim(excludeInputs, "")
	var excludes []string
	if excludeInputsTrimmed != "" {
		excludes = strings.Split(excludeInputsTrimmed, ",")
	}

	searchRequest := &SearchRequestText{
		Query: query,
		TextOptions: SearchRequestTextOptions{
			IsNLP:    false,
			IsFuzzy:  isFuzzy,
			Excludes: excludes,
			IsAnd:    isAnd,
		},
	}

	return searchRequest, nil
}

func (h *WebHandler) IDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	searchReq, err := validateAndProcessIDRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var (
		searchResults []es.DocumentStructure
		count         int64
	)

	searchResults = make([]es.DocumentStructure, 0)
	searchResults, count, err = h.searcher.SearchSimilarByID(ctx, searchReq.GetID())
	searchResp := &SearchResponse{
		Images:     searchResults,
		TotalCount: count,
	}

	jsonResponse, err := json.Marshal(searchResp)
	if err != nil {
		http.Error(w, "failed to serialize search results", http.StatusInternalServerError)
		return
	}

	log.Println("img request", searchReq, "response", searchResp.TotalCount)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *WebHandler) IDsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	searchReq, err := validateAndProcessIDsRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var (
		searchResults []es.DocumentStructure
		count         int64
	)

	searchResults = make([]es.DocumentStructure, 0)
	searchResults, count, err = h.searcher.SearchByIDs(ctx, searchReq.GetIDs())
	searchResp := &SearchResponse{
		Images:     searchResults,
		TotalCount: count,
	}

	jsonResponse, err := json.Marshal(searchResp)
	if err != nil {
		http.Error(w, "failed to serialize search results", http.StatusInternalServerError)
		return
	}

	log.Println("img request", searchReq, "response", searchResp.TotalCount)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func (h *WebHandler) ImageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	searchReq, err := validateAndProcessImgRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var (
		searchResults []es.DocumentStructure
		count         int64
	)

	searchResults = make([]es.DocumentStructure, 0)
	searchResults, count, err = h.searcher.SearchTextInImage(ctx, searchReq.GetQuery())
	searchResp := &SearchResponse{
		Images:     searchResults,
		TotalCount: count,
	}

	jsonResponse, err := json.Marshal(searchResp)
	if err != nil {
		http.Error(w, "failed to serialize search results", http.StatusInternalServerError)
		return
	}

	log.Println("img request", searchReq, "response", searchResp.TotalCount)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
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

	searchResults = make([]es.DocumentStructure, 0)
	if searchReq.GetTextOptions().GetIsFuzzy() {
		searchResults, count, err = h.searcher.SearchTextWithFuzzy(ctx, searchReq.GetQuery(), searchReq.GetTextOptions().GetIsAnd(),
			searchReq.GetTextOptions().GetExcludes())
	} else {
		if searchReq.GetTextOptions().GetIsAnd() {
			searchResults, count, err = h.searcher.SearchTextWithAnd(ctx, searchReq.GetQuery(),
				searchReq.GetTextOptions().GetExcludes())
		} else {
			searchResults, count, err = h.searcher.SearchTextNoFuzzy(ctx, searchReq.GetQuery(), searchReq.GetTextOptions().GetIsAnd(),
				searchReq.GetTextOptions().GetExcludes())
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

	log.Println("request", searchReq, "response", searchResp.TotalCount)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
