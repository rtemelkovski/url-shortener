package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/rtemelkovski/url-shortener/urlshort"
)

/*
decoderAdapter - Function type for converting a json or yaml specific Decoder to our own custom Decoder.
By adapting these types to our own interface, we can generalize all the Decoding work instead of writing
a seperate functions for both yaml and json. Seeing as all the Decoder types take an io.Reader as input,
we follow this trend.

Example:
func(reader io.Reader) urlshort.Decoder {
	decoder := yaml.NewDecoder(reader)
	return decoder
}
*/
type decoderAdapter func(io.Reader) urlshort.Decoder

/*
Extends the paths of a given redirector using the configs in the file at the given path. decoderAdapter
allows us to read in any type of file and convert it to our own custom Decoder. This way, parsing from
new types of files can be done easily by just creating a new decoderAdapter for that file type.
*/
func extendPaths(redirector http.HandlerFunc, path *string, adaptDecoder decoderAdapter) http.HandlerFunc {
	// attempt to open the file. Print an error message and skip this file's registration if there is an error
	reader, err := os.Open(*path)
	if err != nil {
		fmt.Printf("Cannot open file %s, skipping registration for this file\n", *path)
		return redirector
	}

	// attempt to decode the file and append the paths. Print an error message and skip this file's registration if there is an error
	newPaths, err := urlshort.DecodeToMap(adaptDecoder(reader))
	if err != nil {
		fmt.Printf("Cannot decode file %s, skipping registration for this file\n", *path)
		return redirector
	}

	return urlshort.Extend(redirector, newPaths)
}

/*
	Default handler function which just writes Hello with the url path
*/
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}

func main() {
	// read in flags to get the file names of the JSON and YAML files which contain our routes
	yamlPath := flag.String("YAML", "paths.yaml", "A YAML file")
	jsonPath := flag.String("JSON", "paths.json", "A JSON file")
	flag.Parse()

	// default handler in case a path doesn't exist in either config
	redirector := defaultHandler

	// attempt to read in the file given via YAML flag
	redirector = extendPaths(redirector, yamlPath, func(reader io.Reader) urlshort.Decoder {
		decoder := yaml.NewDecoder(reader)
		return decoder
	})

	// attempt to read in the file given via JSON flag
	redirector = extendPaths(redirector, jsonPath, func(reader io.Reader) urlshort.Decoder {
		decoder := json.NewDecoder(reader)
		return decoder
	})

	// set up handle func and run server
	http.HandleFunc("/", redirector)
	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Unable to start server", err)
	}
}
