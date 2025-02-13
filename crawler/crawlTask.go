//------------------------------------//
// Represents a task provided by dave //
//------------------------------------//

package main

import (
	"sync"
)

type crawlTaskType int

const (
	pageCrawlTask crawlTaskType = iota
	robotsCrawlTask
	sitemapCrawlTask
)

type crawlTask struct {
	crawlType crawlTaskType
	url       url
}

type safeCrawlTaskSlice struct {
	mutex sync.Mutex
	slice []crawlTask
}
