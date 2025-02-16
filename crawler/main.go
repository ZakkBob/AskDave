package main

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"ZakkBob/AskDave/crawler/tasks"
	"ZakkBob/AskDave/crawler/url"
	"fmt"
	"strconv"
	"time"
)

func main() {
	r := tasks.TaskRunner{
		Fetcher: &fetcher.DummyFetcher{
			Response:  "Dummy",
			Delay:     time.Second,
			RandDelay: time.Second * 2,
			Debug:     true,
		},
	}

	for i := 0; i < 2; i++ {
		s := "https://example.com/robots.txt" + strconv.Itoa(i)

		u, _ := url.ParseAbs(s)
		r.Tasks.Robots.Slice = append(r.Tasks.Robots.Slice, u)
	}

	for i := 0; i < 2; i++ {
		u, _ := url.ParseAbs(fmt.Sprintf("https://example.com/sitemap%d.xml", i))
		r.Tasks.Sitemaps.Slice = append(r.Tasks.Sitemaps.Slice, u)
	}

	for i := 0; i < 2; i++ {
		u, _ := url.ParseAbs(fmt.Sprintf("https://example.com/page%d.html", i))
		r.Tasks.Pages.Slice = append(r.Tasks.Pages.Slice, u)
	}

	r.Run(50)
}
