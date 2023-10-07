package app

type entry struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Source []entry

type SourceRelease struct {
	Data         Source
	Version      string
	LastUpdateAt string
}

type SourceClient interface {
	GetLatestVersion() (string, error)
	Download(version string) (SourceRelease, error)
}
