package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/ericywl/gophercises/quiethn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("./quiethn/index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

type result struct {
	idx  int
	item item
	err  error
}

func getStories(ids []int) []item {
	client := hn.NewClient()
	resultCh := make(chan result)

	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
				return
			}
			resultCh <- result{idx: idx, item: parseHNItem(idx, hnItem)}
		}(i, ids[i])
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}

		if isStoryLink(res.item) {
			stories = append(stories, res.item)
		}
	}

	return stories
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func getTopStories(numStories int) ([]item, error) {
	client := hn.NewClient()
	ids, err := client.TopItems()
	if err != nil {
		return nil, err
	}

	var stories []item
	var at int
	for len(stories) < numStories {
		if at >= len(ids)-1 {
			break
		}
		need := (numStories - len(stories)) * 5 / 4
		end := min(len(ids), at+need)
		stories = append(stories, getStories(ids[at:end])...)
		at = end
	}

	if len(stories) < numStories {
		return stories, nil
	}
	return stories[:numStories], nil
}

type storyCache struct {
	stories    []item
	timeout    time.Duration
	expiration time.Time
	lock       sync.Mutex
}

func (c *storyCache) getCachedTopStories(numStories int) ([]item, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if time.Now().Sub(c.expiration) < 0 {
		return c.stories, nil
	}

	stories, err := getTopStories(numStories)
	if err != nil {
		return nil, err
	}

	c.expiration = time.Now().Add(c.timeout)
	c.stories = stories
	return c.stories, nil
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	timeout := 4 * time.Second
	cache := storyCache{timeout: timeout}

	go func() {
		ticker := time.NewTicker(timeout / 2)
		for {
			tempCache := storyCache{timeout: timeout}
			tempCache.getCachedTopStories(numStories)
			cache.lock.Lock()
			cache.stories = tempCache.stories
			cache.expiration = tempCache.expiration
			cache.lock.Unlock()
			<-ticker.C
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := cache.getCachedTopStories(numStories)
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(idx int, hnItem hn.Item) item {
	ret := item{Item: hnItem}
	retUrl, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(retUrl.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
