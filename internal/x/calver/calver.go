package calver

import "github.com/loadsmart/calver-go/calver"

type CalVer struct{}

func New() *CalVer {
	return &CalVer{}
}

func (c *CalVer) Lower(v1, v2 string) (bool, error) {
	pattern := "YYYY.MM.DD"

	v1New, err := calver.Parse(pattern, v1)
	if err != nil {
		return false, err
	}

	v2New, err := calver.Parse(pattern, v2)
	if err != nil {
		return false, err
	}

	if res := v1New.CompareTo(v2New); res != 1 {
		return true, nil
	}

	return false, nil
}
