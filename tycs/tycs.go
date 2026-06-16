// Package tycs is the library behind the tycs CLI.
package tycs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

const DefaultUserAgent = "tycs-cli/dev (+https://github.com/tamnd/tycs-cli)"

// Config holds per-client settings.
type Config struct {
	BaseURL   string
	Rate      time.Duration
	Timeout   time.Duration
	Retries   int
	UserAgent string
}

// DefaultConfig returns a Config pointed at teachyourselfcs.com.
func DefaultConfig() Config {
	return Config{
		BaseURL:   "https://teachyourselfcs.com",
		Rate:      500 * time.Millisecond,
		Timeout:   30 * time.Second,
		Retries:   3,
		UserAgent: DefaultUserAgent,
	}
}

// Client fetches data from teachyourselfcs.com.
type Client struct {
	cfg  Config
	http *http.Client
	last time.Time
}

// NewClient returns a Client configured by cfg.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
	}
}

var (
	h3Re   = regexp.MustCompile(`<h3 class="h3 mb0" id="([\w-]+)">([^<]+)</h3>`)
	hrefRe = regexp.MustCompile(`href="([^"]+)"`)
)

// Subjects fetches the teachyourselfcs.com homepage and returns all 9 CS
// subject guides in document order.
func (c *Client) Subjects(ctx context.Context) ([]*Subject, error) {
	body, err := c.get(ctx, c.cfg.BaseURL+"/")
	if err != nil {
		return nil, err
	}
	html := string(body)

	matches := h3Re.FindAllStringIndex(html, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("no subject headings found")
	}

	subs := h3Re.FindAllStringSubmatch(html, -1)
	var subjects []*Subject
	for i, loc := range matches {
		slug := subs[i][1]
		title := subs[i][2]

		// Slice the block from this h3 to the next (or end of document).
		start := loc[0]
		end := len(html)
		if i+1 < len(matches) {
			end = matches[i+1][0]
		}
		block := html[start:end]

		bookURL := ""
		if m := hrefRe.FindStringSubmatch(block); m != nil {
			bookURL = m[1]
		}

		subjects = append(subjects, &Subject{
			Rank:    i + 1,
			Slug:    slug,
			Title:   title,
			URL:     "https://teachyourselfcs.com/#" + slug,
			BookURL: bookURL,
		})
	}
	return subjects, nil
}

// SiteInfo returns site-level stats.
func (c *Client) SiteInfo(ctx context.Context) (*Info, error) {
	subjects, err := c.Subjects(ctx)
	if err != nil {
		return nil, err
	}
	return &Info{
		Site:     Host,
		Subjects: len(subjects),
		Source:   c.cfg.BaseURL,
	}, nil
}

func (c *Client) get(ctx context.Context, url string) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= c.cfg.Retries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}
		body, retry, err := c.do(ctx, url)
		if err == nil {
			return body, nil
		}
		lastErr = err
		if !retry {
			return nil, err
		}
	}
	return nil, fmt.Errorf("get %s: %w", url, lastErr)
}

func (c *Client) do(ctx context.Context, url string) ([]byte, bool, error) {
	c.pace()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, true, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
		return nil, true, fmt.Errorf("http %d", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("http %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, true, err
	}
	return b, false, nil
}

func (c *Client) pace() {
	if c.cfg.Rate <= 0 {
		return
	}
	if wait := c.cfg.Rate - time.Since(c.last); wait > 0 {
		time.Sleep(wait)
	}
	c.last = time.Now()
}

func backoff(attempt int) time.Duration {
	return min(time.Duration(attempt)*500*time.Millisecond, 5*time.Second)
}
