package http

import (
	"errors"
	"net/http"
	"wb-challenge/internal/query"

	"github.com/gin-gonic/gin"
)

type LocateRequest struct {
	ID int `json:"id" binding:"required"`
}

type LocateResponse struct {
	VehicleID int `json:"vehicle_id"`
}

func (s *Server) PostLocate(c *gin.Context) {
	req := LocateRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vehicleID, err := s.groupQS.FindAssignedVehicle(req.ID)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			c.JSON(http.StatusNotFound, struct{}{})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if vehicleID == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, LocateResponse{VehicleID: vehicleID})
}
