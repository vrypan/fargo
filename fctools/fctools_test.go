package fctools

import (
	"testing"
)

func Test_GetFidByUsername_vrypan(t *testing.T) {
    var username string = "vrypan"
    var expected_fid uint64 = 280

    t.Logf("Looking up fid for username=%v", username)
    hub := NewFarcasterHub(); defer hub.Close()
    fid, err := hub.GetFidByUsername(username)
    if err != nil {
        t.Error(err)
    }
    if fid != expected_fid {
        t.Errorf("Expected fid=%v, got fid=%v", expected_fid, fid)
    }
}
