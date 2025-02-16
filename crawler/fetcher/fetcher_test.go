package fetcher_test

import (
	"ZakkBob/AskDave/crawler/fetcher"
	"ZakkBob/AskDave/crawler/url"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFileFetcher(t *testing.T) {
	var f = fetcher.FileFetcher{
		Debug: true,
	}
	s, _ := url.ParseAbs("https://fetchertest.com/index.html")
	got, _ := f.Fetch(&s)
	want := "Success"
	if got != want {
		t.Errorf("got '%s', wanted '%s'", got, want)
	}
}

func TestDummyFetcher(t *testing.T) {
	var f = fetcher.DummyFetcher{
		Debug:    true,
		Response: "Success",
	}
	s, _ := url.ParseAbs("fetchertest.com/index.html")
	got, _ := f.Fetch(&s)
	want := "Success"
	if got != want {
		t.Errorf("got '%s', wanted '%s'", got, want)
	}
}

func TestNetFetcher(t *testing.T) {
	mockHandler := http.NewServeMux()
	mockHandler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	})

	mockServer := httptest.NewServer(mockHandler)
	defer mockServer.Close()

	split := strings.Split(mockServer.URL, ".") //I apologise for this being so cursed. i really need to support ips

	fetcher := &fetcher.NetFetcher{}
	testUrl := url.Url{
		Domain: split[0] + "." + split[1],
		Tld:    split[2] + "." + split[3],
		Path:   []string{"/test"},
	}

	want := "Success"
	got, err := fetcher.Fetch(&testUrl)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	if got != want {
		t.Errorf("got '%s', wanted '%s'", got, want)
	}
}
