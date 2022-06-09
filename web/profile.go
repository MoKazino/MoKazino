package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/btcsuite/btcutil"
	"github.com/monero-ecosystem/go-monero-rpc-client/wallet"
)

type ProfileResponse struct {
	User User  `json:"user"`
	Bets []Bet `json:"bets"`
}

func apiv1switchcurrencyxmr(rw http.ResponseWriter, r *http.Request) {
	auth, u := IsAuth(rw, r)
	if !auth {
		return
	}
	u.Currency = "XMR"
	db.Save(&u)
	rw.Write(errorit(`<a href="/static/login.html">Login</a><meta http-equiv="Refresh" content="0; url='/static/home/` + RandomString(32) + `'" />`))
}

func apiv1switchcurrencybtc(rw http.ResponseWriter, r *http.Request) {
	auth, u := IsAuth(rw, r)
	if !auth {
		return
	}
	u.Currency = "BTC"
	db.Save(&u)
	rw.Write(errorit(`<a href="/static/login.html">Login</a><meta http-equiv="Refresh" content="0; url='/static/home/` + RandomString(32) + `'" />`))
}

func apiv1switchcurrencyxmrjson(rw http.ResponseWriter, r *http.Request) {
	auth, u := IsAuth(rw, r)
	if !auth {
		return
	}
	u.Currency = "XMR"
	db.Save(&u)
	rw.Write(okitjson(`ok`))
}

func apiv1switchcurrencybtcjson(rw http.ResponseWriter, r *http.Request) {
	auth, u := IsAuth(rw, r)
	if !auth {
		return
	}
	u.Currency = "BTC"
	db.Save(&u)
	rw.Write(okitjson(`ok`))
}


func apiv1profile(rw http.ResponseWriter, r *http.Request) {
	auth, u := IsAuth(rw, r)
	if !auth {
		return
	}
	if u.AddressIndex == 0 {
		r, err := w.CreateAddress(&wallet.RequestCreateAddress{})
		if err != nil {
			//rw.WriteHeader(500)
			//rw.Write(errorit(err.Error()))
			//return
			log.Println(err)
		} else {
			u.AddressIndex = r.AddressIndex
			log.Println("new address", r.Address, r.AddressIndex)
			db.Save(&u)
		}
	}
	mr, err := w.GetAddress(&wallet.RequestGetAddress{
		AccountIndex: 0,
		AddressIndex: []uint64{u.AddressIndex},
	})
	if err != nil {
		log.Println(err)
		//rw.WriteHeader(500)
		//rw.Write(errorit(err.Error()))
		//return
	} else {
		u.MoneroAddress = mr.Addresses[0].Address
		db.Save(&u)
	}
	amt := checkDeposit(u.AddressIndex)
	if amt != 0 {
		u.Balance += int64(amt)
		messages.Publish("[wallet] Deposit of " + Printxmr(amt) + " xmr just got credited!")
		db.Save(&u)
	}
	//db.Save(&u)
	var bets []Bet
	if u.BitcoinAddress == "" {
		//u.BitcoinAddress =
		address, err := btcclient.GetNewAddress("user" + strconv.Itoa(u.ID))
		if err != nil {
			log.Println(err)
		} else {
			u.BitcoinAddress = address.EncodeAddress()
			db.Save(&u)
		}
	}
	amt = checkDeposit_btc(u.BitcoinAddress)
	if amt != 0 {
		u.BitcoinBalance += int64(amt)
		messages.Publish("[wallet] Deposit of " + btcutil.Amount(amt).String() + " BTC just got credited!")
		db.Save(&u)
	}
	if u.Currency == "" {
		u.Currency = "XMR"
	}
	db.Save(&u)

	//db.Find(&bets, "user = ?", u.ID).Order("id desc").Limit(50)
	b, _ := json.Marshal(ProfileResponse{
		User: u,
		Bets: bets,
	})
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(b)
}

func apiv1profileupdate(rw http.ResponseWriter, r *http.Request) {
	auth, u := IsAuth(rw, r)
	if !auth {
		return
	}
	u.ClientSeed = r.PostFormValue("clientseed")
	if len(u.ClientSeed) > 32 {
		u.ClientSeed = u.ClientSeed[0:31]
	}
	db.Save(&u)
	rw.Write([]byte(`{"status": "Okay!"}`))
}
