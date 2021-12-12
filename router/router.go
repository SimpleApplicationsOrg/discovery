package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Handler struct {
	method    string
	handlerFn func(w http.ResponseWriter, r *http.Request)
}

type Middleware func(http.Handler) http.Handler

type Router struct {
	routes      map[string][]Handler
	registered  []string
	middlewares []Middleware
}

func New() *Router {
	router := &Router{
		routes:     make(map[string][]Handler),
		registered: make([]string, 0),
	}
	return router
}

func (r *Router) Use(m Middleware) {
	r.middlewares = append(r.middlewares, m)
}

func (r *Router) Add(method, path string, handler func(http.ResponseWriter, *http.Request)) {
	if method == "" {
		panic("router: empty method")
	}
	if path == "" {
		panic("router: empty path")
	}
	if handler == nil {
		panic("router: nil handler")
	}
	if len(path) > 1 && isParam(strings.Split(path[1:], "/")[0]) {
		panic("router: path starts with param")
	}

	handlers, found := r.routes[path]
	if !found {
		handlers = make([]Handler, 0)
		r.registered = append(r.registered, path)
	}
	h := Handler{
		method:    method,
		handlerFn: handler,
	}
	handlers = append(handlers, h)
	r.routes[path] = handlers
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	route, params := r.findRoute(req)
	methodNotAllowed := false
	if handlers, found := r.routes[route]; found {
		for _, handler := range handlers {
			if handler.method != req.Method {
				methodNotAllowed = true
				continue
			}

			v := context.WithValue(context.Background(), "pathParams", params)
			rCtx := req.WithContext(v)

			_handler := wrapHandler(handler)
			for i := len(r.middlewares) - 1; i >= 0; i-- {
				_handler = r.middlewares[i](_handler)
			}

			_handler.ServeHTTP(w, rCtx)
			return
		}
	}

	if methodNotAllowed {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (r *Router) HasRoute(req *http.Request) bool {
	if r, _ := r.findRoute(req); r != "" {
		return true
	}
	return false
}

func (r *Router) Start() {

	http.Handle("/", r)
	srv := &http.Server{Addr: ":8080"}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")

	wait(srv, done)
}

func (r *Router) findRoute(req *http.Request) (string, map[string]string) {
	p := req.URL.Path
	parts := strings.Split(p[1:], "/")
	for _, route := range r.registered {
		rParts := strings.Split(route[1:], "/")
		if len(parts) == len(rParts) {
			params := make(map[string]string)
			for i, rPart := range rParts {
				if !isParam(rPart) && parts[i] != rPart {
					break
				}
				if isParam(rPart) {
					params[getParam(rPart)] = parts[i]
				}
			}

			return route, params
		}
	}
	return "", nil
}

func isParam(s string) bool {
	return strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")
}

func getParam(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
}

func wrapHandler(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.handlerFn(w, r)
	})
}

func wait(srv *http.Server, done chan os.Signal) {

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
}
