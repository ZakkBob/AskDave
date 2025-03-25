package orm

import (
	"context"
	"fmt"

	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/jackc/pgx/v5"
)

func NextTasks(n int) (*tasks.Tasks, error) {
	query := `WITH ordered_crawls AS (
		SELECT page.id AS page_id, site.id AS site_id, CONCAT(site.url, page.path) AS url,
		rank() OVER (PARTITION BY site.id ORDER BY page.next_crawl ASC, page.id DESC) as crawl_rank,
		(site.last_robots_crawl < CURRENT_DATE OR site.last_robots_crawl IS NULL) AS recrawl_robots
		FROM page 
		JOIN site on page.site = site.id   
		WHERE page.next_crawl <= CURRENT_DATE AND page.assigned IS FALSE 
		ORDER BY page.next_crawl ASC, crawl_rank ASC 
		LIMIT $1
	)
	UPDATE page
	SET assigned = TRUE
	FROM ordered_crawls
	WHERE page.id = ordered_crawls.page_id 
	RETURNING ordered_crawls.url, ordered_crawls.recrawl_robots;`

	var t tasks.Tasks

	rows, err := dbpool.Query(context.Background(), query, n)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var urlS string
	var recrawl_robots bool

	_, err = pgx.ForEachRow(rows, []any{&urlS, &recrawl_robots}, func() error {
		u, err := url.ParseAbs(urlS)
		if err != nil {
			return fmt.Errorf("failed to parse url: %w", err)
		}

		if recrawl_robots {
			t.Robots.Slice = append(t.Robots.Slice, u)
		}

		t.Pages.Slice = append(t.Pages.Slice, u)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to process rows: %w", err)
	}

	return &t, nil
}
