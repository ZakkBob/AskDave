//---------------------------------------//
// A minimal abstraction over crypto/md5 //
//---------------------------------------//

package hash

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

type Hash [16]byte

// Returns md5 hash of a string
func Hashs(s string) Hash {
	return md5.Sum([]byte(s))
}

func (h *Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

func (h *Hash) UnmarshalJSON(data []byte) error {
	var s string
	json.Unmarshal(data, &s)
	temp, err := StrToHash(s)
	if err != nil {
		return err
	}
	*h = temp
	return nil
}

func StrToHash(s string) (Hash, error) {
	d, err := hex.DecodeString(s)
	if err != nil {
		return Hash{}, err
	}

	return Hash(d), nil
}

func (h *Hash) String() string {
	return hex.EncodeToString(h[:])
}
