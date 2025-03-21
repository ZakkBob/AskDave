package orm

import (
	"context"
	"fmt"
	"time"

	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/jackc/pgx/v5"
)

func NextTasks(n int) (*tasks.Tasks, error) {
	query := `WITH ordered_crawls AS (
		SELECT page.id, page.url, page.next_crawl, robots.allowed_patterns, robots.disallowed_patterns,
		rank() OVER (PARTITION BY site.id ORDER BY page.next_crawl ASC, page.id DESC) as crawl_rank,
		(robots.last_crawl < CURRENT_DATE OR robots.last_crawl IS NULL) AS recrawl_robots
		FROM page 
		JOIN site on page.site_id = site.id 
		JOIN robots on robots.site_id = site.id  
		WHERE page.next_crawl <= CURRENT_DATE AND page.assigned IS FALSE 
		ORDER BY page.next_crawl ASC, crawl_rank ASC 
		LIMIT $1
	)
	UPDATE page
	SET assigned = TRUE
	FROM ordered_crawls
	WHERE page.id = ordered_crawls.id 
	RETURNING ordered_crawls.url, ordered_crawls.next_crawl, ordered_crawls.recrawl_robots, ordered_crawls.allowed_patterns, ordered_crawls.disallowed_patterns;`

	var t tasks.Tasks

	rows, err := dbpool.Query(context.Background(), query, n)
	if err != nil {
		return nil, fmt.Errorf("unable to get next %d tasks: %v", n, err)
	}
	defer rows.Close()

	var urlS string
	var allowed_patterns []string
	var disallowed_patterns []string
	var next_crawl time.Time
	var recrawl_robots bool

	_, err = pgx.ForEachRow(rows, []any{&urlS, &next_crawl, &recrawl_robots, &allowed_patterns, &disallowed_patterns}, func() error {
		u, err := url.ParseAbs(urlS)

		if err != nil {
			return fmt.Errorf("looping rows: %v", err)
		}

		if recrawl_robots {
			t.Robots.Slice = append(t.Robots.Slice, u)
		}

		t.Pages.Slice = append(t.Pages.Slice, u)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to get next %d tasks: %v", n, err)
	}

	return &t, nil
}
