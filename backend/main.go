package main

import (
	"time"

	"github.com/ZakkBob/AskDave/backend/crawlerapi"
	"github.com/ZakkBob/AskDave/backend/orm"
)

func main() {
	orm.Connect("")
	defer orm.Close()

	crawlerapi.Init()

	// u, _ := url.ParseAbs("google.com")

	// p := page.Page{
	// 	Url:           u,
	// 	Title:         "hekki",
	// 	OgTitle:       "lol",
	// 	OgDescription: "egg",
	// 	OgSiteName:    "site name",
	// 	Hash:          hash.Hashs(""),
	// }

	// orm.SaveNewPage(p)

	time.Sleep(time.Second * 100)
}
