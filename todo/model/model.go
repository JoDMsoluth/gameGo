package model

import (
	"time"
)


type Todo struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Completed bool `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type dbHandler interface {
	getTodos() []*Todo
	addTodo(name string) *Todo 
	removeTodo(id int) bool
	completeTodo(id int, complete bool) bool
}

type memoryHandler struct {
	todoMap map[int]*Todo
}

func newMemoryHandler() dbHandler {
	m := &memoryHandler{}
	m.todoMap = make(map[int]*Todo)
	return m
}

var handler dbHandler

func init() {
	handler = newMemoryHandler()
}

func (m *memoryHandler) getTodos() []*Todo {
	list := []*Todo{}
	for _, v := range m.todoMap {
		list = append(list, v)
	}
	return list
}

func (m *memoryHandler) addTodo(name string) *Todo {
	id := len(m.todoMap) 
	todo := &Todo{id, name, false, time.Now()}
	m.todoMap[id] = todo
	return todo
} 

func (m *memoryHandler) removeTodo(id int) bool {
	if _, ok := m.todoMap[id]; ok {
		delete(m.todoMap, id)
		return true
	}
	return false
}
 
func (m *memoryHandler) completeTodo(id int, complete bool) bool {
	 if todo, ok := m.todoMap[id]; ok {
		todo.Completed = complete
		return true
	} 
	return false
}

func GetTodos() []*Todo {
	return handler.getTodos()
}
func AddTodos(name string) *Todo {
	return handler.addTodo(name)
}
func RemoveTodos(id int) bool {
	return handler.removeTodo(id)
}
func CompleteTodos(id int, complete bool) bool {
	return handler.completeTodo(id, complete)
}