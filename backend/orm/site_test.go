package orm_test

import (
	"testing"
	"time"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/ZakkBob/AskDave/gocommon/robots"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func compareCreatedAndRetrievedSites(t *testing.T, createdSite *orm.OrmSite, retrievedSite *orm.OrmSite) {
	assert.Equal(t, createdSite.ID(), retrievedSite.ID(), "Retrieved site should have correct ID")
	assert.Equal(t, createdSite.Url.String(), retrievedSite.Url.String(), "Retrieved site should have correct URL")
	assert.Equal(t, createdSite.Validator.AllowedStrings(), retrievedSite.Validator.AllowedStrings(), "Retrieved site should have the correct allowed strings")
	assert.Equal(t, createdSite.Validator.DisallowedStrings(), retrievedSite.Validator.DisallowedStrings(), "Retrieved site should have the correct disallowed strings")
	assert.WithinDuration(t, createdSite.LastRobotsCrawl, retrievedSite.LastRobotsCrawl, time.Second, "Retrieved timestamps should be reasonably close")
}

func TestCreateSite_And_SiteByUrl_And_SiteByID(t *testing.T) {
	defer resetDB(t)

	siteTests := []struct {
		name string
		in   string
		out  string
	}{
		{"nopath", "https://testcreatesite_and_sitebyurl_and_sitebyid.com", "https://testcreatesite_and_sitebyurl_and_sitebyid.com"},
		{"withpath", "https://testcreatesite_and_sitebyurl_and_sitebyid2.com/path/1", "https://testcreatesite_and_sitebyurl_and_sitebyid2.com"},
	}

	for _, tt := range siteTests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.ParseAbs(tt.in)
			require.NoError(t, err, "ParseAbs should not return an error")

			v := robots.UrlValidator{}
			now := time.Now()

			createdSite, err := orm.CreateSite(u, v, now)
			require.NoError(t, err, "CreateSite should not return an error")

			assert.NotZero(t, createdSite.ID(), "Returned OrmSite should have a non-zero ID")
			assert.Equal(t, tt.out, createdSite.Url.String(), "Returned OrmSite should have the correct URL")
			assert.Equal(t, v.AllowedStrings(), createdSite.Validator.AllowedStrings(), "Retrieved site should have the correct allowed strings")
			assert.Equal(t, v.DisallowedStrings(), createdSite.Validator.DisallowedStrings(), "Retrieved site should have the correct disallowed strings")
			assert.WithinDuration(t, now, createdSite.LastRobotsCrawl, time.Second, "Timestamps should be reasonably close")

			t.Run("SiteByUrl", func(t *testing.T) {
				retrievedSite, err := orm.SiteByUrl(u)
				require.NoError(t, err, "SiteByUrl should not return an error")

				compareCreatedAndRetrievedSites(t, &createdSite, &retrievedSite)
			})

			t.Run("SiteByID", func(t *testing.T) {
				retrievedSite, err := orm.SiteByID(createdSite.ID())
				require.NoError(t, err, "SiteByID should not return an error")

				compareCreatedAndRetrievedSites(t, &createdSite, &retrievedSite)
			})
		})
	}
}

func TestSiteSave(t *testing.T) {
	defer resetDB(t)

	u := makeURL(t, "")
	v := robots.UrlValidator{}
	now := time.Now()

	createdSite, err := orm.CreateSite(u, v, now)
	require.NoError(t, err, "CreateSite should not return an error")

	updatedValidator, err := robots.FromStrings([]string{"/allowed"}, []string{"disallowed"})
	require.NoError(t, err, "FromStrings should not return an error")

	updatedUrl, err := url.ParseAbs("https://testsitesave-updatedurl.com")
	require.NoError(t, err, "ParseAbs should not return an error")

	updatedTime := time.Now().Add(time.Second * 2)

	createdSite.Validator = *updatedValidator
	createdSite.Url = updatedUrl
	createdSite.LastRobotsCrawl = updatedTime

	createdSite.Save()

	savedSite, err := orm.SiteByID(createdSite.ID())
	require.NoError(t, err, "SiteByID should not return an error")

	assert.Equal(t, createdSite.ID(), savedSite.ID(), "Retrieved site should have correct ID")
	assert.Equal(t, updatedUrl.String(), savedSite.Url.String(), "Retrieved site should have correct URL")
	assert.Equal(t, updatedValidator.AllowedStrings(), savedSite.Validator.AllowedStrings(), "Retrieved site should have the correct allowed strings")
	assert.Equal(t, updatedValidator.DisallowedStrings(), savedSite.Validator.DisallowedStrings(), "Retrieved site should have the correct disallowed strings")
	assert.WithinDuration(t, updatedTime, savedSite.LastRobotsCrawl, time.Second, "Retrieved timestamps should be reasonably close")
}

func TestCreateEmptySite(t *testing.T) {
	defer resetDB(t)

	u := makeURL(t, "")

	defaultTime, err := time.Parse("2006-01-02", "0000-01-01")
	require.NoError(t, err, "Parse should not return an error")

	createdSite, err := orm.CreateEmptySite(u)
	require.NoError(t, err, "CreateEmptySite should not return an error")

	assert.NotZero(t, createdSite.ID(), "Returned OrmSite should have a non-zero ID")
	assert.Equal(t, u.String(), createdSite.Url.String(), "Returned OrmSite should have the correct URL")
	assert.Equal(t, []string{}, createdSite.Validator.AllowedStrings(), "Retrieved site should have the correct allowed strings")
	assert.Equal(t, []string{}, createdSite.Validator.DisallowedStrings(), "Retrieved site should have the correct disallowed strings")
	assert.WithinDuration(t, defaultTime, createdSite.LastRobotsCrawl, time.Second, "Timestamps should be reasonably close")
}

func TestSiteByUrlOrCreateEmpty(t *testing.T) {
	defer resetDB(t)

	t.Run("Exists", func(t *testing.T) {
		u := makeURL(t, "")
		v, err := robots.FromStrings([]string{"test"}, []string{"testing"})
		require.NoError(t, err, "FromStrings should not return an error")

		now := time.Now()

		createdSite, err := orm.CreateSite(u, *v, now)
		require.NoError(t, err, "CreateSite should not return an error")

		retrievedSite, err := orm.SiteByUrlOrCreateEmpty(u)
		require.NoError(t, err, "SiteByUrlOrCreateEmpty should not return an error")

		compareCreatedAndRetrievedSites(t, &createdSite, &retrievedSite)
	})

	t.Run("NotExists", func(t *testing.T) {
		u := makeURL(t, "")

		defaultTime, err := time.Parse("2006-01-02", "0000-01-01")
		require.NoError(t, err, "Parse should not return an error")

		createdSite, err := orm.SiteByUrlOrCreateEmpty(u)
		require.NoError(t, err, "SiteByUrlOrCreateEmpty should not return an error")

		assert.NotZero(t, createdSite.ID(), "Returned OrmSite should have a non-zero ID")
		assert.Equal(t, u.String(), createdSite.Url.String(), "Returned OrmSite should have the correct URL")
		assert.Equal(t, []string{}, createdSite.Validator.AllowedStrings(), "Retrieved site should have the correct allowed strings")
		assert.Equal(t, []string{}, createdSite.Validator.DisallowedStrings(), "Retrieved site should have the correct disallowed strings")
		assert.WithinDuration(t, defaultTime, createdSite.LastRobotsCrawl, time.Second, "Timestamps should be reasonably close")
	})
}
