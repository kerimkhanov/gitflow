package server

import (
	"gitflow/config"
	"gitflow/internal/email"
	"gitflow/internal/email/delivery"
	"gitflow/internal/email/repository"
	"gitflow/internal/email/usecase"
	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type Server struct {
	httpServer  http.Server
	Logger      *logrus.Logger
	AuthUseCase email.UseCase
}

func NewApp(config config.Config) *Server {
	_, err := initRedis()
	if err != nil {
		log.Fatalf("Redis Db initialization error %v", err)
	}

	repository := repository.NewEmailRepository()
	usecase := usecase.NewEmailUseCase(repository)
	handler := delivery.NewHandler(usecase)
	handler.InitRoutes()
	return nil
}

func initRedis() (*gocelery.CeleryClient, error) {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	cli, _ := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		1,
	)
	return cli, nil
}

func (s *Server) Run(addr string) error {
	router := http.NewServeMux()
	s.httpServer = http.Server{
		Addr:    addr,
		Handler: router,
	}
	log.Printf("Server started: http://localhost%s", addr)
	return s.httpServer.ListenAndServe()
}
