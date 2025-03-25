package orm_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/ZakkBob/AskDave/backend/orm"
	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/require"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func makeURL(t *testing.T, path string) url.Url {
	s := t.Name() + randString(10)
	s = strings.ReplaceAll(s, "/", "-")
	urlS := fmt.Sprintf("https://%s.com%s", s, path)
	u, err := url.ParseAbs(urlS)
	require.NoErrorf(t, err, "ParseAbs should not return error with url '%s'", urlS)
	return u
}

func resetDB(t *testing.T) {
	err := orm.ClearDB()
	if err != nil {
		t.Fatalf("failed to reset database: %v", err)
	}
}
