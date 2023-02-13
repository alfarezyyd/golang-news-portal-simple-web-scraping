package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"golang-news-portal-simple-web-scraping/helper"
	"os"
	"time"
)

type News struct {
	Id    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
	Link  string `json:"link,omitempty"`
	Type  string `json:"type,omitempty"`
}

func main() {
	collyScraper := colly.NewCollector()
	collyScraper.SetRequestTimeout(120 * time.Second)
	allNews := make([]News, 0)

	collyScraper.OnHTML("article", func(element *colly.HTMLElement) {
		element.ForEach("a", func(i int, e *colly.HTMLElement) {
			newsItem := News{}
			newsItem.Id = e.Attr("dtr-id")
			newsItem.Title = e.Attr("dtr-ttl")
			newsItem.Link = "https://www.cnnindonesia.com/ " + e.Attr("href")
			newsItem.Type = e.ChildText("span.kanal")
			allNews = append(allNews, newsItem)
		})
	})

	collyScraper.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})

	collyScraper.OnResponse(func(request *colly.Response) {
		fmt.Println("Got a response from", request.Request.URL)
	})

	collyScraper.OnError(func(response *colly.Response, err error) {
		fmt.Println("Error :", err)
	})

	collyScraper.OnScraped(func(response *colly.Response) {
		fmt.Println("Finished", response.Request.URL)
		js, err := json.MarshalIndent(allNews, "", "     ")
		helper.PanicIfError(err)
		fmt.Println("Writing data to file")
		err = os.WriteFile("all_news.json", js, 0664)
		helper.PanicIfError(err)
		fmt.Println("Data written to file successfully")
	})

	err := collyScraper.Visit("https://www.cnnindonesia.com/")
	helper.PanicIfError(err)
}
