package main

import (
	"context"
	"fmt"
	config2 "github.com/metildachee/imageinn/server/config"
	"testing"
)

func getTestElasticClient(t *testing.T) *Searcher {
	conf := config2.LoadConfig("config/config.yml")
	searcher, err := NewSearcher(*conf)
	if err != nil {
		t.FailNow()
	}
	return searcher
}

func Test_SearchByKeyword_OneHit(t *testing.T) {
	funcName := "Test_SearchByKeyword_OneHit"

	ctx := context.Background()
	keywords := []string{"example"}
	searcher := getTestElasticClient(t)

	results, totalHits, err := searcher.SearchByKeywordsAndOr(ctx, keywords)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchByKeyword_PartialKeywordHit(t *testing.T) {
	funcName := "Test_SearchByKeyword_PartialKeywordHit"

	ctx := context.Background()
	keywords := []string{"example", "apple"}
	searcher := getTestElasticClient(t)

	results, totalHits, err := searcher.SearchByKeywordsAndOr(ctx, keywords)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchByKeyword_NilHit(t *testing.T) {
	funcName := "Test_SearchByKeyword_NilHit"

	ctx := context.Background()
	keywords := []string{"apple"}
	searcher := getTestElasticClient(t)

	results, totalHits, err := searcher.SearchByKeywordsAndOr(ctx, keywords)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	if totalHits != 0 {
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchByOneCategory(t *testing.T) {
	funcName := "Test_SearchByCategory"
	searcher := getTestElasticClient(t)

	ctx := context.Background()
	results, totalHits, err := searcher.SearchByCategoryID(ctx, 100)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchByPartialCategory(t *testing.T) {
	funcName := "Test_SearchByPartialCategory"
	searcher := getTestElasticClient(t)

	ctx := context.Background()
	categories := []int64{100, 88}

	results, totalHits, err := searcher.SearchByCategoryIDsAndOr(ctx, categories)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchByAllCategory(t *testing.T) {
	funcName := "Test_SearchByAllCategory"
	searcher := getTestElasticClient(t)

	ctx := context.Background()
	categories := []int64{100, 101}

	results, totalHits, err := searcher.SearchByCategoryIDsAndOr(ctx, categories)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchByCategories(t *testing.T) {
	funcName := "Test_SearchByCategories"
	searcher := getTestElasticClient(t)

	ctx := context.Background()
	categoryIDs := []int64{100, 101}
	results, totalHits, err := searcher.SearchByCategoryIDsAndOr(ctx, categoryIDs)
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}
