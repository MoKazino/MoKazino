package web

import (
	"log"
	"net/http"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

func apiv1depositqr(rw http.ResponseWriter, r *http.Request) {
	auth, u := IsAuth(rw, r)
	png, err := qrcode.Encode("https://youtu.be/dQw4w9WgXcQ?t=42", qrcode.Medium, 512)
	if auth {
		png, err = qrcode.Encode(strings.ToLower(u.GetCurrencyFull())+":"+u.GetDepositAddress(), qrcode.Medium, 512)
	}
	if err != nil {
		log.Println(err)
		return
	}
	rw.Header().Add("Content-Type", "image/png")
	rw.Write(png)

}

func apiv1depositqrcustom(rw http.ResponseWriter, r *http.Request) {
	auth, _ := IsAuth(rw, r)
	png, err := qrcode.Encode("https://youtu.be/dQw4w9WgXcQ?t=42", qrcode.Medium, 512)
	if auth {
		png, err = qrcode.Encode(r.URL.Query().Get("qr"), qrcode.Medium, 512)
	}
	if err != nil {
		log.Println(err)
		return
	}
	rw.Header().Add("Content-Type", "image/png")
	rw.Write(png)

}
