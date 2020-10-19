package myapp

// 메인앱 이 파일을 수정하지 않고 데코레이터만을 이용해서 확장을 한다.

import (
	"fmt"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	return mux
}
