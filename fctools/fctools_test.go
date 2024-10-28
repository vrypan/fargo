package fctools
import (
        "testing"
)

func Test_GetUsernameByFid_280(t *testing.T) {
    var fid uint64 = 280
    var expected_name = "vrypan"
    t.Logf("Looking up usernames for fid=%v", fid)
    names, err := GetUsernameByFid(fid)
    if err != nil {
        t.Error(err)
    }
    if expected_name != names[0] && expected_name != names[1] {
        t.Errorf("'%v' not found in names: %v", expected_name, names)
    }
}

func Test_GetFidByUsername_vrypan(t *testing.T) {
    var username string = "vrypan"
    var expected_fid uint64 = 280

    t.Logf("Looking up fid for username=%v", username)
    fid, err := GetFidByUsername(username)
    if err != nil {
        t.Error(err)
    }
    if fid != expected_fid {
        t.Errorf("Expected fid=%v, got fid=%v", expected_fid, fid)
    }
}