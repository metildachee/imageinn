package server

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"testing"
)

func Test_SearchByKeyword(t *testing.T) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		t.FailNow()
	}

	ctx := context.Background()
	keywords := []string{"example"}

	results, totalHits, err := SearchByKeyword(ctx, client, keywords)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println("len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchByCategory(t *testing.T) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		t.FailNow()
	}

	ctx := context.Background()
	results, totalHits, err := SearchByCategoryID(ctx, client, 100)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println("len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchByCategories(t *testing.T) {
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		t.FailNow()
	}

	ctx := context.Background()
	categoryIDs := []int64{100, 101}
	results, totalHits, err := SearchByCategoryIDs(ctx, client, categoryIDs)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println("len(results)", totalHits)
	fmt.Println(results)
}
