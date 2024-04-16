package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7" // imports as package "elastic"
	"log"
)

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

	// Use a context to manage cancellation
	ctx := context.Background()

	// Perform a search for a specific URL in the "my_index" index
	termQuery := elastic.NewTermQuery("category_ids", 100)
	searchResult, err := client.Search().
		Index("images").  // search in index "my_index
		Query(termQuery). // specify the query
		Sort("id", true). // sort by "id" field, ascending
		From(0).Size(10). // take documents 0-9
		Pretty(true).     // pretty print request and response JSON
		Do(ctx)           // execute
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
