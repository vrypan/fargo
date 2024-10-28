package localdb

import (
    "testing"
    "strconv"
    "encoding/json"
)

func TestBasic(t *testing.T) {
    Open()
    defer Close()
    t.Logf("kv_store.top %v", kv_store.Top)
    Set("key1", "value1")
    t.Logf("kv_store.top %v", kv_store.Top)
    v, _ := Get("key1")
    if v != "value1" {
        t.Errorf("Expected 'value1', got %#v", v)
    }
    t.Logf("kv_store.top %v", kv_store.Top)

    Set("key2", "value2")
    t.Logf("kv_store.top %v", kv_store.Top)

    v, _ = Get("key1")
    if v != "value1" {
        t.Errorf("Expected 'value1', got %#v", v)
    }

    var b []byte
    b, err := json.Marshal(kv_store.Kv); if err != nil {
        panic(err)
    }
    t.Log(string(b))
    t.Logf("%#v", kv_store)
}


func TestInsert(t *testing.T) {
    Open()
    for i := range 100 {
        Set(strconv.Itoa(i), "Hello, " + strconv.Itoa(i*10))
    }
    Close()
    Open()
    defer Close()
    v,_ := Get("10")
    if v != "Hello, 100" {
        t.Errorf("Expected 'Hello 100', got %#v", v)
    }
}