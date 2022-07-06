package web

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var briarcookie = ""

func init() {
	log.Println("init(): Loading briar")
	b, err := ioutil.ReadFile("/tmp/briar_cookie")
	if err != nil {
		log.Fatalln(err)
	}
	briarcookie = string(b)
	// go briarProxyMessages()
}

type BriarLink struct {
	Link string `json:"link"`
}

var briar_auth string
var briar_host = "http://127.0.0.1:7000"

func LoadBriar() {
	data, err := ioutil.ReadFile("/tmp/briar_cookie")
	if err != nil {
		log.Println("LoadBriar:", err)
		return
	}
	briar_auth = string(data)
}

func init() {
	go LoadBriar()
	go func() {
		for {
			time.Sleep(time.Second * 5)
			BriarPool()
		}
	}()
	go func() {
		log.Println("init(): briarProxyMessages")
		for msg := range messages.Subscribe() {
			var users []User
			db.Find(&users, "briar_id != ?", 0)
			for i := range users {
				go BriarSendMessage(users[i].BriarID, msg)
			}
		}
	}()
}

func BriarPool() {
	ctcs := GetAllContacts()
	for i := range ctcs {
		if ctcs[i].UnreadCount == 0 {
			break
		}
		msgs := BriarGetMessages(int(ctcs[i].ContactID))
		for i := range msgs {
			if msgs[i].Type != "PrivateMessage" || msgs[i].Local {
				continue
			}
			BriarSendReq("/v1/messages/"+strconv.Itoa(int(msgs[i].ContactID))+"/all", "DELETE", nil)

			if msgs[i].Text == nil {
				continue
			}
			txt := msgs[i].Text.(string)

			spltxt := strings.Split(txt, " ")
			if spltxt[0] == "login" {
				if len(spltxt) != 3 {
					BriarSendMessage(int(msgs[i].ContactID), "Usage: login username password")
				}
				username := spltxt[1]
				password := spltxt[2]
				var u User
				db.First(&u, "username = ?", username)
				if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
					BriarSendMessage(int(msgs[i].ContactID), "Invalid username/password combination")
					return
				}
				u.BriarID = int(msgs[i].ContactID)
				db.Save(&u)
				BriarSendMessage(int(msgs[i].ContactID), "Account linked!")
			} else {
				var u User
				db.First(&u, "briar_id = ?", msgs[i].ContactID)
				if u.BriarID != int(msgs[i].ContactID) {
					BriarSendMessage(int(msgs[i].ContactID), "um- sir your chat account is not connected. Please type 'login username password' to link your account.")
					return
				}
				messages.Publish(u.Username + "(briar): " + txt)
				continue
			}
		}
	}
}

type BriarV1ContctsAddLink struct {
	Link string `json:"link"`
}

func BriarGetLink() string {
	var resp BriarV1ContctsAddLink
	r := BriarSendReq("/v1/contacts/add/link", "GET", nil)
	json.Unmarshal(r, &resp)
	return resp.Link
}

type BriarV1Messages struct {
	ContactID int64       `json:"contactId"`
	GroupID   string      `json:"groupId"`
	ID        string      `json:"id"`
	Local     bool        `json:"local"`
	Read      bool        `json:"read"`
	Seen      bool        `json:"seen"`
	Sent      bool        `json:"sent"`
	Text      interface{} `json:"text"`
	Timestamp int64       `json:"timestamp"`
	Type      string      `json:"type"`
}

func BriarGetMessages(id int) []BriarV1Messages {
	var resp []BriarV1Messages
	r := BriarSendReq("/v1/messages/"+strconv.Itoa(id), "GET", nil)
	err := json.Unmarshal(r, &resp)
	if err != nil {
		log.Println(err)
	}
	// BriarSendReq("/v1/messages/"+strconv.Itoa(id)+"/all", "DELETE", nil)
	return resp
}

type BriarV1MessagePost struct {
	Text string `json:"text"`
}

func BriarSendMessage(id int, text string) {
	log.Println(id, text)
	log.Println(string(BriarSendReq("/v1/messages/"+strconv.Itoa(id), "POST", BriarV1MessagePost{
		Text: text,
	})))
}

type BriarV1Contacts struct {
	Alias  string `json:"alias"`
	Author struct {
		FormatVersion int64  `json:"formatVersion"`
		ID            string `json:"id"`
		Name          string `json:"name"`
		PublicKey     string `json:"publicKey"`
	} `json:"author"`
	Connected          bool   `json:"connected"`
	ContactID          int64  `json:"contactId"`
	HandshakePublicKey string `json:"handshakePublicKey"`
	LastChatActivity   int64  `json:"lastChatActivity"`
	UnreadCount        int64  `json:"unreadCount"`
	Verified           bool   `json:"verified"`
}

func GetAllContacts() []BriarV1Contacts {
	var resp []BriarV1Contacts
	r := BriarSendReq("/v1/contacts", "GET", nil)
	json.Unmarshal(r, &resp)
	return resp
}

// GET /v1/contacts

func briaradd(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Cache-Control", "no-store")
	auth, _ := IsAuth(rw, r)
	if !auth {
		return
	}
	if r.PostFormValue("link") != "" {
		BriarAddByLink(r.PostFormValue("link"), RandomString(20))
		rw.Write(errorit("Request sent! Make sure to add my briar too! <code>" + BriarGetLink() + "</code>"))
		return
	}
	rw.Write(errorit(`Link: <code>` + BriarGetLink() + `</code>
<br />
Below you need to type your address so I'll be able to add you
<form method="post" action="/api/v1/briaradd">
	<div class="row gtr-uniform">
		<div class="col-6 col-12-xsmall">
			<input type="text" name="link" id="link" placeholder="Your briar link." />
		</div>
		<!-- Break -->
		<div class="col-12">
			<ul class="actions">
				<li><input type="submit" value="Submit" class="primary" /></li>
			</ul>
		</div>
	</div>
</form>`))
}

type BriarPostV1ContactsAddPending struct {
	Link  string `json:"link"`
	Alias string `json:"alias"`
}

func BriarAddByLink(link string, alias string) {
	BriarSendReq("/v1/contacts/add/pending", "POST", BriarPostV1ContactsAddPending{
		Link:  link,
		Alias: alias,
	})
}

func BriarSendReq(endpoint string, method string, data interface{}) []byte {
	client := &http.Client{}
	var b []byte
	if method == "GET" {
	} else {
		b, _ = json.Marshal(data)
	}

	req, err := http.NewRequest(method, briar_host+endpoint, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+briar_auth)
	if err != nil {
		log.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	// log.Println(string(b))
	return b
}
