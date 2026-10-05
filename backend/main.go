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

type FiboResponse struct {
	Method string `json:"method"`
	N      int    `json:"n"`
	// Sent as a string because int64 values above 2^53 lose precision as JS numbers
	Result     string `json:"result"`
	DurationNs int64  `json:"durationNs"`
}

func main() {
	e := newServer()

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func newServer() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

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
		duration := time.Since(start)

		log.Printf("Time the Recursive Fibonizer took: %d", duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "recursive",
			N:          parsedNum,
			Result:     strconv.FormatInt(fibo[1], 10),
			DurationNs: duration.Nanoseconds(),
		})
	})

	e.GET("/loop/:num", func(c echo.Context) error {
		num := c.Param("num")
		parsedNum, err := strconv.Atoi(num)

		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		fibo := FibonizeLoop(parsedNum)
		duration := time.Since(start)

		log.Printf("Time the Loop Fibonizer took: %d", duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "loop",
			N:          parsedNum,
			Result:     strconv.FormatInt(fibo, 10),
			DurationNs: duration.Nanoseconds(),
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	return e
}

func FibonizeRecursive(num int) [2]int64 {
	if num == 0 {
		return [2]int64{0, 0}
	}

	if num == 1 {
		return [2]int64{0, 1}
	}

	fibonized := FibonizeRecursive(num - 1)

	return [2]int64{fibonized[1], fibonized[0] + fibonized[1]}
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
