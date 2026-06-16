package tycs

import (
	"context"
	"fmt"
	"strings"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

// Host is the site this client talks to.
const Host = "teachyourselfcs.com"

// BaseURL is the root every request is built from.
const BaseURL = "https://" + Host

func init() { kit.Register(Domain{}) }

// Domain is the tycs driver for the kit host (ant).
type Domain struct{}

func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme: "tycs",
		Hosts:  []string{Host},
		Identity: kit.Identity{
			Binary: "tycs",
			Short:  "Browse the Teach Yourself CS curriculum from the command line",
			Long: `tycs reads the Teach Yourself CS curriculum from teachyourselfcs.com,
shapes it into clean records, and prints output that pipes into the rest of your
tools. No API key, nothing to run alongside it.`,
			Site: Host,
			Repo: "https://github.com/tamnd/tycs-cli",
		},
	}
}

func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)

	// list: all subjects.
	kit.Handle(app, kit.OpMeta{Name: "list", Group: "read", List: true,
		Summary: "List all CS subject guides",
		URIType: "subject"}, listSubjects)

	// subjects: alias for list.
	kit.Handle(app, kit.OpMeta{Name: "subjects", Group: "read", List: true,
		Summary: "List all CS subject guides from Teach Yourself CS",
		URIType: "subject"}, listSubjects)

	// topic: fetch one subject by slug.
	kit.Handle(app, kit.OpMeta{Name: "topic", Group: "read", Single: true,
		Summary: "Fetch a CS subject guide by slug", URIType: "subject",
		Args: []kit.Arg{{Name: "slug", Help: "subject slug (e.g. programming)"}}}, getSubject)

	// subject: resolver op.
	kit.Handle(app, kit.OpMeta{Name: "subject", Group: "read", Single: true,
		Summary: "Fetch a CS subject guide by slug", URIType: "subject",
		Resolver: true,
		Args:     []kit.Arg{{Name: "slug", Help: "subject slug (e.g. programming)"}}}, getSubject)

	// info: site stats.
	kit.Handle(app, kit.OpMeta{Name: "info", Group: "read", Single: true,
		Summary: "Print site stats (subject count)"}, getSiteInfo)
}

func newClient(_ context.Context, cfg kit.Config) (any, error) {
	dcfg := DefaultConfig()
	if cfg.UserAgent != "" {
		dcfg.UserAgent = cfg.UserAgent
	}
	if cfg.Rate > 0 {
		dcfg.Rate = cfg.Rate
	}
	if cfg.Retries > 0 {
		dcfg.Retries = cfg.Retries
	}
	if cfg.Timeout > 0 {
		dcfg.Timeout = cfg.Timeout
	}
	return NewClient(dcfg), nil
}

type subjectRef struct {
	Slug   string  `kit:"arg" help:"subject slug (e.g. programming)"`
	Client *Client `kit:"inject"`
}

type subjectsIn struct {
	Client *Client `kit:"inject"`
}

type infoIn struct {
	Client *Client `kit:"inject"`
}

func getSubject(ctx context.Context, in subjectRef, emit func(*Subject) error) error {
	subjects, err := in.Client.Subjects(ctx)
	if err != nil {
		return err
	}
	for _, s := range subjects {
		if s.Slug == in.Slug {
			return emit(s)
		}
	}
	return fmt.Errorf("subject %q not found", in.Slug)
}

func listSubjects(ctx context.Context, in subjectsIn, emit func(*Subject) error) error {
	subjects, err := in.Client.Subjects(ctx)
	if err != nil {
		return err
	}
	for _, s := range subjects {
		if err := emit(s); err != nil {
			return err
		}
	}
	return nil
}

func getSiteInfo(ctx context.Context, in infoIn, emit func(*Info) error) error {
	info, err := in.Client.SiteInfo(ctx)
	if err != nil {
		return err
	}
	return emit(info)
}

func (Domain) Classify(input string) (uriType, id string, err error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", "", errs.Usage("unrecognized tycs reference: %q", input)
	}
	// Treat any slug-like input as a subject id.
	slug := strings.Trim(input, "/#")
	if slug == "" {
		return "", "", errs.Usage("unrecognized tycs reference: %q", input)
	}
	return "subject", slug, nil
}

func (Domain) Locate(uriType, id string) (string, error) {
	if uriType != "subject" {
		return "", errs.Usage("tycs has no resource type %q", uriType)
	}
	return fmt.Sprintf("%s/#%s", BaseURL, id), nil
}
