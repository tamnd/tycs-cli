package tycs

import (
	"testing"

	"github.com/tamnd/any-cli/kit"
)

// These tests are offline: they exercise the URI driver's pure string functions
// and the host wiring, which need no network. The client's HTTP behaviour is
// covered in tycs_test.go.

func TestDomainInfo(t *testing.T) {
	info := Domain{}.Info()
	if info.Scheme != "tycs" {
		t.Errorf("Scheme = %q, want tycs", info.Scheme)
	}
	if len(info.Hosts) == 0 || info.Hosts[0] != Host {
		t.Errorf("Hosts = %v, want [%s]", info.Hosts, Host)
	}
	if info.Identity.Binary != "tycs" {
		t.Errorf("Identity.Binary = %q, want tycs", info.Identity.Binary)
	}
}

func TestClassify(t *testing.T) {
	cases := []struct{ in, typ, id string }{
		{"programming", "subject", "programming"},
		{"#databases", "subject", "databases"},
	}
	for _, tc := range cases {
		typ, id, err := Domain{}.Classify(tc.in)
		if err != nil || typ != tc.typ || id != tc.id {
			t.Errorf("Classify(%q) = (%q, %q, %v), want (%q, %q, nil)",
				tc.in, typ, id, err, tc.typ, tc.id)
		}
	}
}

func TestLocate(t *testing.T) {
	got, err := Domain{}.Locate("subject", "operating-systems")
	want := "https://" + Host + "/#operating-systems"
	if err != nil || got != want {
		t.Errorf("Locate = (%q, %v), want (%q, nil)", got, err, want)
	}
}

func TestHostWiring(t *testing.T) {
	h, err := kit.Open()
	if err != nil {
		t.Fatal(err)
	}

	s := &Subject{Rank: 1, Slug: "programming", Title: "Programming",
		URL: "https://" + Host + "/#programming", BookURL: "https://example.com/"}
	u, err := h.Mint(s)
	if err != nil {
		t.Fatalf("Mint: %v", err)
	}
	if want := "tycs://subject/programming"; u.String() != want {
		t.Errorf("Mint = %q, want %q", u.String(), want)
	}
}
