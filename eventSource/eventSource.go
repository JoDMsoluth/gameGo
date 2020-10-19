package eventSource

import (
	"log"
	"net/http"
	"github.com/urfave/negroni"
	"github.com/gorilla/pat"
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
	log.Println("postMessageHandler ", msg, name)
}
func EventSource() {
	mux := pat.New()
	mux.Post("/messages", postMessageHandler)

	n:=negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}