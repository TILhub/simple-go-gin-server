package main

import (
	"awesomeProject2/server/schema"
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/logharbour/logharbour"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fallbackWriter := logharbour.NewFallbackWriter(logFile, os.Stdout)
	lctx := logharbour.NewLoggerContext(logharbour.Debug2)
	logger := logharbour.NewLogger(lctx, "sampleCode", fallbackWriter)

	router := gin.Default()

	// do panic recovery as well
	router.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Err().LogDebug("panic recovered", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":       "Internal Server Error",
					"erroMessage": err,
				})
			}
		}()
		c.Next()
	})

	s := service.NewService(router).
		WithLogHarbour(logger)
	apiV1Group := router.Group("/api/v1")

	s.RegisterRouteWithGroup(apiV1Group, http.MethodPost, "/calculator", schema.CalculatorFunction)
	s.RegisterRouteWithGroup(apiV1Group, http.MethodGet, "", schema.HealthCheck)

	appServerPortStr := strconv.Itoa(8080)
	err = router.Run(":" + appServerPortStr)
	if err != nil {
		logger.LogActivity("Failed to start server", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}
