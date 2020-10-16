package urlshort

import (
	"net/http"
)

// struct representing the URL path searched on our server and the location URL we want to redirect to
type redirect struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

/*
Decoder - Shared interface for Decoders. This allows us to write functions that work for all Decoder types instead
of copying the same functions for json and yaml decoders.
*/
type Decoder interface {
	Decode(v interface{}) (err error)
}

/*
Extend the startingHandler using the supplied paths. Requests with URL paths that are keys in the map will be redirected
to the values for their respective keys. If the URL path is not in the map, we pass the request on to the next handler.
*/
func Extend(startingHandler http.HandlerFunc, paths map[string]string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// if the url path is a key in the map
		if paths[r.URL.Path] != "" {
			// redirect to the value store for that key
			http.Redirect(rw, r, paths[r.URL.Path], http.StatusFound)
		}

		// pass request on to the next handler
		startingHandler.ServeHTTP(rw, r)
	}
}

/*
DecodeToMap - Given an incoming Decoder of any file type, decode it's contents into a map of the following format:
	key - request url path (IE the url path the user will search)
	value - redirect url (IE the url we will redirect to)
*/
func DecodeToMap(decoder Decoder) (map[string]string, error) {
	var redirects []redirect

	// attempt to decode into an array of redirects, return nil with error if an error is encountered
	err := decoder.Decode(&redirects)
	if err != nil {
		return nil, err
	}

	// iterate through the redirects and insert them as key-value pairs in the map
	paths := make(map[string]string)
	for _, path := range redirects {
		paths[path.Path] = path.URL
	}

	return paths, nil
}
