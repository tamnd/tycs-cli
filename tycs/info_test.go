package tycs_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSiteInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fakeHTML)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	info, err := c.SiteInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if info.Subjects != 2 {
		t.Errorf("Subjects = %d, want 2", info.Subjects)
	}
	if info.Site == "" {
		t.Error("Site should not be empty")
	}
}
