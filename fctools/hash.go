package fctools

import (
	"encoding/hex"
)

type Hash [20]byte

func (h Hash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}
