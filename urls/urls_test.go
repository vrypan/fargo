package urls

import (
	"testing"
)

func Test_Url_Image(t *testing.T) {
	u := NewUrl("https://www.vrypan.net/assets/img/vrypan.jpg").UpdateExt().UpdateType()
	if u.Extension != "jpg" {
		t.Fatalf("Expected jpg, got [%s]", u.Extension)
	}
	if u.ContentType != "image/jpeg" {
		t.Fatalf("Expected image/jpeg, got [%s]", u.ContentType)
	}
	t.Logf("Link: %s\n", u.Link)
	t.Logf("Id: %s\n", u.Id())
	t.Logf("Extension: %s\n", u.Extension)
	t.Logf("ContentType: %s\n", u.ContentType)
	t.Logf("Filename: %s\n", u.Filename())
	t.Logf("String: %s\n", u)
}

func Test_Url_TLD(t *testing.T) {
	u := NewUrl("https://www.vrypan.net/").UpdateExt().UpdateType()
	if u.Extension != "" {
		t.Fatalf("Expected none, got [%s]", u.Extension)
	}
	if u.ContentType != "text/html" {
		t.Fatalf("Expected text/html, got [%s]", u.ContentType)
	}
	t.Logf("Link: %s\n", u.Link)
	t.Logf("Id: %s\n", u.Id())
	t.Logf("Extension: %s\n", u.Extension)
	t.Logf("ContentType: %s\n", u.ContentType)
	t.Logf("Filename: %s\n", u.Filename())
	t.Logf("String: %s\n", u)
}

func Test_Url_Path(t *testing.T) {
	u := NewUrl("https://blog.vrypan.net/archive/").UpdateExt().UpdateType()
	if u.Extension != "" {
		t.Fatalf("Expected none, got [%s]", u.Extension)
	}
	if u.ContentType != "text/html" {
		t.Fatalf("Expected text/html, got [%s]", u.ContentType)
	}
	t.Logf("Link: %s\n", u.Link)
	t.Logf("Id: %s\n", u.Id())
	t.Logf("Extension: %s\n", u.Extension)
	t.Logf("ContentType: %s\n", u.ContentType)
	t.Logf("Filename: %s\n", u.Filename())
	t.Logf("String: %s\n", u)
}
