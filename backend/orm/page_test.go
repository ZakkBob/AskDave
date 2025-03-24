package orm_test

import (
	"testing"
	"time"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func comparePageData(t *testing.T, expected *page.Page, actual *page.Page) {
	assert.Equal(t, expected.Url.String(), actual.Url.String(), "Page data should have correct URL")
	assert.Equal(t, expected.Title, actual.Title, "Page data should have correct Title")
	assert.Equal(t, expected.OgTitle, actual.OgTitle, "Page data should have correct OgTitle")
	assert.Equal(t, expected.OgDescription, actual.OgDescription, "Page data should have correct OgDescription")
	assert.Equal(t, expected.OgSiteName, actual.OgSiteName, "Page data should have correct OgSiteName")
	assert.Equal(t, expected.Hash.String(), actual.Hash.String(), "Page data should have correct Hash string")
	assert.Equal(t, expected.Links, actual.Links, "Page data should have correct links")
}

func compareCreatedAndRetrievedPages(t *testing.T, expected *orm.OrmPage, actual *orm.OrmPage) {
	assert.Equal(t, expected.ID(), actual.ID(), "Retrieved page should have correct ID")
	assert.Equal(t, expected.SiteID(), actual.SiteID(), "Retrieved site should have the correct site ID")

	comparePageData(t, &expected.Page, &actual.Page)

	assert.WithinDuration(t, expected.NextCrawl, actual.NextCrawl, time.Second, "Retrieved page's next crawl should be reasonably close to expected")
	assert.Equal(t, expected.CrawlInterval, actual.CrawlInterval, "Retrieved page should have the correct crawl interval")
	assert.Equal(t, expected.IntervalDelta, actual.IntervalDelta, "Retrieved page should have the correct interval delta")
	assert.Equal(t, expected.Assigned, actual.Assigned, "Retrieved page should have the correct assigned value")
}

func TestCreatePage_And_PageByUrl_And_PageByID(t *testing.T) {
	u := uniqueURL(t)
	p := page.Page{
		Url:           u,
		Title:         "title",
		OgTitle:       "og title",
		OgDescription: "og description",
		OgSiteName:    "og site name",
		Hash:          hash.Hashs("hashed"),
	}

	now := time.Now()
	nextCrawl := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location()).AddDate(0, 1, 2)
	crawlInterval := 30
	intervalDelta := -2
	assigned := true

	createdPage, err := orm.CreatePage(p, nextCrawl, crawlInterval, intervalDelta, assigned)
	require.NoError(t, err, "CreatePage should not return an error")

	assert.NotZero(t, createdPage.ID(), "Returned OrmPage should have a non-zero ID")
	assert.NotZero(t, createdPage.SiteID(), "Returned OrmPage should have a non-zero Site ID")

	comparePageData(t, &p, &createdPage.Page)

	assert.WithinDuration(t, nextCrawl, createdPage.NextCrawl, time.Second, "Returned OrmPage's next crawl should be reasonably close to expected")
	assert.Equal(t, crawlInterval, createdPage.CrawlInterval, "Returned OrmPage should have the correct crawl interval")
	assert.Equal(t, intervalDelta, createdPage.IntervalDelta, "Returned OrmPage should have the correct interval delta")
	assert.Equal(t, assigned, createdPage.Assigned, "Returned OrmPage should have the correct assigned value")

	t.Run("PageByUrl", func(t *testing.T) {
		retrievedPage, err := orm.PageByUrl(u)
		require.NoError(t, err, "PageByUrl should not return an error")

		compareCreatedAndRetrievedPages(t, &createdPage, &retrievedPage)
	})

	t.Run("SiteByID", func(t *testing.T) {
		retrievedPage, err := orm.PageByID(createdPage.ID())
		require.NoError(t, err, "PageByID should not return an error")

		compareCreatedAndRetrievedPages(t, &createdPage, &retrievedPage)
	})
}

