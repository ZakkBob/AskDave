package orm

import (
	"context"
	"date"
	"time"

	"fmt"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/jackc/pgx/v5"
)

type OrmPage struct {
	page.Page

	id            int
	NextCrawl     time.Time
	CrawlInterval int
	IntervalDelta int
	Assigned      bool
}

func (o *OrmPage) SaveCrawl(datetime date.Date, success bool, failureReason tasks.FailureReason, contentChanged bool, hash hash.Hash) error {
	query := `INSERT INTO crawl (page, datetime, success, failure_reason, content_changed, hash
		VALUES ($1, $2, $3, $4, $5, %6);`

	_, err := dbpool.Exec(context.Background(), query, o.id, datetime, success, failureReason, contentChanged, hash)
	if err != nil {
		return fmt.Errorf("unable to save crawl '%s' '%s' '%s' '%s' '%s' '%s': %v", datetime, success, failureReason, contentChanged, hash, err)
	}

	return nil
}

func SaveNewPage(p page.Page) (OrmPage, error) {
	query := `INSERT INTO page (site, path, title, og_title, og_description, og_sitename, next_crawl, crawl_interval, interval_delta, assigned)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id;`

	nextCrawl := time.Now().AddDate(0, 0, 7)

	o := OrmPage{
		Page:          p,
		NextCrawl:     nextCrawl,
		CrawlInterval: 7,
		IntervalDelta: 1,
		Assigned:      false,
	}

	s, err := SiteByUrl(p.Url.StringNoPath())
	if err != nil {
		return o, fmt.Errorf("unable to save new page '%v': %v", p, err)
	}

	row := dbpool.QueryRow(context.Background(), query, s.id, o.Url.PathString(), o.Title, o.OgTitle, o.OgDescription, o.OgSiteName, o.NextCrawl, o.CrawlInterval, o.IntervalDelta, o.Assigned)

	err = row.Scan(&o.id)

	if err != nil {
		return o, fmt.Errorf("unable to save new page '%v': %v", p, err)
	}

	err = o.updateLinks()

	if err != nil {
		return o, fmt.Errorf("unable to save new page '%v': %v", o, err)
	}

	return o, nil
}

func (o *OrmPage) updateLinks() error {
	DeleteLinksBySrc(o.Url.String())

	for _, dst := range o.Links { // Could be optimised if removing the orm
		p, err := PageByUrl(dst.String())
		if err == pgx.ErrNoRows {

		}
		if err != nil {
			return fmt.Errorf("unable to save link '%v': %v", o, err)
		}

		SaveNewLink(*o, p)
	}
	return nil
}

func (o *OrmPage) Save() error {
	s, err := SiteByUrl(o.Url.FQDN())
	if err != nil {
		return fmt.Errorf("unable to save page '%v': %v", o, err)
	}

	query := `UPDATE page
		SET site = $2, path = $3, title = $4, og_title = $5, og_description = $6, og_sitename = $7, next_crawl = $8, crawl_interval = $9, interval_delta = $10, assigned = $11
		WHERE link.id = $1;`

	_, err = dbpool.Exec(context.Background(), query, o.id, s.id, o.Url.PathString(), o.Title, o.OgDescription, o.OgSiteName, o.NextCrawl, o.CrawlInterval, o.IntervalDelta, o.Assigned)
	if err != nil {
		return fmt.Errorf("unable to save page '%v': %v", o, err)
	}

	err = o.updateLinks()
	if err != nil {
		return fmt.Errorf("unable to save page '%v': %v", o, err)
	}

	return nil
}

func pageFromRow(row pgx.Row) (OrmPage, error) {
	var p OrmPage
	var siteId int
	var path string

	err := row.Scan(p.id, siteId, path, p.Title, p.OgTitle, p.OgDescription, p.OgSiteName, p.NextCrawl, p.CrawlInterval, p.IntervalDelta, p.Assigned)

	if err != nil {
		return p, err
	}

	// Get Url
	site, err := SiteByID(siteId)

	if err != nil {
		return p, err
	}

	u, err := url.ParseRel(path, site.Url)

	if err != nil {
		return p, err
	}

	p.Url = u

	// Get Links
	dsts, err := LinkDstsBySrc(p.Url.String())
	if err != nil {
		return p, err
	}
	p.Links = dsts

	return p, nil
}

func PageByID(id int) (OrmPage, error) {
	query := `SELECT id, site, path, title, og_title, og_description, og_sitename, next_crawl, crawl_interval, interval_delta, assigned
		FROM page
		WHERE page.id = $1;`

	row := dbpool.QueryRow(context.Background(), query, id)
	p, err := pageFromRow(row)

	if err != nil {
		return p, fmt.Errorf("unable to get page from database for id '%d': %v", id, err)
	}

	return p, nil

}

func PageByUrl(urlS string) (OrmPage, error) {
	query := `SELECT id, site, path, title, og_title, og_description, og_sitename, next_crawl, crawl_interval, interval_delta, assigned
		FROM page
		WHERE page.url = $1;`

	row := dbpool.QueryRow(context.Background(), query, urlS)
	p, err := pageFromRow(row)

	if err != nil {
		return p, fmt.Errorf("unable to get page from database for url '%s': %v", urlS, err)
	}

	return p, nil

}
