package serviceServer

import (
	"fmt"
	"log"
	"net/http"
)

func (server SimpleServer) validateMethod(method string, handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			w.Header().Add("Allow", method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			errorMessage := fmt.Sprintf("%d %s", http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
			w.Write([]byte(errorMessage))

			log.Println(r.Method + " " + r.URL.Path + " " + errorMessage)

			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (server SimpleServer) AddMappingWithMethod(path string, method string, handler http.HandlerFunc) {
	server.AddMapping(path, server.validateMethod(method, handler))
}
