package utilPackage

import (
	"time"
	"encoding/json"
	"net/http"
	"github.com/gorilla/pat"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// 1. 먼저 전역변수로 render
var rd *render.Render

type User struct {
	Name string `json:"name"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name :"Jod", Email:"jod@gmail.com"}
	
	// 아래와 같은 뜻
	rd.JSON(w, http.StatusOK, user)
	// w.Header().Add("Content-type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// data, _ := json.Marshal(user)
	// fmt.Fprint(w, string(data))
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)

	// 아래와 같은 뜻
	rd.JSON(w, http.StatusOK, user)
	 err := json.NewDecoder(r.Body).Decode(user)
	 if err != nil {
		rd.Text(w, http.StatusBadRequest, err.Error())
	 	// w.WriteHeader(http.StatusBadRequest)
		//  fmt.Fprint(w, err)
	 	return
	 }
	 user.CreatedAt = time.Now()
	rd.JSON(w, http.StatusOK, user)
}

// 템플릿 활용
func helloHandler(w http.ResponseWriter, r *http.Request) {
	user := User{Name :"Jod", Email:"jod@gmail.com"}
	// 아래와 같은 뜻
	// 주의 ) 확장자를 빼고 등록하자
	rd.HTML(w, http.StatusOK, "body", user)
	// tmpl, err := template.New("Hello").ParseFiles("utilPackage/template/hello.tmpl")
	// if err != nil {
	// 	rd.Text(w, http.StatusInternalServerError, err.Error())
		// w.WriteHeader(http.StatusInternalServerError)
		// fmt.Fprint(w, err)
	// 	return
	// }
	// tmpl.ExecuteTemplate(w, "hello.tmpl", "jod")
}

func UtilPackage() {
	// 2. render 객체 생성
	rd = render.New(render.Options{
		Directory: "utilPackage/template", // 기본값 : templates
		Extensions: []string{".html",".tmpl"},
		Layout: "hello"})
		
	// pat 패키지 추가하면 더 편하게 라우터를 추가할 수 있다.
	mux := pat.New()
	
	// pat을 사용하면 다음과 같이 표현
	mux.Get("/users", getUserInfoHandler)
	mux.Post("/users", addUserHandler)
	mux.Get("/hello", helloHandler)

	// 기본 파일서버
	n := negroni.Classic()
	n.UseHandler(mux)

	http.ListenAndServe(":3000", n)
}