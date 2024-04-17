package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7" // imports as package "elastic"
	"log"
)

type DocumentStructure struct {
	URL         string  `json:"url"`
	Caption     string  `json:"caption"`
	ID          string  `json:"id"`
	CategoryIDs []int64 `json:"category_ids"`
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

func SearchByCategoryID(ctx context.Context, client *elastic.Client, categoryID int64) ([]DocumentStructure, int64, error) {
	termQuery := elastic.NewTermQuery("category_ids", categoryID)
	searchResult, err := client.Search().
		Index("images").
		Query(termQuery).
		Sort("id", true).
		From(0).Size(10).
		Pretty(true).
		Do(ctx)
	if err != nil {
		return nil, 0, err
	}

	totalHits := searchResult.TotalHits()

	hits, unmarshalErr := unmarshalResults(searchResult.Hits.Hits)
	return hits, totalHits, unmarshalErr
}

func SearchByCategoryIDs(ctx context.Context, client *elastic.Client, categoryIDs []int64) ([]DocumentStructure, int64, error) {
	andCategoryQuery := elastic.NewBoolQuery()
	for _, category := range categoryIDs {
		andCategoryQuery.Must(elastic.NewTermQuery("category_ids", category))
	}

	//orCategoryQuery := elastic.NewBoolQuery()
	//for _, category := range categoryIDs {
	//	orCategoryQuery.Should(elastic.NewTermQuery("categories", category))
	//}

	searchResult, err := client.Search().
		Index("images").
		Query(andCategoryQuery).
		Size(10).
		Do(ctx)
	if err != nil {
		return nil, 0, err
	}
	totalHits := searchResult.TotalHits()

	hits, unmarshalErr := unmarshalResults(searchResult.Hits.Hits)
	return hits, totalHits, unmarshalErr
}

func SearchByKeyword(ctx context.Context, client *elastic.Client, keywords []string) ([]DocumentStructure, int64, error) {
	andQuery := elastic.NewBoolQuery()
	for _, keyword := range keywords {
		andQuery.Must(elastic.NewMatchQuery("caption", keyword))
	}

	//orQuery := elastic.NewBoolQuery()
	//for _, keyword := range keywords {
	//	orQuery.Should(elastic.NewMatchQuery("caption", keyword))
	//}

	searchResult, err := client.Search().
		Index("images").
		Query(andQuery).
		Size(10).
		Do(ctx)
	if err != nil {
		return nil, 0, err
	}

	totalHits := searchResult.TotalHits()

	hits, unmarshalErr := unmarshalResults(searchResult.Hits.Hits)
	return hits, totalHits, unmarshalErr
}

func main() {
	// Create a client and connect to http://127.0.0.1:9200
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetBasicAuth("elastic", "rTur2xp5QmFAvEAPJjfT"),
	)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	ctx := context.Background()

	termQuery := elastic.NewTermQuery("category_ids", 100)
	searchResult, err := client.Search().
		Index("images").
		Query(termQuery).
		Sort("id", true).
		From(0).Size(10).
		Pretty(true).
		Do(ctx)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	fmt.Printf("Found a total of %d documents\n", searchResult.TotalHits())

	// Iterate through results
	for _, hit := range searchResult.Hits.Hits {
		// Deserialize hit.Source into a Document (could also be just a map[string]interface{}).
		var t struct {
			URL     string `json:"url"`
			Caption string `json:"caption"`
		}
		err := json.Unmarshal(hit.Source, &t)
		if err != nil {
			// Deserialization failed
			log.Fatalf("Deserialization failed: %s", err)
		}

		fmt.Printf("Document by URL: %s, Caption: %s\n", t.URL, t.Caption)
	}
}
