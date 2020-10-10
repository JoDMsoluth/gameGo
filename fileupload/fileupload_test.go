package fileupload

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadTest(t *testing.T) {
	assert := assert.New(t)
	path := "../public/img.png"

	// 0. 파일 스트림 열기
	file, _ := os.Open(path)
	defer file.Close()
	fmt.Println("file", file)

	// 1. 해당 폴더 내용 전부 지운다.
	// os.RemoveAll("../uploads")

	// 2. 읽기 버퍼 생성 - 스트림에 있는 데이터를 읽어오기 위해
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	// file input에 있는 key와 파일 이름을 매개변수로 넣는다.
	// 3. 읽어온다.
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path))
	assert.NoError(err)

	// 4. 파일 스트림에서 받은 파일을 multi에 복사한다.
	io.Copy(multi, file)
	writer.Close()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf)
	req.Header.Set("Content-type", writer.FormDataContentType())

	// fileupload가 mux이므로 mux로 바꿔서 해야 제대로 적용된다.
	mux := FileUpload()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)

	uploadFilePath := "../uploads/" + filepath.Base(path)
	// 업로드된 파일 정보
	_, err = os.Stat(uploadFilePath)
	assert.NoError(err)

	// 업로드된 파일 열기
	uploadFile, _ := os.Open(uploadFilePath)
	// 기존 파일 열기
	originFile, _ := os.Open(path)
	defer uploadFile.Close()
	defer originFile.Close()

	uploadData := []byte{}
	originData := []byte{}

	// 연 파일 각각 uploadData, originData에 저장
	uploadFile.Read(uploadData)
	originFile.Read(originData)

	// 비교
	assert.Equal(originData, uploadData)

}
