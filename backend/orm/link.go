package orm

import (
	"context"

	"fmt"

	"github.com/ZakkBob/AskDave/gocommon/url"
)

type OrmLink struct {
	id  int
	Src url.Url
	Dst url.Url
}

func SaveNewLink(src OrmPage, dst OrmPage) (OrmLink, error) {
	var l OrmLink

	query := `INSERT INTO link (src, dst)
		VALUES ($1, $2)
		RETURNING id;`

	row := dbpool.QueryRow(context.Background(), query, src.id, dst.id)

	err := row.Scan(&l.id)
	if err != nil {
		return l, fmt.Errorf("unable to save new link with src '%s', dst '%s': %v", src.Url.String(), dst.Url.String(), err)
	}

	l.Src = src.Url
	l.Dst = dst.Url

	return l, nil
}

func (l *OrmLink) Save() error {
	query := `UPDATE link
		SET src = $2, dst = $3
		WHERE link.id = $1;`

	srcPage, err := PageByUrl(l.Src.String())

	if err != nil {
		return fmt.Errorf("unable to save link '%v': %v", l, err)
	}

	dstPage, err := PageByUrl(l.Src.String())

	if err != nil {
		return fmt.Errorf("unable to save link '%v': %v", l, err)
	}

	_, err = dbpool.Exec(context.Background(), query, l.id, srcPage.id, dstPage.id)
	if err != nil {
		return fmt.Errorf("unable to save link '%v': %v", l, err)
	}

	return nil
}

type scanInterface interface {
	Scan(...interface{}) error
}

func linkFromRow(row scanInterface) (OrmLink, error) {
	var l OrmLink
	var srcId int
	var dstId int

	err := row.Scan(&l.id, &srcId, &dstId)

	if err != nil {
		return OrmLink{}, fmt.Errorf("%v", err)
	}

	srcPage, err := PageByID(srcId)

	if err != nil {
		return OrmLink{}, fmt.Errorf("%v", err)
	}

	dstPage, err := PageByID(dstId)

	if err != nil {
		return OrmLink{}, fmt.Errorf("%v", err)
	}

	l.Src = srcPage.Url
	l.Dst = dstPage.Url

	return l, nil
}

func LinkByID(id int) (OrmLink, error) {
	var l OrmLink

	query := `SELECT id, src, dst
		FROM link
		WHERE link.id = $1;`

	row := dbpool.QueryRow(context.Background(), query, id)
	l, err := linkFromRow(row)
	if err != nil {
		return OrmLink{}, fmt.Errorf("unable to get link from database for id '%d': %v", id, err)
	}

	return l, nil
}

func LinksBySrc(src string) ([]OrmLink, error) {
	links := make([]OrmLink, 0)

	query := `SELECT id, src, dst
		FROM link
		WHERE link.src = $1;`

	rows, err := dbpool.Query(context.Background(), query, src)

	if err != nil {
		return []OrmLink{}, fmt.Errorf("unable to get links from database for src '%s': %v", src, err)
	}

	for rows.Next() {
		l, err := linkFromRow(rows)
		if err != nil {
			return []OrmLink{}, fmt.Errorf("unable to get links from database for src '%s': %v", src, err)
		}

		links = append(links, l)
	}

	return links, nil
}

func LinkDstsBySrc(src string) ([]url.Url, error) {
	var urlS string
	urls := make([]url.Url, 0)

	query := `SELECT dst
		FROM link
		WHERE link.src = $1;`

	rows, err := dbpool.Query(context.Background(), query, src)

	if err != nil {
		return []url.Url{}, fmt.Errorf("unable to get destinations from database for src '%s': %v", src, err)
	}

	for rows.Next() {
		err := rows.Scan(&urlS)

		if err != nil {
			return []url.Url{}, fmt.Errorf("unable to get destinations from database for src '%s': %v", src, err)
		}

		u, err := url.ParseAbs(urlS)

		if err != nil {
			return []url.Url{}, fmt.Errorf("unable to get destinations from database for src '%s': %v", src, err)
		}

		urls = append(urls, u)
	}

	return urls, nil
}

func DeleteLinksBySrc(src string) error {
	query := `DELETE
		FROM link
		WHERE link.src = $1;`

	_, err := dbpool.Exec(context.Background(), query)

	if err != nil {
		return fmt.Errorf("unable to get destinations from database for src '%s': %v", src, err)
	}

	return nil
}
