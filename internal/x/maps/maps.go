package maps

import "strconv"

/*
Update map a using map b
*/
func Update(a map[string]string, b map[string]interface{}) {
	for k, v := range b {
		switch t := v.(type) {
		case int:
			a[k] = strconv.Itoa(t)
		case string:
			a[k] = t
		}
	}
}
