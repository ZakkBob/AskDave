package url_test

import (
	"ZakkBob/AskDave/crawler/url"
	"encoding/json"
	"testing"
)

func TestUrlString(t *testing.T) {
	testUrls := []url.Url{
		{
			url.HttpsProtocol, "www", "example", "com", 8080,
			[]string{
				"privacy",
			}, true,
		},
		{
			Domain: "example",
			Tld:    "com",
		},
		{
			Protocol: url.HttpProtocol,
			Domain:   "example",
			Tld:      "com",
			Port:     80,
		},
	}
	expectedUrls := []string{
		"https://www.example.com:8080/privacy/",
		"example.com",
		"http://example.com",
	}

	for i, u := range testUrls {
		got := u.String()
		expected := expectedUrls[i]
		if got != expected {
			t.Errorf("got '%s', expected '%s'", got, expected)
		}
	}
}

func TestParseAbsoluteUrl(t *testing.T) {
	testStrings := []string{
		"https://www.google.com:1234/privacy-policy/help/",
		"http://www.google.com:1234/",
		"https://www.google.com:1234",
		"www.google.com/privacy-policy/help",
		"google.com/privacy-policy/help",
		"http://google.com/privacy-policy/",
		"https://example.com/robots.txt",
	}

	for _, expected := range testStrings {
		parsed, err := url.ParseAbs(expected)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		got := parsed.String()
		if got != expected {
			t.Errorf("got '%s', expected '%s'", got, expected)
		}
	}
}

func TestParseRelativeUrl(t *testing.T) {
	baseUrl, _ := url.ParseAbs("example.com/subdir")

	testStrings := []string{
		"home.php",
		"/home.html",
		"./home",
		"http://example.org/examples",
		"http://example.org",
		"../subdir-2",
		"../../..",
		"../../e/./../subdir/",
	}

	expectedStrings := []string{
		"example.com/subdir/home.php",
		"example.com/home.html",
		"example.com/subdir/home",
		"http://example.org/examples",
		"http://example.org",
		"example.com/subdir-2",
		"example.com",
		"example.com/subdir/",
	}

	for i, testString := range testStrings {
		parsed, err := url.ParseRel(testString, baseUrl)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		got := parsed.String()
		expected := expectedStrings[i]
		if got != expected {
			t.Errorf("got '%s', expected '%s'", got, expected)
			t.Error(parsed)
		}
	}
}

func TestIsFile(t *testing.T) {
	testStrings := []string{
		"example.com/subdir/home.php",
		"example.com/home.html",
		"example.com/subdir/home",
		"http://example.org/examples",
		"http://example.org",
		"example.com/subdir-2",
		"example.com",
		"example.com/subdir/",
	}

	wantStrings := []bool{
		true,
		true,
		false,
		false,
		false,
		false,
		false,
		false,
	}

	for i, testString := range testStrings {
		parsed, err := url.ParseAbs(testString)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		got := parsed.IsFile()
		want := wantStrings[i]
		if got != want {
			t.Errorf("got '%t', expected '%t'", got, want)
		}
	}
}

func TestParseRelativeUrlWithFile(t *testing.T) {
	baseUrl, _ := url.ParseAbs("example.com/subdir/index.html")

	testStrings := []string{
		"home.php",
		"/home.html",
		"./home",
		"../subdir-2",
		"../../..",
		"../../e/./../subdir/",
	}

	expectedStrings := []string{
		"example.com/subdir/home.php",
		"example.com/home.html",
		"example.com/subdir/home",
		"example.com/subdir-2",
		"example.com",
		"example.com/subdir/",
	}

	for i, testString := range testStrings {
		parsed, err := url.ParseRel(testString, baseUrl)
		if err != nil {
			t.Error(err.Error())
			continue
		}
		got := parsed.String()
		expected := expectedStrings[i]
		if got != expected {
			t.Errorf("got '%s', expected '%s'", got, expected)
			t.Error(parsed)
		}
	}
}

func TestMarshal(t *testing.T) {
	u, _ := url.ParseAbs("https://www.example.com/e")
	got, _ := json.MarshalIndent(&u, "", "  ")
	want := `"https://www.example.com/e"`

	if string(got) != want {
		t.Errorf("expected '%s', but got '%s'", want, string(got))
	}
}

func TestUnmarshal(t *testing.T) {
	data := []byte(`"https://www.example.com/e"`)
	var u url.Url

	json.Unmarshal(data, &u)

	want := "https://www.example.com/e"

	if u.String() != want {
		t.Errorf("expected '%s', but got '%s'", want, u.String())
	}
}
