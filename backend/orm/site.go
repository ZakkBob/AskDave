package orm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ZakkBob/AskDave/gocommon/robots"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/jackc/pgx/v5"
)

type OrmSite struct {
	id              int
	Url             url.Url
	Validator       robots.UrlValidator
	LastRobotsCrawl time.Time
}

func (s *OrmSite) ID() int {
	return s.id
}

func (o *OrmSite) Save() error {
	query := `UPDATE site
				SET url = $2, allowed_patterns = $3, disallowed_patterns = $4, last_robots_crawl = $5
				WHERE id = $1;`

	_, err := dbpool.Exec(context.Background(), query, o.id, o.Url.String(), o.Validator.AllowedStrings(), o.Validator.DisallowedStrings(), o.LastRobotsCrawl)
	if err != nil {
		return fmt.Errorf("failed to execute query '%s': %w", query, err)
	}

	return nil
}

func CreateSite(u url.Url, validator robots.UrlValidator, lastRobotsCrawl time.Time) (OrmSite, error) {
	var s OrmSite
	s.Url = u
	s.Validator = validator
	s.LastRobotsCrawl = lastRobotsCrawl

	query := `INSERT INTO site (url, allowed_patterns, disallowed_patterns, last_robots_crawl)
				VALUES ($1, $2, $3, $4)
				RETURNING id;`

	row := dbpool.QueryRow(context.Background(), query, u.String(), validator.AllowedStrings(), validator.DisallowedStrings(), lastRobotsCrawl)
	err := row.Scan(&s.id)
	if err != nil {
		return s, fmt.Errorf("failed to scan query results': %w", err)
	}

	return s, nil
}

func CreateEmptySite(u url.Url) (OrmSite, error) {
	v, err := robots.FromStrings([]string{}, []string{})
	if err != nil {
		return OrmSite{}, fmt.Errorf("failed to contruct empty UrlValidator: %w", err)
	}

	t, err := time.Parse("2006-01-02", "0000-01-01")
	if err != nil {
		return OrmSite{}, fmt.Errorf("failed to contruct empty time: %w", err)
	}

	s, err := CreateSite(u, *v, t)
	if err != nil {
		return s, fmt.Errorf("failed to create site: %w", err)
	}
	return s, nil
}

// eg. SELECT * FROM site WHERE id=1;
func SiteFromQuery(query string, args ...interface{}) (OrmSite, error) {
	var urlS string
	var allowedStrings []string
	var disallowedStrings []string
	var s OrmSite

	row := dbpool.QueryRow(context.Background(), query, args...)
	err := row.Scan(&s.id, &urlS, &allowedStrings, &disallowedStrings, &s.LastRobotsCrawl)
	if err != nil {
		return s, fmt.Errorf("failed to scan query results': %w", err)
	}

	v, err := robots.FromStrings(allowedStrings, disallowedStrings)
	if err != nil {
		return s, fmt.Errorf("failed to construct UrlValidator: %w", err)
	}

	s.Validator = *v

	u, err := url.ParseAbs(urlS)
	if err != nil {
		return s, fmt.Errorf("failed to parse url '%s': %w", urlS, err)
	}

	s.Url = u
	return s, nil
}

// Queries for site by url
// creates new empty site if not found
func SiteByUrlOrCreateEmpty(u url.Url) (OrmSite, error) {
	s, err := SiteByUrl(u)

	if err == nil {
		return s, nil
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return s, fmt.Errorf("failed to get site with url '%s': %w", u.String(), err)
	}

	s, err = CreateEmptySite(u)
	if err != nil {
		return s, fmt.Errorf("failed to create empty site: %w", err)
	}
	return s, nil
}

func SiteByID(id int) (OrmSite, error) {
	query := `SELECT *
				FROM site
				WHERE id = $1;`

	s, err := SiteFromQuery(query, id)
	if err != nil {
		return s, fmt.Errorf("failed to get site with query '%s': %w", query, err)
	}

	return s, nil
}

func SiteByUrl(u url.Url) (OrmSite, error) {
	query := `SELECT *
				FROM site
				WHERE url = $1;`

	s, err := SiteFromQuery(query, u.String())
	if err != nil {
		return s, fmt.Errorf("failed to get site with query '%s': %w", query, err)
	}

	return s, nil
}
