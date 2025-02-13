package main

import (
	"fmt"
	"os"
	"testing"
)

func TestAddLink(t *testing.T) {
	var p page
	expected := "https://google.com"
	u, _ := parseAbsoluteUrl(expected)
	p.addLink(u)
	got := p.links[0].String()
	if got != expected {
		t.Errorf("got '%s', expected '%s'", got, expected)
	}
}

func readPageFile(t *testing.T, n int) string {
	content, err := os.ReadFile(fmt.Sprintf("testdata/pages/page_%d.html", n))
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	return string(content)
}

func TestParseBody(t *testing.T) {
	b := readPageFile(t, 1)
	u, _ := parseAbsoluteUrl("https://www.example.com/home")
	p, err := parseBody(b, u)

	if err != nil {
		t.Fatal(err.Error())
	}

	t.Log(p)
}
