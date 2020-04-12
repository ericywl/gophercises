package urlshort

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if url, found := pathsToUrls[r.URL.Path]; found {
			http.Redirect(w, r, url, 308)
		} else {
			fallback.ServeHTTP(w, r)
		}
	})
}

// DBHandler with return an http.HandlerFunc that will attempt
// to map any paths (keys in the map) to their corresponding
// URL (values that each key in the map points to, in string format).
func DBHandler(db *bolt.DB, bucketName []byte, fallback http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		db.View(func(tx *bolt.Tx) error {
			c := tx.Bucket(bucketName).Cursor()
			key := []byte(r.URL.Path)
			if k, v := c.Seek(key); bytes.Equal(k, key) {
				http.Redirect(w, r, string(v), 308)
			} else {
				fallback.ServeHTTP(w, r)
			}
			return nil
		})
	})
}

// pathURL contains path and url pair
type pathURL struct {
	Path string `yaml:"path" json:"path"`
	URL  string `yaml:"url" json:"url"`
}

func buildMap(pathURLs []pathURL) map[string]string {
	urlMap := make(map[string]string)
	for _, x := range pathURLs {
		urlMap[x.Path] = x.URL
	}

	return urlMap
}

func parseYAML(yml []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := yaml.Unmarshal(yml, &pathURLs)
	if err != nil {
		return nil, err
	}

	return pathURLs, nil
}

func parseJSON(jsn []byte) ([]pathURL, error) {
	var pathURLs []pathURL
	err := json.Unmarshal(jsn, &pathURLs)
	if err != nil {
		return nil, err
	}

	return pathURLs, nil
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}

	urlMap := buildMap(pathURLs)
	return MapHandler(urlMap, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
func JSONHandler(jsn []byte, fallback http.Handler) (http.HandlerFunc, error) {
	pathURLs, err := parseJSON(jsn)
	if err != nil {
		return nil, err
	}

	urlMap := buildMap(pathURLs)
	return MapHandler(urlMap, fallback), nil
}
