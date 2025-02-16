//-------------------------------//
// Works through all crawl tasks //
//-------------------------------//

package tasks

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"sync"
)

type TaskRunner struct {
	Tasks   Tasks
	Results Results
	Fetcher fetcher.Fetcher
}

func (r *TaskRunner) Run(concurrency int) {
	loopConcurrently(r.crawlNextRobots, concurrency)
	loopConcurrently(r.crawlNextSitemap, concurrency)
	loopConcurrently(r.crawlNextPage, concurrency)
}

// Attempts to crawl the next robots.txt in the taskList, returns false if no more
func (r *TaskRunner) crawlNextRobots() bool {
	u := r.Tasks.Robots.Next()
	if u == nil {
		return false
	}
	r.Fetcher.Fetch(u)

	return true
}

// Attempts to crawl the next sitemap in the taskList, returns false if no more
func (r *TaskRunner) crawlNextSitemap() bool {
	u := r.Tasks.Sitemaps.Next()
	if u == nil {
		return false
	}
	r.Fetcher.Fetch(u)
	return true
}

// Attempts to crawl the next page in the taskList, returns false if no more
func (r *TaskRunner) crawlNextPage() bool {
	u := r.Tasks.Pages.Next()
	if u == nil {
		return false
	}
	r.Fetcher.Fetch(u)
	return true
}

// Runs n goroutines, each calling t until it return false
func loopConcurrently(t func() bool, workers int) {
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t() {
			}
		}()
	}

	wg.Wait()
}
