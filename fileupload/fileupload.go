package fileupload

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadsHandler(w http.ResponseWriter, r *http.Request) {

	// 1. 보내준 파일 Read
	// input에서 넣은 key랑 동일해야 됨
	uploadFile, header, err := r.FormFile("upload_file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	defer uploadFile.Close()

	// 2. 저장할 공간의 파일 생성
	dirname := "uploads"
	// 폴더 없으면 폴더 만들기
	os.MkdirAll("../"+dirname, 0777)
	filepath := fmt.Sprintf("../%s/%s", dirname, header.Filename)

	// 파일 생성 핸들러 사용 - 사용 끝나면 자원 반납
	file, err := os.Create(filepath)
	defer file.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	// 3. 만들 파일 Copy
	io.Copy(file, uploadFile)

	// 4. response 전송
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, filepath)

}

func FileUpload() http.Handler {
	// 라우터 인스턴스 만들기
	mux := http.NewServeMux()

	mux.HandleFunc("/uploads", uploadsHandler)
	mux.Handle("/", http.FileServer(http.Dir("public")))

	return mux
}
