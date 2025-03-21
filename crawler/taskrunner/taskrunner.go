//-------------------------------//
// Works through all crawl tasks //
//-------------------------------//

package taskrunner

import (
	"fmt"
	"sync"

	"github.com/ZakkBob/AskDave/crawler/fetcher"
	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/robots"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
)

type TaskRunner struct {
	Tasks   tasks.Tasks
	Results tasks.Results
	Fetcher fetcher.Fetcher
}

func (r *TaskRunner) Run(concurrency int) {
	go r.Results.ListenToChans()

	loopConcurrently(r.crawlNextRobots, concurrency)
	close(r.Results.RobotsChan)

	<-r.Results.RobotsFinished

	// loopConcurrently(r.crawlNextSitemap, concurrency)
	// close(r.Results.SitemapsChan)

	loopConcurrently(r.crawlNextPage, concurrency)
	close(r.Results.PagesChan)

	<-r.Results.PagesFinished
}

// Attempts to crawl the next robots.txt in the taskList, returns false if no more
func (r *TaskRunner) crawlNextRobots() bool {
	u, ok := r.Tasks.Robots.Next()
	if !ok {
		return false
	}

	robotsTxtUrl, _ := url.ParseRel("/robots.txt", u)
	robotsTxtUrl.TrailingSlash = false //normalise the url

	res, err := r.Fetcher.Fetch(robotsTxtUrl.String())
	if err != nil || res.StatusCode != 200 {
		robotsResult := tasks.RobotsResult{
			Url:           &robotsTxtUrl,
			Success:       false,
			FailureReason: tasks.FetchFailed,
		}

		r.Results.RobotsChan <- &robotsResult
		return true
	}

	validator, _ := robots.Parse(res.Body)
	robotsResult := tasks.RobotsResult{
		Url:       &robotsTxtUrl,
		Success:   true,
		Changed:   true,
		Hash:      hash.Hashs(res.Body),
		Validator: &validator,
	}

	r.Results.RobotsChan <- &robotsResult

	return true
}

// Attempts to crawl the next sitemap in the taskList, returns false if no more
// func (r *TaskRunner) crawlNextSitemap() bool {
// 	u := r.Tasks.Sitemaps.Next()
// 	if u == nil {
// 		return false
// 	}
// 	_, err := r.Fetcher.Fetch(u)
// 	if err != nil {
// 		return true
// 	}
// 	return true
// }

// Attempts to crawl the next page in the taskList, returns false if no more
func (r *TaskRunner) crawlNextPage() bool {
	u, ok := r.Tasks.Pages.Next()
	if !ok {
		return false
	}

	robotsAllowed, err := r.Results.CheckRobots(u)

	if err != nil {
		fmt.Printf("error: crawling next robots: %v\n", err)
		return true
	}

	if !robotsAllowed {
		pageResult := tasks.PageResult{
			Url:           &u,
			Success:       false,
			FailureReason: tasks.RobotsDisallowed,
		}

		r.Results.PagesChan <- &pageResult
		return true
	}

	res, err := r.Fetcher.Fetch(u.String())
	if err != nil {
		pageResult := tasks.PageResult{
			Url:           &u,
			Success:       false,
			FailureReason: tasks.FetchFailed,
		}

		r.Results.PagesChan <- &pageResult
		return true
	}

	p := page.Parse(res.Body, u)
	pageResult := tasks.PageResult{
		Url:     &u,
		Success: true,
		Changed: true,
		Page:    &p,
	}

	r.Results.PagesChan <- &pageResult

	return true
}

// Runs n goroutines, each calling t until it return false
func loopConcurrently(t func() bool, workers int) {
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				isNext := t()
				if !isNext {
					break
				}
			}
		}()
	}

	wg.Wait()
}
