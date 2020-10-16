package urlshort

import (
	"net/http"
)

type redirect struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

// Decoder Generic
type Decoder interface {
	Decode(v interface{}) (err error)
}

// Extend handler to have new endpoints
func Extend(startingHandler http.HandlerFunc, paths map[string]string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if paths[r.URL.Path] != "" {
			http.Redirect(rw, r, paths[r.URL.Path], http.StatusFound)
		}

		startingHandler.ServeHTTP(rw, r)
	}
}

// ToMap with Decoder
func ToMap(decoder Decoder) (map[string]string, error) {
	var redirects []redirect

	err := decoder.Decode(&redirects)
	if err != nil {
		return nil, err
	}

	paths := make(map[string]string)
	for _, path := range redirects {
		paths[path.Path] = path.URL
	}

	return paths, nil
}
