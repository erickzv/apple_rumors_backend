package main

import (
	"encoding/json"
	"net/http"
	"os"
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

type ApiRes struct {
	News     [][]News `json:"news"`
	Websites []string `json:"websites"`
}

// finds every <tag> inside tags, we look for an <a>
func ParseSoup(tags []soup.Root) []News {
	var news []News
	aTagCount := 0
	for i := 0; i < len(tags) && aTagCount < 10; i++ {
		aTag := tags[i].Find("a")
		if aTag.Error == nil {
			aTagCount++
			news = append(news, News{
				Title: aTag.Text(), Href: aTag.Attrs()["href"],
			})
		}
	}
	return news
}

// Requests Html from a website and returns all tags that match the findTag
func GetSoup(website Website) []soup.Root {
	res, err := soup.Get(website.Url)
	if err != nil {
		panic(err)
	}
	doc := soup.HTMLParse(res)
	return doc.FindAll(website.FindTag)
}

func ParseUrl(url string) string {
	website := ""
	for i := 8; i < len(url)-4; i++ {
		website += string(url[i])
	}
	return website
}

func Scrape() ApiRes {
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
		{
			Url:     "https://machash.com",
			FindTag: "h2",
		},
	}

	var urls []string
	var news [][]News
	var wg = sync.WaitGroup{}

	for i := 0; i < len(websites); i++ {
		website := websites[i]
		wg.Add(1)
		go func() {
			s := GetSoup(website)
			news = append(news, ParseSoup(s))
			urls = append(urls, ParseUrl(website.Url))
			wg.Done()
		}()
	}

	wg.Wait()
	return ApiRes{News: news, Websites: urls}
}

// Defines a port to listen on
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
