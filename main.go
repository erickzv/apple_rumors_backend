package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/poximy/ohmyapple/route"
)

func main() {
	s := NewServer()
	s.MountMiddleware()
	s.MountHandlers()

	err := http.ListenAndServe(":"+Port(), s.Router)
	if err != nil {
		panic(err)
	}
}

type Server struct {
	Router *chi.Mux
}

func NewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountMiddleware() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Compress(5, "text/html", "text/css"))
	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "*"},
	}))
}

func (s *Server) MountHandlers() {
	s.Router.Get("/rumors", route.Rumors)
}

func Port() (port string) {
	const defaultPort string = "8080"

	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	return
}
