package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// 유저 정보를 가지고 있는 JSON
// request body의 json변수와 다를 때를 위해, 어노테이션도 붙일 수 있다.
type User struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type fooHandler struct{}

// 인터페이스 형태로 사용하기 위해 ServeHTTP(고정) 재정의
func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(User)

	// body가 user struct 구조에 맞는지 파악
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user.CreatedAt = time.Now()

	// json 형태로 바꿈
	data, _ := json.Marshal(user)
	// response의 형태를 알려준다.
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, string(data))
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

func barHandler(w http.ResponseWriter, r *http.Request) {
	// query 가져오기
	name := r.URL.Query().Get("name")

	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s", name)
}

func NewHttpHandler() http.Handler {
	// 라우터 인스턴스 만들기
	mux := http.NewServeMux()

	// request에 대한 핸들러를 등록한다.
	// 등록방법1
	mux.HandleFunc("/", indexHandler)
	// 등록방법2
	mux.HandleFunc("/bar", barHandler)

	// 등록방법3
	// 인터페이스를 구현한 애를 등록하는 방법
	mux.Handle("/foo", &fooHandler{})

	return mux
}
