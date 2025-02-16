//-------------------------------------------------------------//
// An interface for fetching webpages and some implementations //
//-------------------------------------------------------------//

package fetcher

import (
	"ZakkBob/AskDave/crawler/url"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Fetcher interface {
	Fetch(*url.Url) (string, error)
}

type NetFetcher struct{}

func (f *NetFetcher) Fetch(u *url.Url) (string, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
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
func (f *DummyFetcher) Fetch(u *url.Url) (string, error) {
	totalDelay := f.sleep()

	if f.Debug {
		fmt.Printf("fetched dummy url '%s' after %s\n", u.String(), totalDelay.String())
	}
	return f.Response, nil
}

type FileFetcher struct {
	Debug bool
}

func (f *FileFetcher) Fetch(u *url.Url) (string, error) {
	path := "../testdata/sites/"
	if u.ProtocolString() != "" {
		path += u.ProtocolString() + "/"
	}
	path += u.FQDN() + u.PathString()
	content, err := os.ReadFile(path)
	if f.Debug {
		fmt.Printf("fetching file url '%s'\n", u.String())
		fmt.Printf("path: '%s'\n", path)
		fmt.Printf("content: '%s'\n", content)
	}
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	return string(content), nil
}
