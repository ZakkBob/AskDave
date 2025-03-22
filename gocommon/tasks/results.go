package tasks

import (
	"fmt"
	"sync"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/robots"
	"github.com/ZakkBob/AskDave/gocommon/url"
)

type FailureReason = int

const (
	NoFailure FailureReason = iota
	RobotsDisallowed
	FetchFailed
)

type Results struct {
	Robots         map[string]*RobotsResult `json:"robots,omitempty"`
	RobotsChan     chan *RobotsResult       `json:"-"`
	RobotsFinished chan bool                `json:"-"`
	robotMu        sync.RWMutex             `json:"-"`
	Pages          map[string]*PageResult   `json:"pages,omitempty"`
	PagesChan      chan *PageResult         `json:"-"`
	PagesFinished  chan bool                `json:"-"`
	// Sitemaps     map[string]*string       `json:"sitemaps,omitempty"`
	// SitemapsChan chan *string           `json:"-"`
}

type PageResult struct {
	Url           *url.Url      `json:"url"`
	Success       bool          `json:"success"`
	FailureReason FailureReason `json:"failure_reason,omitempty"`
	Changed       bool          `json:"changed,omitempty"`
	Page          *page.Page    `json:"page,omitempty"`
}

type RobotsResult struct {
	Url           *url.Url             `json:"robots"`
	Success       bool                 `json:"success"`
	FailureReason FailureReason        `json:"failure_reason,omitempty"`
	Hash          hash.Hash            `json:"hash,omitempty"`
	Changed       bool                 `json:"changed,omitempty"`
	Validator     *robots.UrlValidator `json:"validator,omitempty"`
}

func (r *Results) CheckRobots(u url.Url) (bool, error) {
	robotsUrl, err := url.ParseRel("/robots.txt", u)
	robotsUrl.TrailingSlash = false //normalise the url

	if err != nil {
		return false, fmt.Errorf("checking robots: %w", err)
	}

	r.robotMu.RLock()
	defer r.robotMu.RUnlock()
	robotResult, ok := r.Robots[robotsUrl.String()]
	if !ok {
		return true, nil
	}

	if !robotResult.Success { //robots.txt couldnt be fetched
		return true, nil
	}

	return robotResult.Validator.ValidateUrl(&u), nil
}

func (r *Results) ListenToChans() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			robotsResult, ok := <-r.RobotsChan
			if !ok && robotsResult == nil {
				r.RobotsFinished <- true
				return
			}
			r.robotMu.Lock()
			r.Robots[robotsResult.Url.String()] = robotsResult
			r.robotMu.Unlock()
		}
	}()
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	for {
	// 		sitemapResult, ok := <-r.SitemapsChan
	// 		if !ok && sitemapResult == nil {
	// 			return
	// 		}
	// 		r.Sitemaps[sitemapResult.Url.String()] = sitemapResult
	// 	}
	// }()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			pageResult, ok := <-r.PagesChan
			if !ok && pageResult == nil {
				r.PagesFinished <- true
				return
			}
			r.Pages[pageResult.Url.String()] = pageResult
		}
	}()
	wg.Wait()
}
