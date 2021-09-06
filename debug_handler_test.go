package zql

import (
	"net/http"
)

func ExampleDebugHandler_ServeHTTP() {
	h := new(DebugHandler)
	http.ListenAndServe(":8080", h)
}
