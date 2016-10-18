package serviceServer

import (
	"context"
	"net/http"
)

func (server SimpleServer) injectContext(ctx context.Context, handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (server SimpleServer) AddMappingWithContext(path string, method string, handler http.HandlerFunc, ctx context.Context) {
	server.AddMappingWithMethod(path, method, server.injectContext(ctx, handler))
}
