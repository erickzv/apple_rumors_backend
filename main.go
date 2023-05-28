package main

import (
	"net/http"
	"os"
	"sync"

	"github.com/anaskhan96/soup"
	"github.com/goccy/go-json"
)

var websites [3]Website = [3]Website{
	{
		Domain:  "macrumors",
		FindTag: "h2",
	},
	{
		Domain:  "appleinsider",
		FindTag: "h2",
	},
	{
		Domain:  "9to5mac",
		FindTag: "h2",
	},
}

type Website struct {
	Domain  string
	FindTag string
}

// GetSoup requests Html from a website and returns all tags that match the findTag
func (w Website) GetSoup() (string, []soup.Root) {
	url := "https://" + w.Domain + ".com"
	res, err := soup.Get(url)
	if err != nil {
		panic(err)
	}
	doc := soup.HTMLParse(res)
	return url, doc.FindAll(w.FindTag)
}

// ParseSoup finds every <a> inside tags
func (w Website) ParseSoup(tags []soup.Root) []News {
	TagCount, maximum := 0, 16
	news := make([]News, 0, maximum)

	deleteAmount := 0

	for i := 0; i < len(tags) && TagCount < maximum; i++ {
		aTag := tags[i].Find("a")
		if aTag.Error != nil {
			continue
		}

		href := aTag.Attrs()["href"]

		if deleteAmount == 0 {
			cleanHref(href, &deleteAmount)
		}

		path := href[deleteAmount:]
		news = append(news, News{aTag.Text(), path})
		TagCount++
	}
	return news
}

func cleanHref(href string, amount *int) {
	slashCount := 0
	for _, character := range href {
		if string(character) == "/" {
			slashCount++
		}

		if slashCount == 3 {
			break
		}

		*amount++
	}
}

type News struct {
	Title string `json:"title"`
	Href  string `json:"href"`
}

func Scrape() map[string][]News {
	data := make(map[string][]News, len(websites))
	wg := sync.WaitGroup{}

	for i := 0; i < len(websites); i++ {
		website := websites[i]

		wg.Add(1)
		go func() {
			defer wg.Done()

			url, htmlSoup := website.GetSoup()
			data[url] = website.ParseSoup(htmlSoup)
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
