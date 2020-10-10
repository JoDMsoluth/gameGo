package main

import (
	"net/http"

	"github.com/JoDMsoluth/webGo/restful"
)

func main() {
	// 1. router 등록 및 테스트 환경, 없으면 nil
	// http.ListenAndServe(":3000", myapp.NewHttpHandler())

	// 2. 파일 업로드
	// http.ListenAndServe(":3000", fileupload.FileUpload())

	// 3. Restful
	http.ListenAndServe(":3000", restful.RestfulHandler())
}
