package orm_test

// import (
// 	"testing"

// 	"github.com/ZakkBob/AskDave/backend/orm"
// 	"github.com/ZakkBob/AskDave/gocommon/hash"
// 	"github.com/ZakkBob/AskDave/gocommon/page"
// 	"github.com/ZakkBob/AskDave/gocommon/url"
// 	"github.com/stretchr/testify/assert"
// )

// func TestSaveNewPage(t *testing.T) {
// 	u, err := url.ParseAbs("www.test.com/TestSaveNewPage")

// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	p := page.Page{
// 		Url:           u,
// 		Title:         "title",
// 		OgTitle:       "og titlke",
// 		OgDescription: "descriptiuon",
// 		OgSiteName:    "site_name",
// 		Links:         []url.Url{},
// 		Hash:          hash.Hashs(""),
// 	}

// 	_, err = orm.SaveNewSite("www.test.com")
// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	_, err = orm.SaveNewPage(p)
// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	ormPage, err := orm.PageByUrl(u.String())
// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	assert.Equal(t, p.Url.String(), ormPage.Url.String())
// 	assert.Equal(t, p.Title, ormPage.Title)
// 	assert.Equal(t, p.OgTitle, ormPage.OgTitle)
// 	assert.Equal(t, p.OgDescription, ormPage.OgDescription)
// 	assert.Equal(t, p.OgSiteName, ormPage.OgSiteName)
// 	assert.Equal(t, p.Links, ormPage.Links)
// 	assert.Equal(t, p.Hash.String(), ormPage.Hash.String())
// }

// func TestSavePage(t *testing.T) {
// 	u, err := url.ParseAbs("www.test.com/TestSavePage")

// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	p := page.Page{
// 		Url:           u,
// 		Title:         "title",
// 		OgTitle:       "og titlke",
// 		OgDescription: "descriptiuon",
// 		OgSiteName:    "site_name",
// 		Links:         []url.Url{},
// 		Hash:          hash.Hashs(""),
// 	}

// 	_, err = orm.SaveNewSite("www.test.com")
// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	ormPage, err := orm.SaveNewPage(p)
// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	links, err := url.ParseMany([]string{"www.TestSaveNewPage.com", "www.TestSaveNewPage2.com/e", "www.TestSaveNewPage.com/e"})
// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	ormPage.Title = "Title 2"
// 	ormPage.OgTitle = "Og Title 2"
// 	ormPage.OgDescription = "Desc 2"
// 	ormPage.OgSiteName = "Name 2"
// 	ormPage.Links = links

// 	err = ormPage.Save(true)
// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	ormPage2, err := orm.PageByUrl(u.String())
// 	if err != nil {
// 		t.Errorf("didn't expect an error: %v", err)
// 	}

// 	assert.Equal(t, ormPage.Url.String(), ormPage2.Url.String())
// 	assert.Equal(t, ormPage.Title, ormPage2.Title)
// 	assert.Equal(t, ormPage.OgTitle, ormPage2.OgTitle)
// 	assert.Equal(t, ormPage.OgDescription, ormPage2.OgDescription)
// 	assert.Equal(t, ormPage.OgSiteName, ormPage2.OgSiteName)
// 	assert.Equal(t, ormPage.Links, ormPage2.Links)
// 	assert.Equal(t, ormPage.Hash.String(), ormPage2.Hash.String())
// }