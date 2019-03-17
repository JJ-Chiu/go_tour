package main

import (
	"fmt"
	"sync"
)

type WebCrawler struct {
	// Cache fetched url, return its body.
	cache map[string]string
	// mutex used to protect cache.
	mux sync.Mutex
}

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func NewWebCrawler() *WebCrawler {
	return &WebCrawler{
		cache: make(map[string]string),
	}
}

// Add url, body into cache if url isn't cached.
// Return false if url already cached.
func (wc *WebCrawler) Add(url string, body string) bool {
	wc.mux.Lock()
	defer wc.mux.Unlock()

	_, hit := wc.cache[url]
	if !hit {
		wc.cache[url] = body
	}
	return !hit
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func (wc *WebCrawler) Crawl(url string, depth int, fetcher Fetcher) {
	if depth <= 0 || !wc.Add(url, "") {
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)
	// It's safe to update the cache without lock it.
	// Because the slot of this url is already set by wc.Add()
	wc.cache[url] = body

	wg := sync.WaitGroup{}
	for _, u := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			wc.Crawl(url, depth-1, fetcher)
		}(u)
	}
	wg.Wait()
	return
}
