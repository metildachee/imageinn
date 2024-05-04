package main

import (
	"github.com/gorilla/mux"
	config "github.com/metildachee/imageinn/server/config"
	"github.com/metildachee/imageinn/server/es"
	"github.com/metildachee/imageinn/server/handler"
	"log"
	"net/http"
)

func main() {
	configPath := "config/config.yml" // change into flag
	serverConfig := config.LoadConfig(configPath)
	searcher, err := es.NewSearcher(*serverConfig)
	if err != nil {
		log.Fatalln("load searcher failed", err)
	}
	webHandler := handler.NewWebHandler(searcher)

	r := mux.NewRouter()
	r.HandleFunc("/search", webHandler.SearchHandler).Methods("GET")

	log.Println("Starting server on :8080")
	if httpError := http.ListenAndServe(":8080", r); httpError != nil {
		log.Fatal(httpError)
	}
}
