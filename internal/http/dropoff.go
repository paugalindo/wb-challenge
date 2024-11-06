package http

import (
	"errors"
	"net/http"
	"wb-challenge/internal/commands"

	"github.com/gin-gonic/gin"
)

type DropoffRequest struct {
	ID int `json:"id" binding:"required"`
}

func (s *Server) PostDropoff(c *gin.Context) {
	req := DropoffRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.DropOffGroupCmd{
		ID: req.ID,
	}
	if err := s.cmdBus.Dispatch(c.Request.Context(), cmd); err != nil {
		if errors.Is(err, commands.ErrNotFound) {
			c.JSON(http.StatusNotFound, struct{}{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
