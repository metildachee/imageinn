package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	config "github.com/metildachee/imageinn/server/config"
	"github.com/metildachee/imageinn/server/es"
	"github.com/metildachee/imageinn/server/handler"
	"log"
	"net/http"
)

func main() {
	configPath := "config/config.yml" // TODO: Change into flag
	serverConfig := config.LoadConfig(configPath)
	searcher, err := es.NewSearcher(*serverConfig)
	if err != nil {
		log.Fatalln("load searcher failed", err)
	}
	webHandler := handler.NewWebHandler(searcher)

	r := mux.NewRouter()
	r.HandleFunc("/search", webHandler.SearchHandler).Methods("GET")

	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // Adjust to match your requirement
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "X-Requested-With"}),
		handlers.ExposedHeaders([]string{"X-My-Custom-Header"}),
		handlers.AllowCredentials(),
	)
	handler := corsObj(r)

	log.Println("Starting server on :8080")
	if httpError := http.ListenAndServe(":8080", handler); httpError != nil {
		log.Fatal(httpError)
	}
}
