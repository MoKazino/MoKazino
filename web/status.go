package web

import (

	"github.com/btcsuite/btcutil"

	"fmt"
	"net/http"
)

func apiguestv1stats(rw http.ResponseWriter, r *http.Request) {
	h, errh := w.GetHeight()
	bh, errbh := btcclient.GetBlockChainInfo()
	rw.Write(errorit(`
Monero Height: ` + fmt.Sprint(h.Height) + perr(errh) + `<br />
Bitcoin Height: ` + fmt.Sprint(bh.Blocks) + ` (` + fmt.Sprintf("%.2f", bh.VerificationProgress*100) + `%) ` + perr(errbh) + `
<hr />
All money XMR = (`+fmt.Sprint(Printxmr(uint64(userbalance)))+` XMR + `+fmt.Sprint(Printxmr(uint64(serverbalance)))+` XMR)<br />
All money BTC = (`+fmt.Sprint(btcutil.Amount(userbalance_btc).String())+` BTC + `+fmt.Sprint(btcutil.Amount(serverbalance_btc).String())+`)`))
}

func perr(err error) string {
	if err != nil {
		return "(" + err.Error() + ")"
	}
	return ""
}
