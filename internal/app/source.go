package app

type entry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Source []entry

type SourceRelease struct {
	Source  Source
	Version string
}

type SourceClient interface {
	GetLatestVersion() (string, error)
	Download(version string) (SourceRelease, error)
}
