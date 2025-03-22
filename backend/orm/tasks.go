package orm

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/jackc/pgx/v5"
)

func NextTasks(n int) (*tasks.Tasks, error) {
	query := `WITH ordered_crawls AS (
		SELECT page.id, page.site, page.path, page.next_crawl, robots.allowed_patterns, robots.disallowed_patterns,
		rank() OVER (PARTITION BY site.id ORDER BY page.next_crawl ASC, page.id DESC) as crawl_rank,
		(robots.last_crawl < CURRENT_DATE OR robots.last_crawl IS NULL) AS recrawl_robots
		FROM page 
		JOIN site on page.site = site.id 
		JOIN robots on robots.site = site.id  
		WHERE page.next_crawl <= CURRENT_DATE AND page.assigned IS FALSE 
		ORDER BY page.next_crawl ASC, crawl_rank ASC 
		LIMIT $1
	)
	UPDATE page
	SET assigned = FALSE
	FROM ordered_crawls
	WHERE page.id = ordered_crawls.id 
	RETURNING ordered_crawls.site, ordered_crawls.path, ordered_crawls.next_crawl, ordered_crawls.recrawl_robots, ordered_crawls.allowed_patterns, ordered_crawls.disallowed_patterns;`

	var t tasks.Tasks

	rows, err := dbpool.Query(context.Background(), query, n)
	if err != nil {
		return nil, fmt.Errorf("unable to get next %d tasks: %w", n, err)
	}
	defer rows.Close()

	var siteID int
	var sitePath string
	var allowed_patterns []string
	var disallowed_patterns []string
	var next_crawl time.Time
	var recrawl_robots bool

	var robotsUrls []string

	_, err = pgx.ForEachRow(rows, []any{&siteID, &sitePath, &next_crawl, &recrawl_robots, &allowed_patterns, &disallowed_patterns}, func() error {
		s, err := SiteByID(siteID)

		if err != nil {
			return fmt.Errorf("looping rows: %w", err)
		}

		u, err := url.ParseAbs(s.Url.StringNoPath() + sitePath)

		if err != nil {
			return fmt.Errorf("looping rows: %w", err)
		}

		t.Pages.Slice = append(t.Pages.Slice, u)

		u.Path = []string{}

		if recrawl_robots && !slices.Contains(robotsUrls, u.String()) {
			robotsUrls = append(robotsUrls, u.String())
			t.Robots.Slice = append(t.Robots.Slice, u)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to get next %d tasks: %w", n, err)
	}

	return &t, nil
}
