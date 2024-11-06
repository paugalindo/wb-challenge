package http

import (
	"context"
	"net/http"
	"wb-challenge/bus"
	"wb-challenge/internal/query"

	"github.com/gin-gonic/gin"
)

type Server struct {
	srv *http.Server

	cmdBus  *bus.CommandBus
	groupQS *query.GroupQS
}

func New(port string, cmdBus *bus.CommandBus, groupQS *query.GroupQS) Server {
	s := Server{
		cmdBus:  cmdBus,
		groupQS: groupQS,
	}

	router := gin.Default()

	router.GET("/status", s.GetStatus)
	router.PUT("/evs", s.PutEVs)
	router.POST("/journey", s.PostJourney)
	router.POST("/dropoff", s.PostDropoff)
	router.POST("/locate", s.PostLocate)

	s.srv = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	return s
}

func (s *Server) Start() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