func TestPageSave(t *testing.T) {
	u := uniqueURL(t)
	p := page.Page{
		Url:           u,
		Title:         "title",
		OgTitle:       "og title",
		OgDescription: "og description",
		OgSiteName:    "og site name",
		Hash:          hash.Hashs("hashed"),
	}

	createdPage, err := orm.CreatePage(p, time.Now(), 2, -1, false)
	require.NoError(t, err, "CreatePage should not return an error")

	updatedUrl, err := url.ParseAbs("https://testpagesave-updatedurl.com")
	require.NoError(t, err, "ParseAbs should not return an error")

	now := time.Now()

	newPageData := page.Page{
		Url:           updatedUrl,
		Title:         "updated title",
		OgTitle:       "updated og title",
		OgDescription: "updated og description",
		OgSiteName:    "updated og site name",
		Hash:          hash.Hashs("upadated hash"),
	}

	createdPage.Page = newPageData

	nextCrawl := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location()).AddDate(0, 1, 2)
	crawlInterval := 3
	intervalDelta := -5
	assigned := true

	createdPage.NextCrawl = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location()).AddDate(0, 1, 2)
	createdPage.CrawlInterval = crawlInterval
	createdPage.IntervalDelta = intervalDelta
	createdPage.Assigned = assigned

	createdPage.Save()

	savedPage, err := orm.PageByID(createdPage.ID())
	require.NoError(t, err, "PageByID should not return an error")

	comparePageData(t, &newPageData, &savedPage.Page)

	assert.WithinDuration(t, nextCrawl, savedPage.NextCrawl, time.Second, "Retrieved page's next crawl should be reasonably close to expected")
	assert.Equal(t, crawlInterval, savedPage.CrawlInterval, "Retrieved page should have the correct crawl interval")
	assert.Equal(t, intervalDelta, savedPage.IntervalDelta, "Retrieved page should have the correct interval delta")
	assert.Equal(t, assigned, savedPage.Assigned, "Retrieved page should have the correct assigned value")

}

// func TestCreateEmptySite(t *testing.T) {
// 	u := uniqueURL(t)

// 	defaultTime, err := time.Parse("2006-01-02", "0000-01-01")
// 	require.NoError(t, err, "Parse should not return an error")

// 	createdSite, err := orm.CreateEmptySite(u)
// 	require.NoError(t, err, "CreateEmptySite should not return an error")

// 	assert.NotZero(t, createdSite.ID(), "Returned OrmSite should have a non-zero ID")
// 	assert.Equal(t, u.String(), createdSite.Url.String(), "Returned OrmSite should have the correct URL")
// 	assert.Equal(t, []string{}, createdSite.Validator.AllowedStrings(), "Retrieved site should have the correct allowed strings")
// 	assert.Equal(t, []string{}, createdSite.Validator.DisallowedStrings(), "Retrieved site should have the correct disallowed strings")
// 	assert.WithinDuration(t, defaultTime, createdSite.LastRobotsCrawl, time.Second, "Timestamps should be reasonably close")
// }

// func TestSiteByUrlOrCreateEmpty(t *testing.T) {
// 	t.Run("Exists", func(t *testing.T) {
// 		u := uniqueURL(t)
// 		v, err := robots.FromStrings([]string{"test"}, []string{"testing"})
// 		require.NoError(t, err, "FromStrings should not return an error")

// 		now := time.Now()

// 		createdSite, err := orm.CreateSite(u, *v, now)
// 		require.NoError(t, err, "CreateSite should not return an error")

// 		retrievedSite, err := orm.SiteByUrlOrCreateEmpty(u)
// 		require.NoError(t, err, "SiteByUrlOrCreateEmpty should not return an error")

// 		compareCreatedAndRetrievedSites(t, &createdSite, &retrievedSite)
// 	})

// 	t.Run("NotExists", func(t *testing.T) {
// 		u := uniqueURL(t)

// 		defaultTime, err := time.Parse("2006-01-02", "0000-01-01")
// 		require.NoError(t, err, "Parse should not return an error")

// 		createdSite, err := orm.SiteByUrlOrCreateEmpty(u)
// 		require.NoError(t, err, "SiteByUrlOrCreateEmpty should not return an error")

// 		assert.NotZero(t, createdSite.ID(), "Returned OrmSite should have a non-zero ID")
// 		assert.Equal(t, u.String(), createdSite.Url.String(), "Returned OrmSite should have the correct URL")
// 		assert.Equal(t, []string{}, createdSite.Validator.AllowedStrings(), "Retrieved site should have the correct allowed strings")
// 		assert.Equal(t, []string{}, createdSite.Validator.DisallowedStrings(), "Retrieved site should have the correct disallowed strings")
// 		assert.WithinDuration(t, defaultTime, createdSite.LastRobotsCrawl, time.Second, "Timestamps should be reasonably close")
// 	})
// }
