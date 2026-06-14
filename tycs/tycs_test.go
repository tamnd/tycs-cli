package tycs_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tamnd/tycs-cli/tycs"
)

const fakeHTML = `<html><body>
<h3 class="h3 mb0" id="programming">Programming</h3>
<div><a href="https://sicp.example.com/">SICP</a></div>
<h3 class="h3 mb0" id="architecture">Computer Architecture</h3>
<div><a href="https://csapp.example.com/">CS:APP</a></div>
</body></html>`

func newTestClient(ts *httptest.Server) *tycs.Client {
	cfg := tycs.DefaultConfig()
	cfg.BaseURL = ts.URL
	cfg.Rate = 0
	return tycs.NewClient(cfg)
}

func TestSubjects(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fakeHTML)
	}))
	defer ts.Close()

	c := newTestClient(ts)
	subjects, err := c.Subjects(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(subjects) != 2 {
		t.Fatalf("want 2, got %d", len(subjects))
	}
	if subjects[0].Slug != "programming" {
		t.Errorf("Slug[0] = %q", subjects[0].Slug)
	}
	if subjects[0].Title != "Programming" {
		t.Errorf("Title[0] = %q", subjects[0].Title)
	}
	if subjects[0].BookURL != "https://sicp.example.com/" {
		t.Errorf("BookURL[0] = %q", subjects[0].BookURL)
	}
	if subjects[0].Rank != 1 {
		t.Errorf("Rank[0] = %d, want 1", subjects[0].Rank)
	}
}
