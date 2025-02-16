package tasks

import (
	"ZakkBob/AskDave/crawler/url"
	"sync"
)

type Tasks struct {
	Robots   taskSlice `json:"robots"` // Bit clunky
	Sitemaps taskSlice `json:"sitemaps"`
	Pages    taskSlice `json:"pages"`
}

type taskSlice struct {
	Mu    sync.Mutex `json:"-"`
	Slice []url.Url  `json:"slice"`
}

// Returns next url in slice, returns nil if slice is empty
func (t *taskSlice) Next() *url.Url {
	t.Mu.Lock()
	defer t.Mu.Unlock()
	if len(t.Slice) == 0 {
		return nil
	}
	u := &t.Slice[0]
	t.Slice = t.Slice[1:len(t.Slice)]
	return u
}
