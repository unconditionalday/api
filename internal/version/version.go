package version

type Versioning interface {
	Lower(v1, v2 string) (bool, error)
}
