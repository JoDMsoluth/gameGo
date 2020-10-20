package app

import (
	"github.com/JoDMsoluth/webGo/todo/model"
	"strconv"
	"net/http"
	"github.com/unrolled/render"
	"github.com/gorilla/mux"
)

var rd *render.Render

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func getTodoListHandler(w http.ResponseWriter, r *http.Request) {
	// list := []*model.Todo{}
	// for _, v := range todoMap {
	// 	list = append(list, v)
	// }
	list := model.GetTodos()
	rd.JSON(w, http.StatusOK, list)
}

func addTodoListHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	todo := model.AddTodo(name)
	// id := len(todoMap) 
	// todo := &model.Todo{id, name, false, time.Now()}
	// todoMap[id] = todo
	rd.JSON(w, http.StatusCreated, todo)
}

type Success struct {
	Success bool `json:"success"`
}

func removeTodoListHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	ok := model.RemoveTodo(id)
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

func completeTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	complete := r.FormValue("complete") == "true"

	ok := model.CompleteTodo(id, complete)
	
	if ok {
		rd.JSON(w, http.StatusOK, Success{true})
	} else {
		rd.JSON(w, http.StatusBadRequest, Success{false})
	} 
}

func MakeNewHandler() http.Handler {
	// todoMap = make(map[int]*model.Todo)

	// addTestTodos()
	
	rd = render.New()
	r := mux.NewRouter()

	r.HandleFunc("/todos", getTodoListHandler).Methods("GET")
	r.HandleFunc("/todos", addTodoListHandler).Methods("POST")
	r.HandleFunc("/todos/{id:[0-9]+}", removeTodoListHandler).Methods("DELETE")
	r.HandleFunc("/complete-todo/{id:[0-9]+}", completeTodoHandler).Methods("GET")
	r.HandleFunc("/", indexHandler)

	return r
}