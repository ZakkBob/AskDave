//------------------------------------------------------//
// An interface for fetching webpages, and some structs //
//------------------------------------------------------//

package fetcher

import (
	"embed"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/ZakkBob/AskDave/gocommon/url"
)

//go:embed testsites/*
var files embed.FS

type Fetcher interface {
	Fetch(string) (Response, error)
}

type Response struct {
	Body       string
	StatusCode int
}

type NetFetcher struct {
	Debug bool
}

func (f *NetFetcher) Fetch(u string) (Response, error) {
	resp, err := http.Get(u)
	if f.Debug {
		fmt.Printf("fetching url '%s'\n", u)
	}
	if err != nil {
		return Response{}, err // Get was unsuccessful, url probably doesnt exist or something, who knows
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	if f.Debug {
		fmt.Printf("fetched url '%s'\n", u)
	}
	return Response{
		Body:       string(body),
		StatusCode: resp.StatusCode,
	}, nil
}

type DummyFetcher struct {
	Response  string
	Delay     time.Duration
	RandDelay time.Duration
	Debug     bool
}

func (f *DummyFetcher) sleep() time.Duration {
	if f.Delay == 0 && f.RandDelay == 0 {
		return 0
	}

	if f.RandDelay == 0 {
		time.Sleep(f.Delay)
		return f.Delay
	}

	extraDelay := time.Duration(rand.Int() % (int(f.RandDelay)))
	totalDelay := f.Delay + extraDelay
	time.Sleep(totalDelay)
	return totalDelay
}

// returns DummyFetcher.response, the url parameter is ignored
func (f *DummyFetcher) Fetch(u string) (Response, error) {
	totalDelay := f.sleep()

	if f.Debug {
		fmt.Printf("fetched dummy url '%s' after %s\n", u, totalDelay.String())
	}
	return Response{
		Body:       f.Response,
		StatusCode: 200,
	}, nil
}

type FileFetcher struct {
	Delay     time.Duration
	RandDelay time.Duration
	Debug     bool
}

func (f *FileFetcher) sleep() time.Duration {
	if f.Delay == 0 && f.RandDelay == 0 {
		return 0
	}

	if f.RandDelay == 0 {
		time.Sleep(f.Delay)
		return f.Delay
	}

	extraDelay := time.Duration(rand.Int() % (int(f.RandDelay)))
	totalDelay := f.Delay + extraDelay
	time.Sleep(totalDelay)
	return totalDelay
}

func (f *FileFetcher) Fetch(s string) (Response, error) {
	totalDelay := f.sleep()

	u, err := url.ParseAbs(s)
	if err != nil {
		return Response{}, fmt.Errorf("fetching file: %w", err)
	}

	path := "testsites/"
	path += u.FQDN() + u.PathString()
	content, err := files.ReadFile(path)
	if f.Debug {
		fmt.Printf("fetching file url '%s' after %s\n", u.String(), totalDelay.String())
		fmt.Printf("path: '%s'\n", path)
		// fmt.Printf("content: '%s'\n", content)
	}
	if err != nil {
		fmt.Printf("fetching file url '%s': %w\n", u.String(), err)
		return Response{
			Body:       "",
			StatusCode: 404,
		}, nil
	}
	return Response{
		Body:       string(content),
		StatusCode: 200,
	}, nil
}
