package orm

import (
	"context"
	"time"

	"errors"
	"fmt"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

type OrmPage struct {
	page.Page

	id            int
	NextCrawl     time.Time
	CrawlInterval int
	IntervalDelta int
	Assigned      bool
}

func (o *OrmPage) SaveCrawl(datetime time.Time, success bool, failureReason tasks.FailureReason, contentChanged bool, hash hash.Hash) error {
	query := `INSERT INTO crawl (page, datetime, success, failure_reason, content_changed, hash
		VALUES ($1, $2, $3, $4, $5, %6);`

	_, err := dbpool.Exec(context.Background(), query, o.id, datetime, success, failureReason, contentChanged, hash.String())
	if err != nil {
		return fmt.Errorf("unable to save crawl '%v' '%v' '%v' '%v' '%v': %w", datetime, success, failureReason, contentChanged, hash, err)
	}

	return nil
}

func SaveNewPage(p page.Page) (OrmPage, error) {
	query := `INSERT INTO page (site, path, title, og_title, og_description, og_sitename, next_crawl, crawl_interval, interval_delta, assigned, hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id;`

	nextCrawl := time.Now().AddDate(0, 0, 7)

	o := OrmPage{
		Page:          p,
		NextCrawl:     nextCrawl,
		CrawlInterval: 7,
		IntervalDelta: 1,
		Assigned:      false,
	}

	var s OrmSite
	var err error

	s, err = SaveNewSite(p.Url.StringNoPath())
	if err != nil {
		return o, fmt.Errorf("unable to save new page '%v': %w", p, err)
	}

	row := dbpool.QueryRow(context.Background(), query, s.id, o.Url.PathString(), o.Title, o.OgTitle, o.OgDescription, o.OgSiteName, o.NextCrawl, o.CrawlInterval, o.IntervalDelta, o.Assigned, o.Hash.String())

	err = row.Scan(&o.id)

	if err != nil {
		return o, fmt.Errorf("unable to save new page '%v': %w", p, err)
	}

	err = o.updateLinks()

	if err != nil {
		return o, fmt.Errorf("unable to save new page '%v': %w", o, err)
	}

	return o, nil
}

func (o *OrmPage) updateLinks() error {
	DeleteLinksBySrc(o.Url.String())
	log.Println("deleted all links")

	var p OrmPage
	var err error

	for _, dst := range o.Links { // Could be optimised if removing the orm
		p, err = PageByUrl(dst.String(), true)
		if errors.Is(err, pgx.ErrNoRows) {
			p, err = SaveNewPage(page.Page{
				Url:           dst,
				Title:         "",
				OgTitle:       "",
				OgDescription: "",
				OgSiteName:    "",
				Links:         []url.Url{},
				Hash:          hash.Hashs(""),
			})
			if err != nil {
				return fmt.Errorf("unable to update links (saving page) '%v': %w", o, err)
			}
		}
		if err != nil {
			return fmt.Errorf("unable to update links '%v': %w", o, err)
		}

		SaveNewLink(*o, p)
		log.Println("saved new link from " + o.Url.String() + "  " + p.Url.String())
	}
	return nil
}

func (o *OrmPage) Save(updateLinks bool) error {
	s, err := SaveNewSite(o.Url.FQDN())
	if err != nil {
		return fmt.Errorf("unable to save page '%v': %w", o, err)
	}

	query := `UPDATE page
		SET site = $2, path = $3, title = $4, og_title = $5, og_description = $6, og_sitename = $7, next_crawl = $8, crawl_interval = $9, interval_delta = $10, assigned = $11, hash = $12
		WHERE page.id = $1;`

	_, err = dbpool.Exec(context.Background(), query, o.id, s.id, o.Url.PathString(), o.Title, o.OgTitle, o.OgDescription, o.OgSiteName, o.NextCrawl, o.CrawlInterval, o.IntervalDelta, o.Assigned, o.Hash.String())
	if err != nil {
		return fmt.Errorf("unable to save page '%v': %w", o, err)
	}

	if updateLinks {
		err = o.updateLinks()
		if err != nil {
			return fmt.Errorf("unable to save page links '%v': %w", o, err)
		}
	}

	return nil
}

func pageFromRow(row pgx.Row, loadLinks bool) (OrmPage, error) {
	var p OrmPage
	var siteId int
	var path string
	var hashS string

	err := row.Scan(&p.id, &siteId, &path, &p.Title, &p.OgTitle, &p.OgDescription, &p.OgSiteName, &p.NextCrawl, &p.CrawlInterval, &p.IntervalDelta, &p.Assigned, &hashS)

	if err != nil {
		return p, err
	}

	p.Hash, err = hash.StrToHash(hashS)

	if err != nil {
		return p, err
	}

	// Get Url
	site, err := SiteByID(siteId)

	if err != nil {
		return p, err
	}

	u, err := url.ParseAbs(site.Url.String() + path)

	if err != nil {
		return p, err
	}

	p.Url = u

	if loadLinks {
		// Get Links
		dsts, err := LinkDstsBySrc(p.Url.String())
		if err != nil {
			return p, err
		}
		p.Links = dsts
	}
	return p, nil
}

func PageByID(id int, loadLinks bool) (OrmPage, error) {
	query := `SELECT id, site, path, title, og_title, og_description, og_sitename, next_crawl, crawl_interval, interval_delta, assigned, hash
		FROM page
		WHERE page.id = $1;`

	row := dbpool.QueryRow(context.Background(), query, id)
	p, err := pageFromRow(row, loadLinks)

	if err != nil {
		return p, fmt.Errorf("unable to get page from database for id '%d': %w", id, err)
	}

	return p, nil

}

func PageByUrl(urlS string, loadLinks bool) (OrmPage, error) {
	query := `SELECT id, site, path, title, og_title, og_description, og_sitename, next_crawl, crawl_interval, interval_delta, assigned, hash
		FROM page
		WHERE page.site = $1 AND page.path = $2;`

	u, err := url.ParseAbs(urlS)
	if err != nil {
		return OrmPage{}, fmt.Errorf("unable to get page from database for url '%s': %w", urlS, err)
	}

	s, err := SiteByUrl(u.StringNoPath())
	if err != nil {
		return OrmPage{}, fmt.Errorf("unable to get page from database for url '%s': %w", urlS, err)
	}

	row := dbpool.QueryRow(context.Background(), query, s.id, u.PathString())
	p, err := pageFromRow(row, loadLinks)

	if err != nil {
		return p, fmt.Errorf("unable to get page from database for url '%s': %w", urlS, err)
	}

	return p, nil

}
