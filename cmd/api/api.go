package api

import (
	"database/sql"
	"log"
	"template/services/user"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := gin.Default()
	apiGroup := router.Group("/api/v1")

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(apiGroup)

	log.Println("Listening on", s.addr)

	return router.Run(s.addr)
}
