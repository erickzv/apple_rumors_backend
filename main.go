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

func ParseSoup(tags []soup.Root) []News {
	// finds every <tag> inside tags, we look for an <a>
	news := []News{}
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

func GetSoup(website Website) []soup.Root {
	res, err := soup.Get(website.Url)
	if err != nil {
		os.Exit(1)
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
	}

	var urls []string
	news := [][]News{}
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

func main() {
	http.HandleFunc("/all_news", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Scrape())
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8000" // Default port if not specified
	} else {
		port = ":" + port
	}
	http.ListenAndServe(port, nil)
}
