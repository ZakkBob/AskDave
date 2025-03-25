package orm

import (
	"context"
	"errors"
	"time"

	"fmt"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/url"

	"github.com/jackc/pgx/v5"
)

type OrmPage struct {
	page.Page

	id            int
	siteId        int
	NextCrawl     time.Time
	CrawlInterval int
	IntervalDelta int
	Assigned      bool
}

func (p *OrmPage) ID() int {
	return p.id
}

func (p *OrmPage) SiteID() int {
	return p.siteId
}

func (p *OrmPage) Site() (OrmSite, error) {
	s, err := SiteByID(p.id)
	if err != nil {
		return s, fmt.Errorf("failed to get page's site: %w", err)
	}
	return s, err
}

// func (p *OrmPage) SaveCrawl(datetime time.Time, success bool, failureReason tasks.FailureReason, contentChanged bool, hash hash.Hash) error {
// 	query := `INSERT INTO crawl (page, datetime, success, failure_reason, content_changed, hash
// 		VALUES ($1, $2, $3, $4, $5, $6);`

// 	_, err := dbpool.Exec(context.Background(), query, p.id, datetime, success, failureReason, contentChanged, hash.String())
// 	if err != nil {
// 		return fmt.Errorf("unable to save crawl '%v' '%v' '%v' '%v' '%v': %w", datetime, success, failureReason, contentChanged, hash, err)
// 	}

// 	return nil
// }

func (p *OrmPage) deleteLinksFrom() error {
	query := `DELETE
				FROM link
				WHERE src = $1;`

	_, err := dbpool.Exec(context.Background(), query, p.id)
	if err != nil {
		return fmt.Errorf("failed to execute query '%s' with arg '%d'", query, p.id)
	}
	return nil
}

// inefficient probably
func (p *OrmPage) saveLinks() error {
	err := p.deleteLinksFrom()
	if err != nil {
		return fmt.Errorf("failed to delete links from page: %w", err)
	}

	query := `SELECT page.id
                FROM page
                LEFT JOIN site ON site.id = page.site
                WHERE site.url = $1 AND page.path = $2;`

	var dstID int
	var dstIDs []int

	for _, dst := range p.Links {
		row := dbpool.QueryRow(context.Background(), query, dst.StringNoPath(), dst.PathString())
		err := row.Scan(&dstID)

		if errors.Is(err, pgx.ErrNoRows) {
			ormPage, err := CreateEmptyPage(dst)
			if err != nil {
				return fmt.Errorf("failed to create empty page with url '%s': %w", dst.String(), err)
			}

			dstID = ormPage.id
		} else if err != nil {
			return fmt.Errorf("failed to get page id with url '%s': %w", dst.String(), err)
		}

		dstIDs = append(dstIDs, dstID)
	}

	_, err = dbpool.CopyFrom(
		context.Background(),
		pgx.Identifier{"link"},
		[]string{"src", "dst"},
		pgx.CopyFromSlice(len(dstIDs), func(i int) ([]any, error) {
			return []any{p.id, dstIDs[i]}, nil
		}),
	)
	if err != nil {
		return fmt.Errorf("failed to copy links to database: %w", err)
	}
	return nil
}

func (p *OrmPage) loadLinks() error {
	query := `SELECT concat(site.url, page.path) AS dst 
		FROM link 
		LEFT JOIN page 
		ON page.id = link.dst 
		LEFT JOIN site 
		ON site.id = page.site 
		WHERE src = $1;`

	rows, err := dbpool.Query(context.Background(), query, p.id)
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var dstS string
	var dst url.Url

	for rows.Next() {
		err := rows.Scan(&dstS)
		if err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		dst, err = url.ParseAbs(dstS)
		if err != nil {
			return fmt.Errorf("failed to parse dst: %w", err)
		}
		p.Links = append(p.Links, dst)
	}
	return nil
}

func CreatePage(p page.Page, nextCrawl time.Time, crawlInterval int, intervalDelta int, assigned bool) (OrmPage, error) {
	query := `INSERT INTO page (site, path, title, og_title, og_description, og_sitename, hash, next_crawl, crawl_interval, interval_delta, assigned)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, next_crawl;`

	s, err := SiteByUrlOrCreateEmpty(p.Url)
	if err != nil {
		return OrmPage{}, fmt.Errorf("failed to get site by url '%s': %w", p.Url.String(), err)
	}

	ormPage := OrmPage{
		Page:          p,
		siteId:        s.id,
		NextCrawl:     nextCrawl,
		CrawlInterval: crawlInterval,
		IntervalDelta: intervalDelta,
		Assigned:      assigned,
	}

	row := dbpool.QueryRow(context.Background(), query, ormPage.siteId, ormPage.Url.PathString(), ormPage.Title, ormPage.OgTitle, ormPage.OgDescription, ormPage.OgSiteName, ormPage.Hash.String(), ormPage.NextCrawl, ormPage.CrawlInterval, ormPage.IntervalDelta, ormPage.Assigned)
	err = row.Scan(&ormPage.id, &ormPage.NextCrawl)
	if err != nil {
		return ormPage, fmt.Errorf("failed to scan query results: %w", err)
	}

	err = ormPage.saveLinks()
	if err != nil {
		return ormPage, fmt.Errorf("could not save page links: %w", err)
	}

	return ormPage, nil
}

