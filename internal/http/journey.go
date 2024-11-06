package http

import (
	"errors"
	"net/http"
	"wb-challenge/internal/commands"

	"github.com/gin-gonic/gin"
)

type JourneyRequest struct {
	ID     int `json:"id" binding:"required"`
	People int `json:"people" binding:"required"`
}

func (s *Server) PostJourney(c *gin.Context) {
	req := JourneyRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.CreateGroupCmd{
		ID:     req.ID,
		People: req.People,
	}
	if err := s.cmdBus.Dispatch(c.Request.Context(), cmd); err != nil {
		if errors.Is(err, commands.ErrGroupAlreadyExist) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, struct{}{})
}
