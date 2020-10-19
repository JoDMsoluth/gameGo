package decoratorHandler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	decoHandler "github.com/JoDMsoluth/webGo/decoratorHandler/handler"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

// decorator 패턴
func logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Print("[LOGGER1] Started")
	h.ServeHTTP(w, r)
	log.Print("[LOGGER1] Completed time:", time.Since(start).Milliseconds())
	log.Println()
}

// decorator 패턴
func logger2(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Print("[LOGGER2] Started")
	h.ServeHTTP(w, r)
	log.Print("[LOGGER2] Completed time:", time.Since(start).Milliseconds())
	log.Println()
}

func DecoratorHandler() http.Handler {
	// 라우터 인스턴스 만들기
	mux := http.NewServeMux()

	// mux.HandleFunc("/", indexHandler)

	h := decoHandler.NewDecoHandler(mux, logger)
	h := decoHandler.NewDecoHandler(mux, logger2)

	return h
}
