package main

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
	"github.com/goccy/go-json"
)

type Website struct {
	Url     string
	FindTag string
}

// GetSoup requests Html from a website and returns all tags that match the findTag
func (w Website) GetSoup() []soup.Root {
	res, err := soup.Get(w.Url)
	if err != nil {
		panic(err)
	}
	doc := soup.HTMLParse(res)
	return doc.FindAll(w.FindTag)
}

// ParseSoup finds every <tag> inside tags, we look for <a> tag
func (w Website) ParseSoup(tags []soup.Root) []News {
	aTagCount, maximum := 0, 16
	news := make([]News, 0, maximum)

	for i := 0; i < len(tags) && aTagCount < maximum; i++ {
		aTag := tags[i].Find("a")
		urlLength := len(w.Url)
		if aTag.Error == nil {
			aTagCount++
			news = append(news, News{
				Title: aTag.Text(), Href: w.ParseHref(aTag.Attrs()["href"], urlLength),
			})
		}
	}
	return news
}

// ParseHref removes "http(s)://" & ".com"
func (w Website) ParseHref(href string, length int) string {
	parsed := href[length:]
	parsed = strings.TrimPrefix(parsed, ".com")
	return parsed
}

type News struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

func Scrape() map[string][]News {
	// FIXME Should be a global variable
	websites := [3]Website{
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
			FindTag: "h2",
		},
	}

	data := map[string][]News{}
	wg := sync.WaitGroup{}

	for i := 0; i < len(websites); i++ {
		website := websites[i]
		wg.Add(1)
		go func() {
			htmlSoup := website.GetSoup()
			data[website.Url] = website.ParseSoup(htmlSoup)
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
