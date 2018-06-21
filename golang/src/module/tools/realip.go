package tools

import (
    "strings"
	"net/http"
)

func GetRealClientIP(r *http.Request ) string {
	xforward := r.Header.Get("X-Forwarded-For")
	if "" == xforward {
		return strings.SplitN(r.RemoteAddr, ":", 2)[0]
	}

	return strings.SplitN(string(xforward), ",", 2)[0]
}