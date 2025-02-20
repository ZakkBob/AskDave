package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/ZakkBob/AskDave/gocommon/hash"
	"github.com/ZakkBob/AskDave/gocommon/page"
	"github.com/ZakkBob/AskDave/gocommon/robots"
	"github.com/ZakkBob/AskDave/gocommon/tasks"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func main() {
	var err error
	dbpool, err = pgxpool.New(context.Background(), "postgres://postgres:password@127.0.0.1:5432/postgres")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// t, err := nextTasks(2)

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// 	os.Exit(1)
	// }

	// data, err := json.MarshalIndent(t, "", "  ")

	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to serialise tasks: %v\n", err)
	// 	os.Exit(1)
	// }

	// fmt.Println(string(data))

	r := tasks.Results{
		Robots: make(map[string]*tasks.RobotsResult),
		Pages:  make(map[string]*tasks.PageResult),
	}

	mateiValidator, _ := robots.Parse(`
#
# robots.txt
#
# This file is to prevent the crawling and indexing of certain parts
# of your site by web crawlers and spiders run by sites like Yahoo!
# and Google. By telling these "robots" where not to go on your site,
# you save bandwidth and server resources.
#
# This file will be ignored unless it is at the root of your host:
# Used:    http://example.com/robots.txt
# Ignored: http://example.com/site/robots.txt
#
# For more information about the robots.txt standard, see:
# http://www.robotstxt.org/robotstxt.html

User-agent: *
Crawl-delay: 10
# CSS, JS, Images
Allow: /misc/*.css$
Allow: /misc/*.css?
Allow: /misc/*.js$
Allow: /misc/*.js?
Allow: /misc/*.gif
Allow: /misc/*.jpg
Allow: /misc/*.jpeg
Allow: /misc/*.png
Allow: /modules/*.css$
Allow: /modules/*.css?
Allow: /modules/*.js$
Allow: /modules/*.js?
Allow: /modules/*.gif
Allow: /modules/*.jpg
Allow: /modules/*.jpeg
Allow: /modules/*.png
Allow: /profiles/*.css$
Allow: /profiles/*.css?
Allow: /profiles/*.js$
Allow: /profiles/*.js?
Allow: /profiles/*.gif
Allow: /profiles/*.jpg
Allow: /profiles/*.jpeg
Allow: /profiles/*.png
Allow: /themes/*.css$
Allow: /themes/*.css?
Allow: /themes/*.js$
Allow: /themes/*.js?
Allow: /themes/*.gif
Allow: /themes/*.jpg
Allow: /themes/*.jpeg
Allow: /themes/*.png
# Directories
Disallow: /includes/
Disallow: /misc/
Disallow: /modules/
Disallow: /profiles/
Disallow: /scripts/
Disallow: /themes/
# Files
Disallow: /CHANGELOG.txt
Disallow: /cron.php
Disallow: /INSTALL.mysql.txt
Disallow: /INSTALL.pgsql.txt
Disallow: /INSTALL.sqlite.txt
Disallow: /install.php
Disallow: /INSTALL.txt
Disallow: /LICENSE.txt
Disallow: /MAINTAINERS.txt
Disallow: /update.php
Disallow: /UPGRADE.txt
Disallow: /xmlrpc.php
# Paths (clean URLs)
Disallow: /admin/
Disallow: /comment/reply/
Disallow: /filter/tips/
Disallow: /node/add/
Disallow: /search/
Disallow: /user/register/
Disallow: /user/password/
Disallow: /user/login/
Disallow: /user/logout/
# Paths (no clean URLs)
Disallow: /?q=admin/
Disallow: /?q=comment/reply/
Disallow: /?q=filter/tips/
Disallow: /?q=node/add/
Disallow: /?q=search/
Disallow: /?q=user/password/
Disallow: /?q=user/register/
Disallow: /?q=user/login/
Disallow: /?q=user/logout/`)

	u, _ := url.ParseAbs("https://mateishome.page")
	matei := tasks.RobotsResult{
		Url:           &u,
		Success:       true,
		FailureReason: tasks.NoFailure,
		Changed:       true,
		Hash:          hash.Hash{},
		Validator:     &mateiValidator,
	}

	r.Robots["https://mateishome.page"] = &matei

	u1, _ := url.ParseAbs("https://www.google.com/e")
	u2, _ := url.ParseAbs("https://mateishome.page/api")

	r.Pages["https://mateishome.page"] = &tasks.PageResult{
		Url:           &u,
		Success:       true,
		FailureReason: tasks.NoFailure,
		Changed:       true,
		Page: &page.Page{
			Url:           u,
			Title:         "Matei's Epic Home Page",
			OgTitle:       "",
			OgDescription: "",
			OgSiteName:    "",
			Hash:          hash.Hash{},
			Links: []url.Url{
				u1, u2,
			},
		},
	}

	data, err := json.MarshalIndent(r, "", "  ")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to serialise results: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))

	err = saveResults(&r)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	validator, err := validatorByUrl("https://mateishome.page")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create validator: %v\n", err)
		os.Exit(1)
	}
	data, err = json.MarshalIndent(validator, "", "  ")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to serialise validator: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(data))
}

// 1 | https://mateishome.page
// 2 | https://games.zakkcarpenter.com
// 3 | https://greendungarees.org.uk
// 4 | https://zakkcarpenter.com
