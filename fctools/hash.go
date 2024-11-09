package fctools

import (
	"encoding/hex"
)

type tHash [20]byte

func (h tHash) String() string {
	return "0x" + hex.EncodeToString(h[:])
}
