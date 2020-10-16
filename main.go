package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rtemelkovski/url-shortener/urlshort"
	"gopkg.in/yaml.v2"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func getYAMLDecoder(yamlPath *string) (*yaml.Decoder, error) {
	reader, err := os.Open(*yamlPath)
	if err != nil {
		return nil, err
	}

	return yaml.NewDecoder(reader), nil
}

func getJSONDecoder(jsonPath *string) (*json.Decoder, error) {
	reader, err := os.Open(*jsonPath)
	if err != nil {
		return nil, err
	}

	return json.NewDecoder(reader), nil
}

func extendPaths(redirector http.HandlerFunc, decoder urlshort.Decoder) http.HandlerFunc {
	newPaths, err := urlshort.ToMap(decoder)
	if err != nil {
		fmt.Printf("Cannot decode file, skipping registration for this file")
		return redirector
	}

	return urlshort.Extend(redirector, newPaths)
}

func main() {
	yamlPath := flag.String("YAML", "paths.yaml", "A YAML file")
	jsonPath := flag.String("JSON", "paths.json", "A JSON file")
	flag.Parse()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	newMux := http.NewServeMux()
	redirector := hello
	// register all the paths
	redirector = urlshort.Extend(redirector, pathsToUrls)

	yamlDecoder, err := getYAMLDecoder(yamlPath)
	if err != nil {
		fmt.Println("Cannot open yaml file, skipping registration for this file")
	} else {
		redirector = extendPaths(redirector, yamlDecoder)
	}

	jsonDecoder, err := getJSONDecoder(jsonPath)
	if err != nil {
		fmt.Println("Cannot open json file, skipping registration for this file")
	} else {
		redirector = extendPaths(redirector, jsonDecoder)
	}

	newMux.HandleFunc("/", redirector)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", newMux); err != nil {
		log.Fatalln("Unable to start server", err)
	}
}
