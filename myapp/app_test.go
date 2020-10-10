package myapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 파일이 갱신될 때마다 테스트를 실행해주는 GoConvey 패키지 이용

// 컨벤션 : 함수명 앞에가 Test로 시작하고 파라미터를 *testing.T 를 받는다.
func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)
	// 이거 이렇게 하나 하나 하기 귀찮으니까 testfy 패키지 이용
	/*
		if res.Code != http.StatusOK {
			t.Fatal("Failed!! ", res.Code)
		}
	*/
	assert.Equal(http.StatusOK, res.Code)

	// body를 모두 읽어온다.
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World", string(data))
}

func TestBarPathHandler_WithoutName(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	// mux로 바꿔서 해야 제대로 적용된다.
	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)

	// body를 모두 읽어온다.
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello World", string(data))
}

func TestBarPathHandler(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=Hyehyong", nil)

	// mux로 바꿔서 해야 제대로 적용된다.
	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)

	// body를 모두 읽어온다.
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal("Hello Hyehyong", string(data))
}

func TestFooHandler_WithoutJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()
	// body를 넣어주지 않으면 BadRequest가 난다.
	req := httptest.NewRequest("GET", "/foo", nil)

	// mux로 바꿔서 해야 제대로 적용된다.
	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
}

func TestFooHandler_WithJson(t *testing.T) {
	assert := assert.New(t)

	res := httptest.NewRecorder()

	// NewReader Method로 Json이 ioreader로 바뀐다.
	req := httptest.NewRequest("POST", "/foo", strings.NewReader(`{
		"first_name" : "Jo",
		"last_name" : "Hyehyeong",
		"email" : "jodmsoluth@gmail.com"
	}`))

	// mux로 바꿔서 해야 제대로 적용된다.
	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusCreated, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)
	assert.Nil(err)
	assert.Equal("Hyehyeong", user.LastName)
	assert.Equal("Jo", user.FirstName)
}
