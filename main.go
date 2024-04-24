package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Wariddon/assessment-tax/docs"
	"github.com/Wariddon/assessment-tax/tax"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

// @title Swagger KBTG Assessment Tax - posttest
// @version 1.0
// @description This is a swagger with Echo.
// @BasePath /
func main() {

	e := echo.New()

	// Middleware for logging and recovering from panics
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/healthcheck", healthCheck)
	// e.GET("/swagger/*", echoSwagger.WrapHandler)
	// e.GET("/test", tax.GetTest)

	// eAdmin := e.Group("/admin")
	// eAdmin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))
	// eAdmin.GET("/test", tax.GetTest)

	eTax := e.Group("/tax")
	eTax.POST("/calculations", tax.CalculationTax)

	// Start server
	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed { // Start server
			e.Logger.Fatal("error starting the server")
		}
	}()

	// Shutdown
	// graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	killSignal := <-signals
	switch killSignal {
	case os.Interrupt:
		fmt.Println("Got SIGINT...")
	case syscall.SIGTERM:
		fmt.Println("got SIGTERM...")
	}
	fmt.Println("shutting down the server")
	err := e.Shutdown(context.Background())
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
}

// HealthCheckHandler godoc
// @summary Health Check
// @description Health checking for the service
// @id HealthCheckHandler
// @produce plain
// @router /healthcheck [get]
func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Ok",
	})
}
