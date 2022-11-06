package server

import (
	"gitflow/config"
	"gitflow/internal/email/delivery"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type Server struct {
	router *gin.Engine
	logger *logrus.Logger
	config *config.Config
	//AuthUseCase email.UseCase
}

func NewApp(cfg config.Config, log logrus.Logger) *Server {
	return &Server{
		config: &cfg,
		logger: &log,
		router: gin.Default(),
	}
}

func (s *Server) Run() error {
	s.initEndpoints()
	srv := http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.router,
	}
	log.Printf("Server started: http://localhost:%s", s.config.Port)
	return srv.ListenAndServe()
}

func (s *Server) initEndpoints() {
	delivery.InitEmailRoutes(s.router)
}
