package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sixsat/assessment/pkg/expense"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Validate authorization token
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("Authorization") != "November 10, 2009" {
				return c.JSON(http.StatusUnauthorized, "Unauthorized")
			}
			return next(c)
		}
	})

	expense.InitDB()

	e.POST("/expenses", expense.CreateExpense)
	e.GET("/expenses", expense.GetAllExpenses)
	e.GET("/expenses/:id", expense.GetExpense)
	e.PUT("/expenses/:id", expense.UpdateExpense)

	port := os.Getenv("PORT")
	log.Println("start at port:", port)

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	log.Println("byeee~")
}
