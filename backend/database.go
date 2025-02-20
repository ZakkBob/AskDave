package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/jackc/pgx/v5"
)

func saveResults(r *tasks.Results) error {
	robotsQuery := `UPDATE robots 
		SET (allowed_patterns, disallowed_patterns, last_crawl) = ($1, $2, CURRENT_DATE) 
		FROM site WHERE site.url = $3 AND site.id = robots.site_id;`

	for urlS, robotsResult := range r.Robots {
		if !robotsResult.Changed {
			continue
		}

		if !robotsResult.Success {
			continue // temporary - should probably do something
		}

		_, err := dbpool.Exec(context.Background(), robotsQuery, robotsResult.Validator.AllowedPatterns, robotsResult.Validator.DisallowedPatterns, urlS)
		if err != nil {
			return fmt.Errorf("unable to save robots result: %v", err)
		}
	}

	for urlS, pageResult := range r.Pages {
		if !pageResult.Changed {
			continue
		}

		// Add crawl to db
		if pageResult.Success {
			if pageResult.Changed {
				query := `INSERT INTO crawl (page_id, datetime, success, content_changed, title, og_title, og_description, hash)
					SELECT page.id, CURRENT_TIMESTAMP, TRUE, TRUE, $1, $2, $3, $4 FROM page WHERE url = $5;`

				_, err := dbpool.Exec(context.Background(), query, pageResult.Page.Title, pageResult.Page.OgTitle, pageResult.Page.OgDescription, pageResult.Page.Hash)
				if err != nil {
					return fmt.Errorf("unable to save page result: %v", err)
				}
			} else {
				query := `INSERT INTO crawl (page_id, datetime, success, content_changed)
					SELECT page.id, CURRENT_TIMESTAMP, TRUE, FALSE FROM page WHERE url = $1;`

				_, err := dbpool.Exec(context.Background(), query, urlS)
				if err != nil {
					return fmt.Errorf("unable to save page result: %v", err)
				}
			}
		} else {
			query := `INSERT INTO crawl (page_id, datetime, success)
				SELECT page.id, CURRENT_TIMESTAMP, FALSE FROM page WHERE url = $1;`

			_, err := dbpool.Exec(context.Background(), query, urlS)

			if err != nil {
				return fmt.Errorf("unable to save page result: %v", err)
			}
		}

		// Set next crawl date
		const maxCrawlInterval = 30
		const minCrawlInterval = 1

		var query string

		if pageResult.Changed {
			query = `UPDATE page SET 
				interval_delta = LEAST(interval_delta - 1, -1), 
				crawl_interval = LEAST(GREATEST(crawl_interval + interval_delta, $1),$2), 
				next_crawl = next_crawl + interval '1' day * crawl_interval 
				WHERE url = $3;`
		} else {
			query = `UPDATE page SET 
				interval_delta = GREATEST(interval_delta + 1, 1), 
				crawl_interval = LEAST(GREATEST(crawl_interval + interval_delta, $1),$2), 
				next_crawl = next_crawl + interval '1' day * crawl_interval 
				WHERE url = $3;`
		}

		_, err := dbpool.Exec(context.Background(), query, maxCrawlInterval, minCrawlInterval, urlS)
		if err != nil {
			return fmt.Errorf("unable to update page crawl scheduling: %v", err)
		}
	}

	return nil
}

func nextTasks(n int) (*tasks.Tasks, error) {
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
		return nil, fmt.Errorf("Unable to get next %d tasks: %v", n, err)
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
		return nil, fmt.Errorf("Unable to get next %d tasks: %v", n, err)
	}

	return &t, nil
}
