package main

import (
	r "awesomeProject2/router"
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/logharbour/logharbour"
	"log"
	"net/http"
	"os"
	"time"
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

	restServer := r.NewRESTServer(router, logger)
	restServer.Init()

	httpServer := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Err().LogDebug("error starting the rest server", err)
	}
}
