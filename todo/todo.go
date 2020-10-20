package todo

import (
	"log"
	"net/http"
	"github.com/JoDMsoluth/webGo/todo/app"
	"github.com/urfave/negroni"
)

func Todo () {
	m := app.MakeNewHandler()
	n := negroni.Classic()
	n.UseHandler(m)


	log.Println("Started App")

	err := http.ListenAndServe(":3000", n)

	if err != nil {
		panic(err)
	}
}