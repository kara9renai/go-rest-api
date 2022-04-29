package params

import (
	"net/http"
	"strings"
)

func GetRouteParams(r *http.Request) []string {
	s := strings.Split(r.RequestURI, "/")
	var params []string
	for i := 0; i < len(s); i++ {
		if len(s[i]) != 0 {
			params = append(params, s[i])
		}
	}
	return params
}
