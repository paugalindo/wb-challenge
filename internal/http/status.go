package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GetStatus(c *gin.Context) {
	c.JSON(http.StatusOK, struct{}{})
}
