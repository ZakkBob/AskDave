//-------------------------------//
// Works through all crawl tasks //
//-------------------------------//

package main

import (
// "fmt"
// "io"
// "net/http"
// "regexp"
// "slices"
// "sync"
// "time"
)

var crawledPages = safePageSlice{slice: make([]page, 0)}
var crawlTasks = safeCrawlTaskSlice{slice: make([]crawlTask, 0)}

// var start time.Time

// func getAbsoluteUrl(url string, pageUrl string) (absoluteUrl string) { //Converts relative urls to fully qualified
// 	if len(url) == 0 { //what the frick happened here
// 		return url
// 	}

// 	if url[0] != '/' {
// 		return url
// 	}

// 	siteUrl := pageUrl

// 	matched, _ := regexp.MatchString("(?i)https://.*?/.*", pageUrl)

// 	if matched { //if url pageUrl contains a directory
// 		siteUrlRegex := regexp.MustCompile("(?i)(https://.*?)/")
// 		matches := siteUrlRegex.FindStringSubmatch(pageUrl)

// 		if len(matches) < 2 {
// 			panic("what the freak!!")
// 		}

// 		siteUrl = matches[1]
// 	}
// 	absoluteUrl = siteUrl + url

// 	return absoluteUrl
// }

// func processDiscoveredUrl(url string) {
// 	discoveredUrls.mutex.Lock()
// 	if slices.Contains(discoveredUrls.slice, url) {
// 		discoveredUrls.mutex.Unlock()
// 		return
// 	}
// 	discoveredUrls.slice = append(discoveredUrls.slice, url)
// 	discoveredUrls.mutex.Unlock()

// 	uncrawledUrls.mutex.Lock()
// 	uncrawledUrls.slice = append(uncrawledUrls.slice, url)
// 	uncrawledUrls.mutex.Unlock()

// }

// func crawlNextUrl() {
// 	uncrawledUrls.mutex.Lock()
// 	uncrawledCount := len(uncrawledUrls.slice)
// 	if len(uncrawledUrls.slice) == 0 {
// 		uncrawledUrls.mutex.Unlock()
// 		return
// 	}

// 	nextUrl := uncrawledUrls.slice[0]
// 	uncrawledUrls.slice = slices.Delete(uncrawledUrls.slice, 0, 1)
// 	uncrawledUrls.mutex.Unlock()

// 	pageData := fetchPageData(nextUrl)

// 	for _, url := range pageData.links {
// 		processDiscoveredUrl(url)
// 	}

// 	crawledPages.mutex.Lock()
// 	crawledPages.slice = append(crawledPages.slice, pageData)
// 	crawledPages.mutex.Unlock()

// 	crawledUrls.mutex.Lock()
// 	crawledCount := len(crawledUrls.slice)
// 	fmt.Printf("\r%f crawls/s %d Crawled, Uncrawled %d   ", float64(crawledCount)/time.Since(start).Seconds(), crawledCount, uncrawledCount)
// 	crawledUrls.slice = append(crawledUrls.slice, nextUrl)
// 	crawledUrls.mutex.Unlock()

// }

// func autoCrawl(count int) {
// 	if count < 1 {
// 		return
// 	}
// 	crawlNextUrl()
// 	autoCrawl(count - 1)
// }

// func logCrawlStats() {
// 	crawledUrls.mutex.Lock()
// 	defer crawledUrls.mutex.Unlock()
// 	//uncrawledUrls.mutex.Lock()
// 	//defer uncrawledUrls.mutex.Unlock()

// 	fmt.Println(len(crawledUrls.slice), "Crawled")
// 	//fmt.Println(len(uncrawledUrls.slice), "Uncrawled")
// }
