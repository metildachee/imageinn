package es

import (
	"context"
	"fmt"
	config2 "github.com/metildachee/imageinn/server/config"
	"testing"
)

func getTestElasticClient(t *testing.T) *SearchClient {
	conf := config2.LoadConfig("../config/config.yml")
	searcher, err := NewSearcher(*conf)
	if err != nil {
		t.FailNow()
		t.FailNow()
	}
	return searcher
}

func Test_SearchText(t *testing.T) {
	funcName := "SearchText"

	ctx := context.Background()
	searcher := getTestElasticClient(t)

	results, totalHits, err := searcher.SearchText(ctx, "iron story")
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchTextWithExclusions(t *testing.T) {
	funcName := "Test_SearchTextWithExclusions"

	ctx := context.Background()
	searcher := getTestElasticClient(t)

	results, totalHits, err := searcher.SearchTextWithExclusions(ctx, "iron", []string{"man"})
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchTextWithFuzzy(t *testing.T) {
	funcName := "Test_SearchTextWithFuzzy"

	ctx := context.Background()
	searcher := getTestElasticClient(t)

	results, totalHits, err := searcher.SearchTextWithFuzzy(ctx, "mountains of the moon", false, []string{})
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}

func Test_SearchTextNoFuzzy(t *testing.T) {
	funcName := "Test_SearchTextNoFuzzy"

	ctx := context.Background()
	searcher := getTestElasticClient(t)

	results, totalHits, err := searcher.SearchTextNoFuzzy(ctx, "mountains of the moon", false, []string{})
	if err != nil {
		fmt.Println("got error", err)
		t.FailNow()
	}

	fmt.Println(funcName, "len(results)", totalHits)
	fmt.Println(results)
}
