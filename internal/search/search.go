package search

type ContextDetails struct {
	Title     string
	Link      string
	Summary   string
	Thumbnail string
	Language  string
}

func (i ContextDetails) IsValid() bool {
	return i.Title != "" && i.Link != "" && i.Summary != "" && i.Thumbnail != "" && i.Language != ""
}

type SearchClient interface {
	FetchContextDetails(query, locale string) (ContextDetails, error)
}
