package orm_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ZakkBob/AskDave/gocommon/url"
	"github.com/stretchr/testify/require"
)

func uniqueURL(t *testing.T) url.Url {
	testName := t.Name()
	testName = strings.ReplaceAll(testName, "/", "-")
	urlS := fmt.Sprintf("https://test-%s.com", testName)
	u, err := url.ParseAbs(urlS)
	require.NoErrorf(t, err, "ParseAbs should not return error with url '%s'", urlS)
	return u
}
