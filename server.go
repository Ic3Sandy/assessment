package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Ic3Sandy/assessment/expenses"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

var db *sql.DB

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		key := c.Request().Header.Get("Authorization")
		if key != os.Getenv("AUTHORIZATION_KEY") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid key")
		}
		return next(c)
	}
}

func main() {
	fmt.Println("Please use server.go for main file")
	fmt.Println("start at port:", os.Getenv("PORT"))

	tableName := "expenses"
	db = ConnectAndCreateTable(tableName)
	expenses.SetDB(db)
	expenses.SetTableName(tableName)
	defer db.Close()

	e := echo.New()
	e.GET("/expenses", expenses.GetExpenses, AuthMiddleware)
	e.GET("/expenses/:id", expenses.GetExpensesById)
	e.POST("/expenses", expenses.CreateExpense)
	e.PUT("/expenses/:id", expenses.UpdateExpense)

	go func() {
		if err := e.Start(":" + os.Getenv("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
