package main

import (
	"github.com/JoDMsoluth/webGo/todo"
	// "net/http"
)

func main() {
	// 1. router 등록 및 테스트 환경, 없으면 nil
	// http.ListenAndServe(":3000", myapp.NewHttpHandler())

	// 2. 파일 업로드
	// http.ListenAndServe(":3000", fileupload.FileUpload())

	// 3. Restful
	// http.ListenAndServe(":3000", restful.RestfulHandler())

	// 4. Decorator
	// decorator.Decorator()

	// 5. DecoratorHandler
	// http.ListenAndServe(":3000", decoratorHandler.DecoratorHandler())
	
	// 6. Template
	// template.Template()
	
	// 7. 유용한 패키지 render(각종 res를 을 편리하게 적용), pat(라우터 추가), negroni(http 각종 미들웨어서버 제공)
	// utilPackage.UtilPackage()

	// 8. EventSource를 이용한 채팅
	// eventSource.EventSource()

	// 9. OAuth Login
	// oauth.OAuth()

	// 10. Todo list
	todo.Todo()
}
