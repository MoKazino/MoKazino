package web

import (
	"encoding/json"
	"net/http"
)

type Stats struct {
	MaxProfit     string `json:"maxprofit"`
	UserBalance   string `json:"userbalance"`
	ServerBalance string `json:"serverbalance"`
}

func apiv1serverstats(rw http.ResponseWriter, r *http.Request) {
	auth, u := IsAuth(rw, r)
	if !auth {
		return
	}
	var b []byte
	if u.Currency == "BTC" {
		b, _ = json.Marshal(Stats{
			MaxProfit:     u.CurToString(maxprofit_btc),
			UserBalance:   u.CurToString(uint64(userbalance_btc)),
			ServerBalance: u.CurToString(uint64(serverbalance_btc)),
		})
	} else {
		b, _ = json.Marshal(Stats{
			MaxProfit:     u.CurToString(maxprofit),
			UserBalance:   u.CurToString(uint64(userbalance)),
			ServerBalance: u.CurToString(uint64(serverbalance)),
		})
	}
	rw.Header().Add("Content-Type", "application/json")
	rw.Write(b)
}
