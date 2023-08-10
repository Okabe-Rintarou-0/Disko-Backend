package utils

import (
	"net/http"
	"strings"
)

func GetToken(r *http.Request) (string, bool) {
	token := r.Header.Get("Authorization")
	if strings.HasPrefix(token, "Bearer ") && len(token) > 7 {
		return token[7:], true
	}
	return "", false
}
