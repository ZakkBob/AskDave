package orm_test

import (
	"testing"
	"time"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/ZakkBob/AskDave/gocommon/robots"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/require"
)

func TestNextTasks(t *testing.T) {
	defer resetDB(t)

	type siteTest struct {
		url             string
		lastRobotsCrawl time.Time
	}

	type pageTest struct {
		url       string
		nextCrawl time.Time
	}

	taskTests := []struct {
		name      string
		taskCount int
		sites     []siteTest
		pages     []pageTest
	}{
		{
			name:      "",
			taskCount: 1,
			pages: []pageTest{
				{
					url:       "https://a.com",
					nextCrawl: time.Now(),
				},
			},
		},
	}

	for _, tt := range taskTests {
		t.Run(tt.name, func(t *testing.T) {
			for _, s := range tt.sites {
				u, err := url.ParseAbs(s.url)
				require.NoError(t, err, "ParseAbs should not return an error")
				orm.CreateSite(u, robots.UrlValidator{}, s.lastRobotsCrawl)

				t, err := orm.NextTasks()
			}
		})
	}

	_, err := orm.CreateEmptyPage(makeURL(t, ""))
	require.NoError(t, err, "CreateEmptySite should not return an error")

	_, err = orm.CreateEmptyPage(makeURL(t, "/a"))
	require.NoError(t, err, "CreateEmptySite should not return an error")

	tasks, err := orm.NextTasks(1)
	require.NoError(t, err, "NextTasks should not return an error")
	t.Log(tasks.Pages.Slice)
	t.Log(tasks.Robots.Slice)
}
