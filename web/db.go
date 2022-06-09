package web

import (
	"log"
	"strings"

	"github.com/btcsuite/btcutil"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID                   int    `json:"id"`
	Username             string `gorm:"uniqueIndex" json:"username"`
	Password             string `json:"-"`
	Cookie               string `json:"cookie"`
	Currency             string `json:"currency"`
	Balance              int64  `json:"balance"`
	BitcoinBalance       int64  `json:"bitcoinbalance"`
	Wagered              uint64 `json:"wagered"`
	WageredBitcoin       uint64 `json:"wageredbitcoin"`
	InvestBalance        uint64 `json:"investbalance"`
	BitcoinInvestBalance uint64 `json:"investbalancebitcoin"`
	AddressIndex         uint64 `json:"-"`
	MoneroAddress        string `json:"moneroaddress"`
	BitcoinAddress       string `json:"bitcoinaddress"`
	Nonce                uint   `json:"nonce"`
	ChatBanned           bool   `json:"chatbanned"`
	ClientSeed           string `json:"clientseed"`
	ServerSeed           string `json:"-"`
	Inviter              int    `json:"-"`
}

func (u *User) AddWagered(x uint64) {
	if u.Currency == "BTC" {
		u.WageredBitcoin += x
	}
	u.Wagered += x
}

func (u *User) GetWagered() uint64 {
	if u.Currency == "BTC" {
		return u.WageredBitcoin
	}
	return u.Wagered
}

func (u *User) AddBalance(x int64) {
	if u.Currency == "BTC" {
		u.BitcoinBalance += x
	}
	u.Balance += x
}

func (u *User) AddInvestBalance(x int64) {
	if u.Currency == "BTC" {
		u.BitcoinInvestBalance += uint64(x)
	}
	u.InvestBalance += uint64(x)
}

func (u *User) GetMaxProfit() uint64 {
	if u.Currency == "BTC" {
		return maxprofit_btc
	}
	return maxprofit
}

func (u *User) CurToString(cur uint64) string {
	if u.Currency == "BTC" {
		return strings.ReplaceAll(btcutil.Amount(cur).String(), " BTC", "")
	}
	return Printxmr(cur)
}

func (u *User) GetCurrencyFull() string {
	if u.Currency == "BTC" {
		return "Bitcoin"
	} else {
		return "Monero"
	}
}

func (u *User) GetDepositAddress() string {
	if u.Currency == "BTC" {
		return u.BitcoinAddress
	} else {
		return u.MoneroAddress
	}
}

func (u *User) GetBalance() int64 {
	if u.Currency == "BTC" {
		return u.BitcoinBalance
	} else {
		return u.Balance
	}
}

func (u *User) GetInvestBalance() uint64 {
	if u.Currency == "BTC" {
		return u.BitcoinInvestBalance
	}
	return u.InvestBalance
}

func (u *User) ResetInvestBalance() {
	if u.Currency == "BTC" {
		u.BitcoinInvestBalance = 0
		return
	}
	u.InvestBalance = 0
}

type Transaction struct {
	TXID string
}
type Bet struct {
	ID         uint64
	Nonce      uint64
	Amount     uint64
	Currency   string
	ClientSeed string
	ServerSeed string `json:"-"`
	User       int    `json:"-"`
}

func init() {
	log.Println("Loading db")
	var err error
	db, err = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("User:", db.AutoMigrate(&User{}))
	log.Println("Transaction:", db.AutoMigrate(&Transaction{}))
	log.Println("Bet:", db.AutoMigrate(&Bet{}))
}
