package tasks

import (
	"ZakkBob/AskDave/crawler/hash"
	"ZakkBob/AskDave/crawler/page"
	"ZakkBob/AskDave/crawler/robots"
	"ZakkBob/AskDave/crawler/url"
)

type FailureReason = int

const (
	NoFailure FailureReason = iota
	RobotsDisallowed
)

type Results struct {
	Robots   map[string]RobotsResult `json:"robots"`
	Sitemaps map[string]string       `json:"sitemaps"`
	Pages    map[string]PageResult   `json:"pages"`
}

type PageResult struct {
	Url           url.Url       `json:"url"`
	Success       bool          `json:"success"`
	FailureReason FailureReason `json:"failure_reason,omitempty"`
	Changed       bool          `json:"changed"`
	Page          page.Page     `json:"page"`
}

type RobotsResult struct {
	Url       url.Url             `json:"robots"`
	Hash      hash.Hash           `json:"hash"`
	Changed   bool                `json:"changed"`
	Validator robots.UrlValidator `json:"validator"`
}
