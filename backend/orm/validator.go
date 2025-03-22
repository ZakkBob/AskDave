package orm

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/ZakkBob/AskDave/gocommon/robots"
	"github.com/jackc/pgx/v5"
)

type OrmValidator struct {
	robots.UrlValidator

	id        int
	site      int
	LastCrawl time.Time
}

func (v *OrmValidator) Save() error {
	query := `UPDATE robots 
		SET (allowed_patterns, disallowed_patterns, last_crawl) = ($2, $3, $4) 
		WHERE robots.id = $1;`

	_, err := dbpool.Exec(context.Background(), query, v.id, v.AllowedStrings(), v.DisallowedStrings(), v.LastCrawl)
	if err != nil {
		return fmt.Errorf("unable to save robots result: %w", err)
	}
	return nil
}

func SaveNewValidator(v robots.UrlValidator, siteID int) (OrmValidator, error) {
	var ormV OrmValidator

	query := `INSERT INTO robots (allowed_patterns, disallowed_patterns, site)
		VALUES ($1, $2, $3)
		RETURNING id;`

	row := dbpool.QueryRow(context.Background(), query, v.AllowedStrings(), v.DisallowedStrings(), siteID)

	err := row.Scan(&ormV.id)
	if err != nil {
		return ormV, fmt.Errorf("unable to save new validator '%v': %w", v, err)
	}

	return ormV, nil
}

func validatorFromRow(row pgx.Row) (OrmValidator, error) {
	var v OrmValidator
	var allowedPatterns []string
	var disallowedPatterns []string

	err := row.Scan(&allowedPatterns, &disallowedPatterns, &v.site, &v.LastCrawl)
	if err != nil {
		return v, err
	}

	for _, allowed := range allowedPatterns {
		r, err := regexp.Compile(allowed)
		if err != nil {
			return v, err
		}
		v.AllowedPatterns = append(v.AllowedPatterns, r)
	}

	for _, disallowed := range disallowedPatterns {
		r, err := regexp.Compile(disallowed)
		if err != nil {
			return v, err
		}
		v.DisallowedPatterns = append(v.DisallowedPatterns, r)
	}

	return v, nil
}

func ValidatorByUrl(urlS string) (OrmValidator, error) {
	query := `SELECT allowed_patterns, disallowed_patterns, site, last_crawl 
		FROM robots 
		JOIN site 
		ON robots.site_id = site.id 
		AND site.url = $1;`

	row := dbpool.QueryRow(context.Background(), query, urlS)
	o, err := validatorFromRow(row)
	if err != nil {
		return o, fmt.Errorf("unable to get validator for url '%s': %w", urlS, err)
	}

	return o, nil
}

func ValidatorByID(id int) (OrmValidator, error) {
	query := `SELECT allowed_patterns, disallowed_patterns, site, last_crawl 
		FROM robots 
		WHERE robots.id = $1;`

	row := dbpool.QueryRow(context.Background(), query, id)
	o, err := validatorFromRow(row)
	if err != nil {
		return o, fmt.Errorf("unable to get validator for id '%d': %w", id, err)
	}

	return o, nil
}
