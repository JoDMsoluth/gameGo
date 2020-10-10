package restful

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	// 핸들러를 등록해서, 목업 서버를 만든다
	ts := httptest.NewServer(RestfulHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	// response는 성공
	assert.Equal(http.StatusOK, resp.StatusCode)

	// body는 Hello World
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", string(data))
}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	// 핸들러를 등록해서, 목업 서버를 만든다
	ts := httptest.NewServer(RestfulHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	// response는 성공
	assert.Equal(http.StatusOK, resp.StatusCode)

	// body는 Hello World
	data, _ := ioutil.ReadAll(resp.Body)
	// 문자열이 포함되어 있어야 한다.
	assert.Contains(string(data), "No Users")
}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	// 핸들러를 등록해서, 목업 서버를 만든다
	ts := httptest.NewServer(RestfulHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	// response는 성공
	assert.Equal(http.StatusOK, resp.StatusCode)

	// body는 Hello World
	data, _ := ioutil.ReadAll(resp.Body)
	// 문자열이 포함되어 있어야 한다.
	assert.Contains(string(data), "No User Id:89")
}

func TestCreateInfo(t *testing.T) {
	assert := assert.New(t)

	// 핸들러를 등록해서, 목업 서버를 만든다
	ts := httptest.NewServer(RestfulHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{
		"first_name" : "Jo",
		"last_name" : "Hyehyeong",
		"email" : "jodmsoluth@gmail.com"
	}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// Post 결과 파스 처리
	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.Equal(1, user.ID)

	// 방금 만든 유저 불러오기
	id := user.ID
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id))
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// 결과를 서로 비교
	user2 := new(User)
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID)
	assert.Equal(user.FirstName, user2.FirstName)

	// 생성한 유저 지우기
	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	// 주의 ) DELETE는 아래와 같이 사용해야 함
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	// 지움
	assert.Contains(string(data), "Deleted User Id:1")
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	// 핸들러를 등록해서, 목업 서버를 만든다
	ts := httptest.NewServer(RestfulHandler())
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	// 지울 아이디가 없었다.
	assert.Contains(string(data), "No User Id:1")
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	// 핸들러를 등록해서, 목업 서버를 만든다
	ts := httptest.NewServer(RestfulHandler())
	defer ts.Close()

	// PUT, DELETE는 다음과 같이 해야함
	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(`{"id":1, "first_name":"updated", "last_name":"updated", "email":"updated@naver.com"}`))
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id:1")

	// Post로 유저 생성
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"Jo", "last_name":"Hyehyeong", "email":"jodmsoluth@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// 생성된 유저정보
	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	// Sprintf는 결과 문자열을 반환
	updateStr := fmt.Sprintf(`{"id":%d, "first_name":"jason"}`, user.ID)

	// 방금 만든 유저의 특정 정보만 수정
	req, _ = http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(updateStr))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// 수정된 유저정보
	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)
	assert.NoError(err)
	// 생성한 유저와, 수정된 유저의 ID 비교
	log.Print(user, updateUser, "_______________________")
	assert.Equal(updateUser.ID, user.ID)
	assert.Equal("jason", updateUser.FirstName)
	assert.Equal(user.LastName, updateUser.LastName)
	assert.Equal(user.Email, updateUser.Email)
}

func TestUsers_WithUsersData(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(RestfulHandler())
	defer ts.Close()

	// 1번 유저 등록
	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"tucker", "last_name":"kim", "email":"tucker@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// 2번 유저 등록
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"jason", "last_name":"park", "email":"jason@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	// 유저리스트 불러오기
	resp, err = http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	users := []*User{}
	err = json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))
}
