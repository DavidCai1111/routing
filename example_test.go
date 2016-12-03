package routing_test

import (
	"fmt"
	"net/http"

	"github.com/DavidCai1993/routing"
)

var httpHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
})

func Example() {
	router := routing.New()

	router.Define("/:type(a|b)/:id(0-9a-f]{24})", httpHandler)

	callback, params, ok := router.Match("/a/8")

	fmt.Println(ok)
	// -> true

	fmt.Println(callback.(http.Handler))

	fmt.Println(params)
	// -> map[string]string{"type": "a", "id": "8"}
}
