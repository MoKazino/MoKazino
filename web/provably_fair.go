package web

import (
	"crypto/sha512"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
)

// GenerateProvablyFairNumber
//  clientSeed - string provided by client
//  serverSeed - string provided by server (secret)
//  nonce      - how many times a number was generated?
//  max        - what is the maximum expected value? - 9999 on site
func GenerateProvablyFairNumber(clientSeed, serverSeed string, nonce uint, max uint) uint {
	payload := clientSeed + ":" + serverSeed + ":" + strconv.Itoa(int(nonce))
	h := sha512.New()
	h.Write([]byte(payload))
	hash := hex.EncodeToString(h.Sum(nil))

	wantedstring := hash[:8]
	value, err := strconv.ParseInt(wantedstring, 16, 64)
	if err != nil {
		log.Panicln("CRITICAL FAILURE - failed to generate number.", err, wantedstring, hash, payload)
		// Just in case.
		<-make(chan int)
	}
	return uint(value) % (max + 1)
}

// If for some dirty, wild reason you want to test the code above:
// <url>/api/v1/number?client_seed=cs
//       &server_seed=ss
//       &nonce=1
//       &max=100
func apiv1number(rw http.ResponseWriter, r *http.Request) {
	if IsBot(rw, r) {
		return
	}
	clientSeed := r.URL.Query().Get("client_seed")
	serverSeed := r.URL.Query().Get("server_seed")
	nonce, err := strconv.Atoi(r.URL.Query().Get("nonce"))
	if err != nil {
		rw.WriteHeader(500)
		rw.Write(errorit(err.Error()))
	}
	max, err := strconv.Atoi(r.URL.Query().Get("max"))
	if err != nil {
		rw.WriteHeader(500)
		rw.Write(errorit(err.Error()))
	}
	rw.Write([]byte(strconv.Itoa(int(GenerateProvablyFairNumber(clientSeed, serverSeed, uint(nonce), uint(max))))))
}
