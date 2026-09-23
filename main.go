package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "WELCOME TO FIBONIZER.\n'/recursive' -> Recursive Fibonizer\n'/loop' -> Loop Fibonizer")
	})

	e.GET("/recursive", func(c echo.Context) error {
		start := time.Now()
		fibo := FibonizeRecursive(8)
		log.Printf("Time the Recursive Fibonizer took: %d", time.Since(start))
		return c.HTML(http.StatusOK, strconv.Itoa(fibo))
	})

	e.GET("/loop", func(c echo.Context) error {
		start := time.Now()
		fibo := FibonizeLoop(8)
		log.Printf("Time the Loop Fibonizer took: %d", time.Since(start))
		return c.HTML(http.StatusOK, strconv.Itoa(fibo))
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func FibonizeRecursive(num int) int {
	if num == 0 {
		return 0
	}

	if num == 1 {
		return 1
	}

	return FibonizeRecursive(num-2) + FibonizeRecursive(num-1)
}

func FibonizeLoop(num int) int {
	if num == 0 {
		return 0
	}

	if num == 1 {
		return 1
	}

	sum := 1
	previous := 0

	for i := 2; i <= num; i++ {
		previousSum := sum
		sum = sum + previous
		previous = previousSum
	}

	return sum
}
