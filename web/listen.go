package web

import (
	"embed"
	"log"
	"net/http"
	"strings"
)

//go:embed node_modules/bootstrap/dist/css/*.css
//go:embed node_modules/bootstrap/dist/js/*.js
//go:embed index.html
//go:embed img/*.svg
//go:embed static/*.html
//go:embed static/*.js
//go:embed provably_fair.go
//go:embed message.txt mokazino.gpg docs.txt
var files embed.FS

func init() {
	log.SetFlags(log.Lshortfile)
}

func Listen() {
	http.HandleFunc("/captcha", ShowCaptcha)
	http.HandleFunc("/captcha_cookie", ShowCaptchaCookie)
	http.HandleFunc("/api/v1/switchcurrency/btc", apiv1switchcurrencybtc)
	http.HandleFunc("/api/v1/switchcurrency/btc/json", apiv1switchcurrencybtcjson)
	http.HandleFunc("/api/v1/switchcurrency/xmr", apiv1switchcurrencyxmr)
	http.HandleFunc("/api/v1/switchcurrency/xmr/json", apiv1switchcurrencyxmrjson)
	http.HandleFunc("/api/v1/register", apiv1register)
	http.HandleFunc("/api/v1/register/json", apiv1registerjson)
	http.HandleFunc("/api/v1/register/oneclick", apiv1registeroneclick)
	http.HandleFunc("/api/v1/login", apiv1login)
	http.HandleFunc("/api/v1/login/json", apiv1loginjson)
	http.HandleFunc("/api/v1/profile", apiv1profile)
	http.HandleFunc("/api/v1/depositqr", apiv1depositqr)
	http.HandleFunc("/api/v1/depositqr/custom", apiv1depositqrcustom)
	http.HandleFunc("/api/v1/profile_update", apiv1profileupdate)
	http.HandleFunc("/api/v1/number", apiv1number)
	http.HandleFunc("/api/v1/bet", apiv1bet)
	http.HandleFunc("/api/v1/withdraw", apiv1withdraw)
	http.HandleFunc("/api/v1/withdraw/xmr", apiv1withdrawxmr)
	http.HandleFunc("/api/v1/withdraw/btc", apiv1withdrawbtc)
	http.HandleFunc("/api/v1/invest", apiv1invest)
	http.HandleFunc("/api/v1/closeinvest", apiv1closeinvest)
	http.HandleFunc("/api/v1/reset_seed", apiv1resetseed)
	http.HandleFunc("/api/v1/reset_seed/all", apiv1resetseedall)
	http.HandleFunc("/api/v1/serverstats", apiv1serverstats)
	http.HandleFunc("/api/v1/chat/send", apiv1chatsend)
	http.HandleFunc("/api/v1/chat/read", apiv1chatread)
	http.HandleFunc("/api/v1/external/majesticbank/exchange", apiv1externalmajesticbankexchnage)
	http.HandleFunc("/api/guest/v1/stats", apiguestv1stats)
	http.HandleFunc("/verify", func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("rubot")
		if !IsBot(rw, r) {
			return
		}
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
			return
		}
		if checkCaptcha(cookie.Value, r.PostFormValue("captcha")) {
			allowed_cookies[cookie.Value] = true
			rw.Header().Add("Location", r.URL.RawPath)
			rw.Header().Add("Content-Type", "text/html")
			rw.WriteHeader(300)
			rw.Write([]byte(`<meta http-equiv="Refresh" content="0; url='` + rw.Header().Get("Location") + `'" />`))
		} else {
			IsBot(rw, r)
		}
	})
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.SetCookie(rw, &http.Cookie{
			Name:     "ref",
			Value:    r.URL.Query().Get("ref"),
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
			HttpOnly: true,
		})
		//log.Println(r.URL.Path)
		if IsBot(rw, r) {
			return
		}
		if r.URL.Path == "/" || len(r.URL.Path) < 1 {
			r.URL.Path = "/index.html"
		}
		if strings.HasPrefix(r.URL.Path, "/static/home") {
			r.URL.Path = "/static/homev2.html"
		}
		f, err := files.ReadFile(r.URL.Path[1:])
		if err != nil {
			rw.WriteHeader(500)
			rw.Write([]byte(err.Error()))
			return
		}
		rw.Header().Add("Content-Type", getCT(r.URL.Path))
		rw.Write(f)
	})
	log.Println("Listening on :2132")
	log.Println(http.ListenAndServe(":2132", nil))
}

func getCT(path string) string {
	if strings.HasSuffix(path, ".js") {
		return "text/javascript"
	} else if strings.HasSuffix(path, ".css") {
		return "text/css"
	} else if strings.HasSuffix(path, ".html") {
		return "text/html"
	} else if strings.HasSuffix(path, ".ico") {
		return "image/x-icon"
	} else if strings.HasSuffix(path, ".txt") || strings.HasSuffix(path, ".go") {
		return "text/plain"
	} else if strings.HasSuffix(path, ".html") {
		return "text/html"
	} else if strings.HasSuffix(path, ".jpg") || strings.HasSuffix(path, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(path, ".png") {
		return "image.png"
	} else if strings.HasSuffix(path, ".svg") {
		return "image/svg+xml"
	}
	return "application/octet-stream"
}
