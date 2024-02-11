package schema

import (
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"net/http"
)

type APPInfo struct {
	Name string
}

func HealthCheck(c *gin.Context, s *service.Service) {
	s.LogHarbour.Debug0().LogDebug("DEBUG healthCheck", "")
	c.JSON(http.StatusOK, APPInfo{
		Name: "WatchDog Bark",
	})
}
