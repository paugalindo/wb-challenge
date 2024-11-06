package http

import (
	"net/http"
	"wb-challenge/internal/commands"

	"github.com/gin-gonic/gin"
)

type EVsRequest []EV

type EV struct {
	ID    int `json:"id" binding:"required"`
	Seats int `json:"seats" binding:"required"`
}

func (s *Server) PutEVs(c *gin.Context) {
	req := EVsRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.LoadVehiclesCmd{
		Vehicles: make([]commands.Vehicle, 0, len(req)),
	}
	for _, v := range req {
		cmd.Vehicles = append(cmd.Vehicles, commands.Vehicle{ID: v.ID, Seats: v.Seats})
	}

	if err := s.cmdBus.Dispatch(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}
