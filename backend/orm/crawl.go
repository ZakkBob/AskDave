package orm

import (
	"context"
	"time"

	"fmt"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
)

type OrmCrawl struct {
	id             int
	Url            url.Url
	Datetime       time.Time
	Success        bool
	FailureReason  tasks.FailureReason
	ContentChanged bool
	Hash           hash.Hash
}

func (c *OrmCrawl) Save() error {
	query := `UPDATE crawl
		SET page = $2, datetime = $3, success = $4, failure_reason = $5, content_changed = $6, hash = $7
		WHERE link.id = $1;`

	p, err := PageByUrl(c.Url.String())

	if err != nil {
		return fmt.Errorf("unable to save crawl '%v': %v", c, err)
	}

	_, err = dbpool.Exec(context.Background(), query, c.id, p.id, c.Datetime, c.Success, c.FailureReason, c.ContentChanged, c.Hash)
	if err != nil {
		return fmt.Errorf("unable to save crawl '%v': %v", c, err)
	}

	return nil
}
