package app
// go-sqlite3를 사용하기 위해서는 -> g c go 필요 -> c를 표준 컴파일 위해서 gcc 필요 -> ms window에서는 c 표준 컴파일 지원 x -> tdm-gcc 설치해야함
import (
	"github.com/JoDMsoluth/webGo/todo/model"
	"fmt"
	"strconv"
	"net/url"
	"net/http"
	"encoding/json"
	"testing"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
)

func TestTodos (t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(MakeNewHandler())
	defer ts.Close()

	// FormValue로 받기 때문에 PostForm
	resp, err := http.PostForm(ts.URL+"/todos", url.Values{"name":{"Test todo"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	var todo model.Todo
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo")
	id1 := todo.ID

	// 하나 더 생성
	resp, err = http.PostForm(ts.URL+"/todos", url.Values{"name":{"Test todo2"}})
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&todo)
	assert.NoError(err)
	assert.Equal(todo.Name, "Test todo2")
	id2 := todo.ID

	
	resp, err = http.Get(ts.URL + "/complete-todo/" + strconv.Itoa(id1) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	resp, err = http.Get(ts.URL+"/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos := []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 2)
	

	for _, t := range todos {
		fmt.Println(t)
		if t.ID == id1 {
			assert.True(t.Completed)
		}
	}

	req, _ := http.NewRequest("DELETE", ts.URL + "/todos/" + strconv.Itoa(id1), nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)


	// Get
	resp, err = http.Get(ts.URL+"/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	todos = []*model.Todo{}
	err = json.NewDecoder(resp.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(len(todos), 1)
	for _, t := range todos {
		assert.Equal(t.ID, id2)
	}
}