func CreateEmptyPage(u url.Url) (OrmPage, error) {
	p := page.Page{
		Url:           u,
		Title:         "",
		OgTitle:       "",
		OgDescription: "",
		OgSiteName:    "",
		Links:         []url.Url{},
		Hash:          hash.Hashs(""),
	}

	ormPage, err := CreatePage(p, time.Now(), 7, 0, false)
	if err != nil {
		return ormPage, fmt.Errorf("failed to create page: %w", err)
	}

	return ormPage, nil
}

func (p *OrmPage) Save() error {
	query := `UPDATE page
		SET site = @site, path = @path, title = @title, og_title = @og_title, og_description = @og_description, og_sitename = @og_sitename, hash = @hash, next_crawl = @next_crawl, crawl_interval = @crawl_interval, interval_delta = @interval_delta, assigned = @assigned
		WHERE page.id = @id;`

	ormSite, err := SiteByUrlOrCreateEmpty(p.Url)
	if err != nil {
		return fmt.Errorf("failed to get site with url '%s' or create empty: %w", p.Url.String(), err)
	}

	_, err = dbpool.Exec(context.Background(),
		query,
		pgx.NamedArgs{
			"id":             p.id,
			"site":           ormSite.id,
			"path":           p.Url.PathString(),
			"title":          p.Title,
			"og_title":       p.OgTitle,
			"og_description": p.OgDescription,
			"og_sitename":    p.OgSiteName,
			"hash":           p.Hash.String(),
			"next_crawl":     p.NextCrawl,
			"crawl_interval": p.CrawlInterval,
			"interval_delta": p.IntervalDelta,
			"assigned":       p.Assigned,
		})
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	err = p.saveLinks()
	if err != nil {
		return fmt.Errorf("unable to save links: %w", err)
	}

	return nil
}

func PageFromQuery(query string, args ...interface{}) (OrmPage, error) {
	var ormPage OrmPage
	var path string
	var hashS string

	row := dbpool.QueryRow(context.Background(), query, args...)
	err := row.Scan(&ormPage.id, &ormPage.siteId, &path, &ormPage.Title, &ormPage.OgTitle, &ormPage.OgDescription, &ormPage.OgSiteName, &hashS, &ormPage.NextCrawl, &ormPage.CrawlInterval, &ormPage.IntervalDelta, &ormPage.Assigned)
	if err != nil {
		return ormPage, fmt.Errorf("failed to scan results from query '%s': %w", query, err)
	}

	ormPage.Hash, err = hash.StrToHash(hashS)
	if err != nil {
		return ormPage, fmt.Errorf("failed to set page hash to '%s': %w", hashS, err)
	}

	ormSite, err := SiteByID(ormPage.siteId)
	if err != nil {
		return ormPage, fmt.Errorf("failed to get site with id '%d': %w", ormPage.siteId, err)
	}

	urlS := ormSite.Url.StringNoPath() + path
	u, err := url.ParseAbs(urlS)
	if err != nil {
		return ormPage, fmt.Errorf("failed to parse url '%s'", urlS)
	}

	ormPage.Url = u

	ormPage.loadLinks()

	return ormPage, nil
}

func PageByID(id int) (OrmPage, error) {
	query := `SELECT *
		FROM page
		WHERE page.id = $1;`

	p, err := PageFromQuery(query, id)
	if err != nil {
		return p, fmt.Errorf("failed to get page from query: %w", err)
	}

	return p, nil
}

func PageByUrl(u url.Url) (OrmPage, error) {
	query := `SELECT page.*
                FROM page
                LEFT JOIN site ON site.id = page.site
                WHERE site.url = $1 AND page.path = $2;`

	p, err := PageFromQuery(query, u.StringNoPath(), u.PathString())
	if err != nil {
		return p, fmt.Errorf("unable to get page from query: %w", err)
	}

	return p, nil
}

func PageByUrlOrCreateEmpty(u url.Url) (OrmPage, error) {
	s, err := PageByUrl(u)

	if err == nil {
		return s, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return s, fmt.Errorf("failed to get page with url '%s': %w", u.String(), err)
	}

	s, err = CreateEmptyPage(u)
	if err != nil {
		return s, fmt.Errorf("failed to create empty page: %w", err)
	}
	return s, nil
}
