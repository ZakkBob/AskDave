package orm

import (
	"time"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/tasks"

	"context"
	"fmt"
)

func SaveResults(r *tasks.Results) error {
	robotsQuery := `UPDATE robots 
		SET (allowed_patterns, disallowed_patterns, last_crawl) = ($1, $2, CURRENT_DATE) 
		FROM site WHERE site.url = $3 
		AND site.id = robots.site;`

	for urlS, robotsResult := range r.Robots {
		if !robotsResult.Changed {
			continue
		}

		if !robotsResult.Success {
			continue // temporary - should probably do something
		}

		_, err := dbpool.Exec(context.Background(), robotsQuery, robotsResult.Validator.AllowedStrings(), robotsResult.Validator.DisallowedStrings(), urlS)
		if err != nil {
			return fmt.Errorf("unable to save robots result: %w", err)
		}
	}

	const maxCrawlInterval = 60
	const minCrawlInterval = 7

	for urlS, pageResult := range r.Pages {
		p, err := PageByUrl(urlS, true)
		if err != nil {
			return fmt.Errorf("unable to save page result: %w", err)
		}

		var h hash.Hash

		if pageResult.Changed {
			p.Title = pageResult.Page.Title
			p.OgTitle = pageResult.Page.OgTitle
			p.OgDescription = pageResult.Page.OgDescription
			p.Hash = pageResult.Page.Hash
			p.Links = pageResult.Page.Links

			p.IntervalDelta--
			if p.IntervalDelta > -1 {
				p.IntervalDelta = -1
			}
		} else {
			p.IntervalDelta++
			if p.IntervalDelta < 1 {
				p.IntervalDelta = 1
			}
		}

		p.CrawlInterval += p.IntervalDelta
		if p.CrawlInterval < minCrawlInterval {
			p.CrawlInterval = minCrawlInterval
		} else if p.CrawlInterval > maxCrawlInterval {
			p.CrawlInterval = maxCrawlInterval
		}
		p.NextCrawl = p.NextCrawl.AddDate(0, 0, p.CrawlInterval)

		err = p.Save(true)
		if err != nil {
			return fmt.Errorf("unable to save page result: %w", err)
		}

		if pageResult.Page == nil {
			h = hash.Hashs("")
		} else {
			h = pageResult.Page.Hash
		}

		err = p.SaveCrawl(time.Now(), pageResult.Success, pageResult.FailureReason, pageResult.Changed, h)
		if err != nil {
			return fmt.Errorf("unable to save page result: %w", err)
		}
	}

	return nil
}
