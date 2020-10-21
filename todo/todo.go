package todo

import (
	"log"
	"net/http"
	"github.com/JoDMsoluth/webGo/todo/app"
)

func Todo () {
	// 경로같은 민감한 부분은 환경변수 등으로 실행인자로 받아올 수 있다.
	m := app.MakeNewHandler("./test.db")
	defer m.Close()

	log.Println("Started App")

	err := http.ListenAndServe(":3000", m)

	if err != nil {
		panic(err)
	}
}