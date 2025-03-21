package orm_test

import (
	"testing"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/assert"
)

func TestSaveNewPage(t *testing.T) {
	u, err := url.ParseAbs("https://www.zakkdev.com")

	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	p := page.Page{
		Url:           u,
		Title:         "title",
		OgTitle:       "og titlke",
		OgDescription: "descriptiuon",
		OgSiteName:    "site_name",
		Links:         []url.Url{},
		Hash:          hash.Hashs(""),
	}

	_, err = orm.SaveNewSite("https://www.zakkdev.com")
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	_, err = orm.SaveNewPage(p)
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	ormPage, err := orm.PageByUrl(u.String())
	if err != nil {
		t.Errorf("didn't expect an error: %v", err)
	}

	assert.Equal(t, p.Url.String(), ormPage.Url.String())
	assert.Equal(t, p.Title, ormPage.Title)
	assert.Equal(t, p.OgTitle, ormPage.OgTitle)
	assert.Equal(t, p.OgDescription, ormPage.OgDescription)
	assert.Equal(t, p.OgSiteName, ormPage.OgSiteName)
	assert.Equal(t, p.Links, ormPage.Links)
	assert.Equal(t, p.Hash.String(), ormPage.Hash.String())
}
