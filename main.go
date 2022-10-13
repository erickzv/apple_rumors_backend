package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
)

type Website struct {
	Url     string
	FindTag string
}

type News struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

// ParseHref removes "http(s)://" & ".com"
func ParseHref(href string, length int) string {
	parsed := href[length:]
	parsed = strings.TrimPrefix(parsed, ".com")
	return parsed
}

// ParseSoup finds every <tag> inside tags, we look for <a> tag
func ParseSoup(tags []soup.Root, length int) []News {
	var news []News
	aTagCount := 0
	for i := 0; i < len(tags) && aTagCount < 10; i++ {
		aTag := tags[i].Find("a")
		if aTag.Error == nil {
			aTagCount++
			news = append(news, News{
				Title: aTag.Text(), Href: ParseHref(aTag.Attrs()["href"], length),
			})
		}
	}
	return news
}

// GetSoup requests Html from a website and returns all tags that match the findTag
func GetSoup(website Website) []soup.Root {
	res, err := soup.Get(website.Url)
	if err != nil {
		panic(err)
	}
	doc := soup.HTMLParse(res)
	return doc.FindAll(website.FindTag)
}

func Scrape() map[string][]News {
	websites := []Website{
		{
			Url:     "https://macrumors.com",
			FindTag: "h2",
		},
		{
			Url:     "https://appleinsider.com",
			FindTag: "h2",
		},
		{
			Url:     "https://9to5mac.com",
			FindTag: "h1",
		},
	}

	data := make(map[string][]News)
	var wg = sync.WaitGroup{}

	for i := 0; i < len(websites); i++ {
		website := websites[i]
		wg.Add(1)
		go func() {
			soup := GetSoup(website)
			data[website.Url] = ParseSoup(soup, len(website.Url))
			wg.Done()
		}()
	}

	wg.Wait()
	return data
}

// Port defines a port to listen on
func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

func main() {
	http.HandleFunc("/all_news", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(Scrape())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	err := http.ListenAndServe(Port(), nil)
	if err != nil {
		panic(err)
	}
}
