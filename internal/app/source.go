package app

type entry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Source []entry

type SourceService interface {
	Download(path string) (Source, error)
}
