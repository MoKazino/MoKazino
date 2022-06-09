package web

import (
	"net/http"
	"strings"
)

var captchas = make(map[string]string)
var timec = make(map[string]int)

func checkCaptcha(id, input string) bool {
	if strings.EqualFold(captchas[id], input) && id != "" && captchas[id] != "" {
		delete(captchas, id)
		delete(timec, id)
		return true
	}
	return false
}

func checkCaptchaCookie(r *http.Request, input string) bool {
	cookie, err := r.Cookie("rubot")
	if err != nil {
		return false
	}
	return checkCaptcha(cookie.Value, input)
}
