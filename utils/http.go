package utils

import (
	"strings"
)

func GetToken(token string) (string, bool) {
	if strings.HasPrefix(token, "Bearer ") && len(token) > 7 {
		return token[7:], true
	}
	return "", false
}
