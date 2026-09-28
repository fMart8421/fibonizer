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

	e.GET("/recursive/:num", func(c echo.Context) error {
		num := c.Param("num")
		parsedNum, err := strconv.Atoi(num)

		if err != nil {
			log.Fatal(err)
			return err
		}

		start := time.Now()
		fibo := FibonizeRecursive(parsedNum)

		log.Printf("Time the Recursive Fibonizer took: %d", time.Since(start))

		return c.HTML(http.StatusOK, strconv.FormatInt(fibo, 10))
	})

	e.GET("/loop/:num", func(c echo.Context) error {
		num := c.Param("num")
		parsedNum, err := strconv.Atoi(num)

		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		fibo := FibonizeLoop(parsedNum)

		log.Printf("Time the Loop Fibonizer took: %d", time.Since(start))

		return c.HTML(http.StatusOK, strconv.FormatInt(fibo, 10))
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

func FibonizeRecursive(num int) int64 {
	if num == 0 {
		return 0
	}

	if num == 1 {
		return 1
	}

	return FibonizeRecursive(num-2) + FibonizeRecursive(num-1)
}

func FibonizeLoop(num int) int64 {
	if num == 0 {
		return 0
	}

	if num == 1 {
		return 1
	}

	sum := int64(1)
	previous := int64(0)

	for i := 2; i <= num; i++ {
		previousSum := sum
		sum = sum + int64(previous)
		previous = previousSum
	}

	return sum
}
