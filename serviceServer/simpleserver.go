package serviceServer

import (
	"fmt"
	"github.com/braintree/manners"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const defaultPort string = "8080"

func NewSimpleServer() SimpleServer {
	log.Println("Creating a server instance...")
	server := SimpleServer{}
	server.mux = http.NewServeMux()
	server.instance = manners.NewServer()
	server.configure()
	return server
}

type SimpleServer struct {
	instance *manners.GracefulServer
	mux      *http.ServeMux
}

func (server SimpleServer) configure() {
	log.Println("Configuring server instance...")

	serverPort := os.Getenv("SERVER_PORT")
	if len(serverPort) == 0 {
		serverPort = defaultPort
	}
	server.instance.Addr = fmt.Sprintf(":%s", serverPort)

	server.instance.Handler = server.loggingHandler()
}

func (server SimpleServer) loggingHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Printf("%s - - [%s] \"%s %s %s\" %s\n", r.RemoteAddr, time.Now().Format(time.RFC1123),
			r.Method, r.URL.Path, r.Proto, r.UserAgent())

		server.mux.ServeHTTP(w, r)

	})
}

func (server SimpleServer) AddMapping(path string, handler http.HandlerFunc) {
	server.mux.HandleFunc(fmt.Sprintf("/%s/", path), handler)
}

func (server SimpleServer) Start() {
	log.Println("Starting server intance on " + server.instance.Addr)

	errChan := make(chan error, 10)

	go func() {
		errChan <- server.instance.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			server.instance.BlockingClose()
			os.Exit(0)
		}
	}
}
