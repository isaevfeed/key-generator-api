package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	gen "isaevfeed/internal/app"
	cache "isaevfeed/internal/cache"

	"github.com/gorilla/mux"
)

type Server struct {
	Addr   string
	Router *mux.Router
}

func New(Addr string, Router *mux.Router) *Server {
	return &Server{Addr, Router}
}

func (s *Server) Listen() {
	srv := &http.Server{
		Handler: s.Router,
		Addr:    s.Addr,
	}

	s.Router.HandleFunc("/api/v1/generate", s.HandleGenerate()).Methods("POST")
	s.Router.HandleFunc("/hello", s.HandleHello()).Methods("GET")

	go func() {
		srv.ListenAndServe()
		log.Printf("Server is stopped")
	}()

	log.Printf("Server is running on the %s", s.Addr)

	// Wait for an interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func (s *Server) HandleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Service's working normally")
	}
}

func (s *Server) HandleGenerate() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		keyType := req.FormValue("type")
		maxKeyLen := req.FormValue("maxLength")
		maxKeyLenNum, err := strconv.ParseFloat(maxKeyLen, 64)
		if err != nil {
			io.WriteString(w, "Parse maxLength error")
		}

		gen := gen.NewGenerator(keyType, maxKeyLenNum, initRedis())

		io.WriteString(w, MakeResponseStringify(200, gen.GenerateKey()))

	}
}

func initRedis() *cache.Cache {
	redisHost, _ := os.LookupEnv("REDIS_HOST")
	redisPort, _ := os.LookupEnv("REDIS_PORT")
	redisPassword, _ := os.LookupEnv("REDIS_PASSWORD")
	redisDatabaseEnv, _ := os.LookupEnv("REDIS_DATABASE")
	redisDatabase, _ := strconv.Atoi(redisDatabaseEnv)

	return cache.NewCache(fmt.Sprintf("%s:%s", redisHost, redisPort), redisPassword, redisDatabase)
}
