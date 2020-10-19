package decoratorHandler

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPage(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(DecoratorHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// body 읽기
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World", data)
}

func TestDecoHandler(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(DecoratorHandler())
	defer ts.Close()

	buf := &bytes.Buffer{}
	// 표준 로그를 버퍼에 출력한다.
	log.SetOutput(buf)

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// 버퍼에 있는 값을 한줄씩 읽어온다.
	r := bufio.NewReader(buf)
	line, _, err := r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Started")

	line, _, err := r.ReadLine()
	assert.NoError(err)
	assert.Contains(string(line), "[LOGGER1] Complted")
}
