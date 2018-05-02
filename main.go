package chatroot

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"crypto/md5"
	"strconv"
	"encoding/json"
)

var upgrader = websocket.Upgrader{}

const TokenName = "token"

func Main() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/ws", wsHandel)
	http.HandleFunc("/login", loginHandel)
	http.HandleFunc("/user/get", getUcInfo)
	http.ListenAndServe(":3000", nil)
}

func getUcInfo(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()
	name := r.FormValue("name")
	if name == "" {
		AjaxReturn(w, 3)
		return
	}
	account ,err := getAccountByName(name)
	if err != nil {
		AjaxReturn(w, 3)
		return
	}
	m := NewMessage(0, strconv.Itoa(account.id))
	ms, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(ms)
}

func loginHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		AjaxReturn(w, 2)
		return
	}
	user := r.PostFormValue("user")
	password := r.PostFormValue("password")
	for _, v := range Accounts {
		if v.user == user && v.password == password {
			salt := md5.New()
			salt.Write([]byte(user))
			readyRandom := fmt.Sprintf("%x", salt.Sum(nil))
			MyCookies.Set(readyRandom, v.id)
			cookie := http.Cookie{
				Name:TokenName,
				Value:readyRandom,
			}
			http.SetCookie(w, &cookie)
			AjaxReturn(w, 0)
			return
		}
	}
	AjaxReturn(w, 1)
}

func indexHandle(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(TokenName)
	if err != nil {
		http.Redirect(w, r, "/public/login.html", http.StatusFound)
		return
	}
	if  MyCookies.Get(token.Value) == 0 {
		fmt.Println("用户未登录1")
		http.Redirect(w, r, "/public/login.html", http.StatusFound)
	} else {
		http.ServeFile(w, r, "public/index.html")
	}
}

func wsHandel(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie(TokenName)
	if err != nil {
		fmt.Println("token不存在: ", err)
		http.Redirect(w, r, "/public/login.html", http.StatusFound)
		return
	}
	if  id := MyCookies.Get(token.Value); id == 0 {
		fmt.Println("用户未登录2")
		http.Redirect(w, r, "/public/login.html", http.StatusFound)
	} else {
		var client *Client
		if client = getClient(id); client == nil {
			if account, err := getAccount(id); err != nil {
				fmt.Println(err)
				return
			} else {
				c, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					fmt.Println("upgrader错误: ", err)
					return
				}
				account.SetCookie(token.Value)
				client = NewUser(account)
				client.conn = c
				client.addUser()
			}
		}
		client.readMes()
	}
}