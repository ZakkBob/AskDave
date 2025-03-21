package orm

import (
	"context"
	"fmt"

	"github.com/ZakkBob/AskDave/gocommon/url"
)

type OrmSite struct {
	id  int
	Url url.Url
}

func SaveNewSite(urlS string) (OrmSite, error) {
	o, err := SiteByUrl(urlS)

	if err == nil {
		return o, nil
	}

	u, err := url.ParseAbs(urlS)

	if err != nil {
		return OrmSite{}, fmt.Errorf("unable to save new site with url '%s': %v", urlS, err)
	}

	o = OrmSite{
		Url: u,
	}

	query := `INSERT INTO site (url)
		VALUES ($1)
		RETURNING id;`

	row := dbpool.QueryRow(context.Background(), query, urlS)

	err = row.Scan(&o.id)

	if err != nil {
		return o, fmt.Errorf("unable to save new site with url '%s': %v", urlS, err)
	}

	return o, nil
}

func (o *OrmSite) Save() error {
	query := `UPDATE site
		SET url = $2
		WHERE site.id = $1;`

	_, err := dbpool.Exec(context.Background(), query, o.id, o.Url)
	if err != nil {
		return fmt.Errorf("unable to save site with id '%d': %v", o.id, err)
	}

	return nil
}

func SiteByID(id int) (OrmSite, error) {
	var s OrmSite
	var urlS string

	query := `SELECT id, url
		FROM site
		WHERE site.id = $1;`

	row := dbpool.QueryRow(context.Background(), query, id)
	err := row.Scan(&s.id, &urlS)

	if err != nil {
		return s, fmt.Errorf("unable to get site from database for id '%d': %v", id, err)
	}

	u, err := url.ParseAbs(urlS)

	if err != nil {
		return OrmSite{}, fmt.Errorf("unable to get site from database for id '%d': %v", id, err)
	}

	s.Url = u
	return s, nil
}

func SiteByUrl(urlS string) (OrmSite, error) {
	var s OrmSite

	query := `SELECT id
		FROM site
		WHERE site.url = $1;`

	row := dbpool.QueryRow(context.Background(), query, urlS)
	err := row.Scan(&s.id)

	if err != nil {
		return s, fmt.Errorf("unable to get site from database for url '%s': %v", urlS, err)
	}

	u, err := url.ParseAbs(urlS)

	if err != nil {
		return OrmSite{}, fmt.Errorf("unable to get site from database for url '%s': %v", urlS, err)
	}

	s.Url = u
	return s, nil
}
