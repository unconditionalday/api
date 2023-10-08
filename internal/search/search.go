package search

type EntityDetails struct {
	Title     string
	Link      string
	Summary   string
	Thumbnail string
	Language  string
	Source    string
}

func (i EntityDetails) IsValid() bool {
	return i.Title != "" && i.Link != "" && i.Summary != "" && i.Thumbnail != "" && i.Language != "" && i.Source != ""
}

type SearchClient interface {
	FetchEntityDetails(query, locale string) (EntityDetails, error)
}
