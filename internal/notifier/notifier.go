package notifier

import "net/http"

type Notifier interface {
	SendMsg(message string, err error, extra *map[string]interface{}) error
	SendHttpMsg(errs string, tags *map[string]string, req *http.Request) error
}
