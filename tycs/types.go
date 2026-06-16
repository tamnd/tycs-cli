package tycs

// Subject is one CS topic from the Teach Yourself CS curriculum.
type Subject struct {
	Rank    int    `json:"rank"     csv:"rank"     tsv:"rank"`
	Slug    string `json:"slug"     csv:"slug"     tsv:"slug"     kit:"id"`
	Title   string `json:"title"    csv:"title"    tsv:"title"`
	URL     string `json:"url"      csv:"url"      tsv:"url"`
	BookURL string `json:"book_url" csv:"book_url" tsv:"book_url"`
}

// Info is site-level stats.
type Info struct {
	Site     string `json:"site"     csv:"site"     tsv:"site"`
	Subjects int    `json:"subjects" csv:"subjects" tsv:"subjects"`
	Source   string `json:"source"   csv:"source"   tsv:"source"`
}
