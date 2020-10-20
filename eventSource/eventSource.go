package eventSource

import (
	"time"
	"strconv"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/urfave/negroni"
	"github.com/gorilla/pat"
	"github.com/antage/eventsource"
)

// 과거 -> 현재
// 과거 : request에 대한 response를 보냄
// 현재 : 동적인 요구가 강해짐 -> html5에서 web socket, event source가 추가됨
// web socket : tcp socket처럼 연결해서 send/recv 하겠다.
// event source : server sent event만 가능 일방향

// evnet source : 구독한 클라이언트에 대한 푸시알림, 이벤트 알림 등에 활용



func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	msg := r.FormValue("msg")
	name := r.FormValue("name")
	fmt.Println(msg, name)
	sendMessage(name, msg)
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("name")
	sendMessage("", fmt.Sprintf("add user: %s", username))
}

func leftUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	sendMessage("", fmt.Sprintf("left user: %s", username))
}

type Message struct {
	Name string `json:"name"`
	Msg string `json:"msg"`
}

var msgCh chan Message

func sendMessage(name, msg string) {
	// send message to every clients
	msgCh <- Message{name, msg}
}

func processMsgCh(es eventsource.EventSource) {
	for msg := range msgCh {
		data, _ := json.Marshal(msg)
		es.SendEventMessage(string(data), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}

func EventSource() {
	// 채널 초기화
	msgCh = make(chan Message)

	// 이벤트 소스 소켓 생성
	es := eventsource.New(nil, nil)
	defer es.Close()

	go processMsgCh(es)

	mux := pat.New()
	mux.Post("/messages", postMessageHandler)
	mux.Handle("/stream", es)
	mux.Post("/users", addUserHandler)
	mux.Delete("/users", leftUserHandler)

	n:=negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}