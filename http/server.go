package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	log    *log.Logger
	router *gin.Engine
	server *http.Server
	close  chan struct{}
}

func Create(log *log.Logger, address string) *HttpServer {
	gin.SetMode(gin.ReleaseMode)

	r := &HttpServer{
		log:    log,
		router: gin.New(),
		server: &http.Server{
			Addr: address,
		},
		close: make(chan struct{}),
	}

	r.server.Handler = r.router

	r.router.NoMethod(methodNotAllowedHandler)
	r.router.NoRoute(notFoundHandler)

	r.router.Use(prepareKeys())
	r.router.Use(gin.Recovery())
	r.router.Use(r.loggingMiddleware())

	return r
}

func (r *HttpServer) Router() *gin.Engine {
	return r.router
}

func (r *HttpServer) Run() {
	r.log.Printf("Starting REST server on %s\n", r.server.Addr)

	go func() {
		<-r.close
		if err := r.server.Close(); err != nil {
			log.Fatalf("Server Close Error: %v\n", err)
		}
	}()

	if err := r.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("HTTP Server closed")
		} else {
			log.Fatalf("HTTP Listen Error: %v\n", err)
		}
	}
}

func (r *HttpServer) Close() {
	close(r.close)
}
