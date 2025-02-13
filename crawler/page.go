//------------------------------//
// Processes page into a struct //
//------------------------------//

package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sync"
)

type page struct {
	url           string
	pageTitle     string
	ogTitle       string
	ogDescription string
	ogSiteName    string
	links         []url
}

type safePageSlice struct {
	mutex sync.Mutex
	slice []page
}

func (p *page) addLink(u url) {
	p.links = append(p.links, u)
}

func crawlPage(u url) (page, error) {
	b, err := fetchPageBody(u)

	if err != nil {
		return page{}, err
	}

	return parseBody(b, u)
}

func parseBody(b string, u url) (page, error) {
	var p page
	p.ogTitle = extractBodyMeta(b, "title")
	p.ogDescription = extractBodyMeta(b, "description")
	p.ogSiteName = extractBodyMeta(b, "site_name")

	p.pageTitle = extractBodyTitle(b)

	p.links = extractBodyLinks(b, u)

	return p, nil
}

func fetchPageBody(u url) (string, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func extractBodyMeta(page string, metaProperty string) (metaContent string) {
	metaElRegexString := fmt.Sprintf("(?s)<meta[^>]*?property=\"og:%s\"[^>]*?>", metaProperty) //Temporary fix, won't work if content contains a '>'

	metaElRegex := regexp.MustCompile(metaElRegexString)
	metaContentRegex := regexp.MustCompile("(?s)content=\"(.*?)\"")

	elMatches := metaElRegex.FindStringSubmatch(page)
	if len(elMatches) < 1 {
		return ""
	}

	metaEl := elMatches[0]

	contentMatches := metaContentRegex.FindStringSubmatch(metaEl)
	if len(contentMatches) < 2 {
		return ""
	}

	metaContent = contentMatches[1]
	return metaContent
}

func extractBodyTitle(page string) (pageTitle string) {
	pageTitleRegex := regexp.MustCompile("(?s)<title.*?>(.*?)</title>") //Temporary, won't match if space are in the tags :(
	matches := pageTitleRegex.FindStringSubmatch(page)

	if len(matches) < 2 {
		return ""
	}

	pageTitle = matches[1]
	return pageTitle
}

func extractBodyLinks(body string, pageUrl url) (pageLinks []url) {
	linkElRegex := regexp.MustCompile("(?s)<a.*?>") //Wont match if '>' is in the tag somewher :shruggie:
	linkHrefRegex := regexp.MustCompile("(?s)href=\"(.*?)\"")

	elMatches := linkElRegex.FindAllString(body, -1)
	if len(elMatches) < 1 {
		return []url{}
	}

	for _, linkEl := range elMatches {
		hrefMatches := linkHrefRegex.FindStringSubmatch(linkEl)
		if len(hrefMatches) < 1 {
			continue
		}

		linkUrl, err := parseRelativeUrl(hrefMatches[1], pageUrl)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		pageLinks = append(pageLinks, linkUrl)
	}

	return pageLinks
}
