package robots_test

import (
	"ZakkBob/AskDave/crawler/robots"
	"os"
	"testing"
)

func readRobotsTxt(t *testing.T, fileName string) string {
	content, err := os.ReadFile("../testdata/" + fileName)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	return string(content)
}

func TestExtractDavebotDirectives(t *testing.T) {
	blocks := map[string]string{
		"*":           "Disallow: /",
		"davebot/0.1": "Allow: /Davebot0.1/",
		"davebot":     "Allow: /Davebot/",
		"bingbot":     "Disallow: /bingbot/",
	}

	got := robots.ExtractDavebotDirectives(blocks)
	want := "Allow: /Davebot0.1/"

	if got != want {
		t.Errorf("got %s, wanted %s", got, want)
	}
}

func TestUrlValidator(t *testing.T) {
	robotsTxt := readRobotsTxt(t, "robots.txt")

	validator, _ := robots.Parse(robotsTxt)

	testUrls := map[string]bool{
		"/allowed-dir":                              true,
		"/disallowed":                               false,
		"/disallowed/subdir":                        false,
		"/disallowed/nevermind-this-is-allowed":     true,
		"/disallowed/nevermind-this-is-allowed/lol": true,
		"/main.php":                                 false,
		"/endline":                                  false,
		"/endline/not-lol":                          true,
		"/allowed-dir/disallowed-php.php":           false,
		"/allow-subdir/":                            false,
		"/allow-subdir/e":                           true,
	}

	for url, want := range testUrls {
		got := validator.Validate(url)
		if got != want {
			t.Errorf("'%s' got %t, wanted %t", url, got, want)
		}
	}
}
