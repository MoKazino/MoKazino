package main
import (
    "net/http"
    "net/url"
    "golang.org/x/net/proxy"
    "net/http/cookiejar"
    "context"
    "os"
    "net"
    "time"
    "encoding/json"
    "io/ioutil"
    "strings"
    "log"
    "fmt"
    "strconv"
)
var client http.Client
var jar *cookiejar.Jar
func loadJar() {
    var err error
    jar, err = cookiejar.New(nil)
    if err != nil {
        log.Fatal(err)
    }
    if (os.Getenv("SOCKS5_SERVER") != "") {
        dialer, err := proxy.SOCKS5("tcp", os.Getenv("SOCKS5_SERVER"), nil, proxy.Direct)
        if err != nil {
            log.Fatal(err)
        }
        dialContext := func(ctx context.Context, network, address string) (net.Conn, error) {
            return dialer.Dial(network, address)
        }
        client = http.Client{
            Transport: &http.Transport{
                Proxy:                 http.ProxyFromEnvironment,
                DialContext:                dialContext,
            },
            Jar: jar,
        }
    } else {
        client = http.Client{
            Jar: jar,
        }
    }
}

type IsOK struct {
    OK bool `json:"ok"`
}


func isok(a []byte) bool {
    var x IsOK
    json.Unmarshal(a, &x)
    return x.OK
}
// var cookie Cookie
func checkLogin(username, password string) bool {
    data := url.Values{}
    data.Set("username", username)
    data.Set("password", password)
    resp, err := client.Post("https://mokazino.net/api/v1/login/json", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
    if err != nil {
        log.Fatal(err)
    }
    u, _ := url.Parse("https://mokazino.net/")
    jar.SetCookies(u, resp.Cookies())
    
    defer resp.Body.Close()
    b, _ := ioutil.ReadAll(resp.Body)
    return isok(b)
}

func sendmessage(text string) {
    data := url.Values{}
    data.Set("text", text)
    resp, _ := client.Post("https://mokazino.net/api/v1/chat/send?from=ssh", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
    defer resp.Body.Close()
}

type Bet struct {
    Nonce int `json:"nonce"` // :26
    Currency string `json:"currency"` // :"XMR"
    Balance uint64 `json:"balance"` // :1325869732,
    Profit int64 `json:"profit"`  // :0,
    Roll int `json:"roll"` // :8378
    Win bool `json:"win"` // :true                                                                                                               │
}
func XMRToDecimal(xmr uint64) string {
	str0 := fmt.Sprintf("%013d", xmr)
	l := len(str0)
	return str0[:l-12] + "." + str0[l-12:]
}
func placebet(amount string, chance int64, betlo bool) string {
    data := url.Values{}
    real_amount, err := StringToXMR(amount)
    if err != nil {
        return err.Error()
    }
    data.Set("amount", fmt.Sprint(real_amount))
    data.Set("chance", fmt.Sprint(chance))
    data.Set("betlo", fmt.Sprint(betlo))
    resp, _ := client.Post("https://mokazino.net/api/v1/bet?from=ssh", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
    defer resp.Body.Close()
    b, _ := ioutil.ReadAll(resp.Body)
    var x Bet
    err = json.Unmarshal(b, &x)
    if err != nil {
        return string(b)
    }
    var wonlost = "[green]won[white] "
    if !x.Win {
        wonlost = "[red]lost[white] -"
    }
    var profit string 
    if x.Profit < 0 {
        profit = XMRToDecimal(uint64(x.Profit*-1))
    } else {
        profit = XMRToDecimal(uint64(x.Profit))
    }
    txt := "You have "+wonlost+profit+" [orange]XMR[white]. Roll: [blue]"+profit+"[white] Balance: [yellow]"+XMRToDecimal(x.Balance)+"[white] [orange]XMR[white]"
    return txt
}
func withdraw(amount uint64, address string) string {
    data := url.Values{}
    data.Set("amount", fmt.Sprint(amount))
    data.Set("address", address)
    resp, _ := client.Post("https://mokazino.net/api/v1/withdraw/xmr", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
    defer resp.Body.Close()
    b, _ := ioutil.ReadAll(resp.Body)
    return string(b)
}
type Profile struct {
	User User  `json:"user"`
	// Bets []Bet `json:"bets"`
}
type User struct {
	ID                   int    `json:"id"`
	Username             string `json:"username"`
	Cookie               string `json:"cookie"`
	Currency             string `json:"currency"`
	Balance              int64  `json:"balance"`
	BitcoinBalance       int64  `json:"bitcoinbalance"`
	Wagered              uint64 `json:"wagered"`
	WageredBitcoin       uint64 `json:"wageredbitcoin"`
	InvestBalance        uint64 `json:"investbalance"`
	BitcoinInvestBalance uint64 `json:"investbalancebitcoin"`
	MoneroAddress        string `json:"moneroaddress"`
	BitcoinAddress       string `json:"bitcoinaddress"`
	Nonce                uint   `json:"nonce"`
	ChatBanned           bool   `json:"chatbanned"`
	ClientSeed           string `json:"clientseed"`
}
func getProfile() (x Profile) {
    resp, _ := client.Get("https://mokazino.net/api/v1/profile")
    defer resp.Body.Close()
    b, _ := ioutil.ReadAll(resp.Body)
    json.Unmarshal(b, &x)
    return x
}

func StringToXMR(xmr string) (uint64, error) {
	f, err := strconv.ParseFloat(xmr, 64)
	if err != nil {
		return 0, err
	}
	return uint64(f * 1e12), nil
}
func getmessage() string {
    resp, err := client.Get("https://mokazino.net/api/v1/chat/read")
    if err != nil {
        time.Sleep(time.Second * 2)
        return getmessage()
    }
    defer resp.Body.Close()
    b, _ := ioutil.ReadAll(resp.Body)
    if err != nil {
        time.Sleep(time.Second * 2)
        return getmessage()
    }
    return string(b)
}