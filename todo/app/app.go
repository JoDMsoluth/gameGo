package app

import (
	"strings"
	"github.com/gorilla/sessions"
	"github.com/JoDMsoluth/webGo/todo/model"
	"strconv"
	"net/http"
	"os"
	"github.com/unrolled/render"
	"github.com/gorilla/mux"
	
	"github.com/urfave/negroni"
)

var rd *render.Render = render.New()
// 쿠키 저장소를 만든다.
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type AppHandler struct {
	http.Handler // 그냥 이름은 안정해도 암싲거으로 포함은 한다는 뜻
	db model.DBHandler
}

var getSessionID = func (r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}

	// Set some session values.
	val := session.Values["id"]
	
	// session id가 비어있다.
	if val == nil {
		return ""
	}
	
	return val.(string)
}

func (a *AppHandler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (a *AppHandler) getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	// list := []*model.Todo{}
	// for _, v := range todoMap {
	// 	list = append(list, v)
	// }
	sessionId := getSessionID(r)
	list := a.db.GetTodo(sessionId)
	rd.JSON(w, http.StatusOK, list)
}

func (a *AppHandler) addTodoListHandler(w http.ResponseWriter, r *http.Request) {
	sessionId := getSessionID(r)
	name := r.FormValue("name")
	todo := a.db.AddTodo(name, sessionId)
	// id := len(todoMap) 
	// todo := &a.db.Todo{id, name, false, time.Now()}
	// todoMap[id] = todo
	rd.JSON(w, http.StatusCreated, todo)
}

type Success struct {
	Success bool `json:"success"`
}

func (a *AppHandler) removeTodoListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := a.db.RemoveTodo(id)
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {

	}
	// if _, ok := todoMap[id]; ok {
	// 		delete(todoMap, id)
	// 		rd.JSON(w, http.StatusOK, Success{true})
	// 	}	else {
	// 		rd.JSON(w, http.StatusOK, Success{false})
	// }
}

func (a *AppHandler) completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"

	ok := a.db.CompleteTodo(id, complete)
	
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusBadRequest, Success{false})
	} 
}

func (a *AppHandler) Close() {
	a.db.Close()
}

func CheckSignin(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// 만약 이미 signin 페이지에서의 요청이었을 때는 적용 x
	if strings.Contains(r.URL.Path, "/signin") || strings.Contains(r.URL.Path, "/auth") {
		next(w, r)
		return
	}

	// if user already signed in
	sessionID := getSessionID(r)
	if sessionID != "" {
		next(w, r)
		return
	}
	
	// if not user sign in
	// redirect signin.html
	http.Redirect(w, r, "/signin.html", http.StatusTemporaryRedirect)
}

func MakeNewHandler(filepath string) *AppHandler {
	r := mux.NewRouter()
	// recovery : panic이 일어난 다음 서버가 안꺼지게 도와줌
	// negroni의 미들웨어 추가
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.HandlerFunc(CheckSignin),  negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	// app handler로 인스턴스를 만들어서 거기에 db와 hander를 포함시켜서 사용한다.
	// 이렇게 하면 db핸들러가 사용한 쪽에서 Close를 사용해야하는데 app에서 사용하기 때문에 거기에 db핸들러를 포함시키면 close를 실행권한을 가지게 된다.
	a := &AppHandler{
		Handler : n,
		db: model.NewDBHandler(filepath),
	}

	
	r.HandleFunc("/auth/google/login", googleLoginHandler)
	r.HandleFunc("/auth/google/callback", googleAuthCallback)
	r.HandleFunc("/todos", a.getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", a.addTodoListHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", a.removeTodoListHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", a.completeTodoHandler).Methods("GET")
	r.HandleFunc("/", a.indexHandler)

	return a
}