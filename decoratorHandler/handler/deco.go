package decoHandler

type DecoratorFunc func (http.ResopnseWriter, *http.Request, http.Handler) {}

type DecoHandler struct {
	fn DecoratorFunc
	h http.Handler
}

// 인자로 받은 function 실행하는 함수
func (self *DecoHandler) ServeHTTP(http.ResopnseWriter, *http.Request, http.Handler) {
	self.fn(w, r, self.h)
}

// Decorator Function 만들기
func NewDecoHandler(h http.Handler, fn DecoratorFunc) http.Handler {
	return &DecoHandler{
		fn: fn,
		h: h }
}