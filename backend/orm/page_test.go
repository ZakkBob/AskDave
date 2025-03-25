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

func assertPageEmpty(t *testing.T, u url.Url, p orm.OrmPage) {
	comparePageData(t, &page.Page{
		Url:           u,
		Title:         "",
		OgTitle:       "",
		OgDescription: "",
		OgSiteName:    "",
		Links:         []url.Url{},
		Hash:          hash.Hashs(""),
	}, &p.Page)

	now := time.Now()

	assert.WithinDuration(t, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.UTC().Location()), p.NextCrawl, time.Second, "OrmPage's next crawl should be reasonably close to expected")
	assert.Equal(t, 7, p.CrawlInterval, "OrmPage should have the default crawl interval")
	assert.Equal(t, 0, p.IntervalDelta, "OrmPage should have the default interval delta")
	assert.Equal(t, false, p.Assigned, "OrmPage should have the default assigned value")
}

func comparePageData(t *testing.T, expected *page.Page, actual *page.Page) {
	assert.Equal(t, expected.Url.String(), actual.Url.String(), "Page data should have correct URL")
	assert.Equal(t, expected.Title, actual.Title, "Page data should have correct Title")
	assert.Equal(t, expected.OgTitle, actual.OgTitle, "Page data should have correct OgTitle")
	assert.Equal(t, expected.OgDescription, actual.OgDescription, "Page data should have correct OgDescription")
	assert.Equal(t, expected.OgSiteName, actual.OgSiteName, "Page data should have correct OgSiteName")
	assert.Equal(t, expected.Hash.String(), actual.Hash.String(), "Page data should have correct Hash string")

	var expectedDstStrings []string
	var actualDstStrings []string

	for _, dst := range expected.Links {
		expectedDstStrings = append(expectedDstStrings, dst.String())
	}
	for _, dst := range actual.Links {
		actualDstStrings = append(actualDstStrings, dst.String())
	}

	assert.Equal(t, expectedDstStrings, actualDstStrings, "Page data should have correct links")
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

	dsts, err := url.ParseMany([]string{"https://www.testcreatepageandbyurlbyidlink1.com/path/e", "https://www.testcreatepageandbyurlbyidlink2.com"})
	assert.NoError(t, err, "ParseMany should not return an error")

	p := page.Page{
		Url:           u,
		Title:         "title",
		OgTitle:       "og title",
		OgDescription: "og description",
		OgSiteName:    "og site name",
		Hash:          hash.Hashs("hashed"),
		Links:         dsts,
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

func TestCreateEmptyPage(t *testing.T) {
	u := uniqueURL(t)

	createdPage, err := orm.CreateEmptyPage(u)
	require.NoError(t, err, "CreateEmptyPage should not return an error")

	assertPageEmpty(t, u, createdPage)
}

func TestPageByUrlOrCreateEmpty(t *testing.T) {
	t.Run("Exists", func(t *testing.T) {
		u := uniqueURL(t)

		createdPage, err := orm.CreateEmptyPage(u)
		require.NoError(t, err, "CreateEmptyPage should not return an error")

		retrievedPage, err := orm.PageByUrlOrCreateEmpty(u)
		require.NoError(t, err, "SiteByUrlOrCreateEmpty should not return an error")

		compareCreatedAndRetrievedPages(t, &createdPage, &retrievedPage)
	})

	t.Run("NotExists", func(t *testing.T) {
		u := uniqueURL(t)

		createdPage, err := orm.PageByUrlOrCreateEmpty(u)
		require.NoError(t, err, "SiteByUrlOrCreateEmpty should not return an error")

		assert.NotZero(t, createdPage.ID(), "Returned OrmPage should have a non-zero ID")
		assert.NotZero(t, createdPage.SiteID(), "Returned OrmPage should have a non-zero Site ID")

		assertPageEmpty(t, u, createdPage)
	})
}
