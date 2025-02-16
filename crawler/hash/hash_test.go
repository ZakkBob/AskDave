package hash_test

import (
	"ZakkBob/AskDave/crawler/hash"
	"encoding/json"
	"testing"
)

func TestHashs(t *testing.T) {
	s := "Testing..."
	hash := hash.Hashs(s)
	got := hash.String()
	expected := "9d770c909c2c69b09eae2372c4cf405d"
	if got != expected {
		t.Errorf("expected '%s', but got '%s'", expected, got)
	}
}

func TestMarshal(t *testing.T) {
	h := hash.Hashs("Testing...")
	got, _ := json.MarshalIndent(&h, "", "  ")
	expected := `"9d770c909c2c69b09eae2372c4cf405d"`

	if string(got) != expected {
		t.Errorf("expected '%s', but got '%s'", expected, got)
	}
}

func TestUnmarshal(t *testing.T) {
	data := []byte(`"9d770c909c2c69b09eae2372c4cf405d"`)
	var h hash.Hash
	json.Unmarshal(data, &h)

	want := "9d770c909c2c69b09eae2372c4cf405d"

	if h.String() != want {
		t.Errorf("expected '%s', but got '%s'", want, h.String())
	}
}
