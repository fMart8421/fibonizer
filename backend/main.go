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

type FiboError struct {
	Method  string `json:"method"`
	N       int    `json:"n"`
	Message string `json:"error"`
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

	e.GET("/recursive/v1/:num", func(c echo.Context) error {
		num := GetParsedNumber(c)

		if num < 0 {
			return c.JSON(http.StatusBadRequest, FiboError{
				Method:  "loop",
				N:       num,
				Message: "Number cannot be negative",
			})
		}

		start := time.Now()
		fibo := FibonizeRecursiveV1(num)
		duration := time.Since(start)

		log.Printf("Time the Recursive Fibonizer took: %d", duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "recursive",
			N:          num,
			Result:     strconv.FormatInt(fibo, 10),
			DurationNs: duration.Nanoseconds(),
		})
	})

	e.GET("/recursive/v2/:num", func(c echo.Context) error {
		num := GetParsedNumber(c)

		if num < 0 {
			return c.JSON(http.StatusBadRequest, FiboError{
				Method:  "loop",
				N:       num,
				Message: "Number cannot be negative",
			})
		}

		start := time.Now()
		fibo := FibonizeRecursiveV2(num)
		duration := time.Since(start)

		log.Printf("Time the Recursive Fibonizer took: %d", duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "recursive",
			N:          num,
			Result:     strconv.FormatInt(fibo[1], 10),
			DurationNs: duration.Nanoseconds(),
		})
	})

	e.GET("/loop/:num", func(c echo.Context) error {
		num := GetParsedNumber(c)

		if num < 0 {
			return c.JSON(http.StatusBadRequest, FiboError{
				Method:  "loop",
				N:       num,
				Message: "Number cannot be negative",
			})
		}

		start := time.Now()
		fibo := FibonizeLoop(num)
		duration := time.Since(start)

		log.Printf("Time the Loop Fibonizer took: %d", duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "loop",
			N:          num,
			Result:     strconv.FormatInt(fibo, 10),
			DurationNs: duration.Nanoseconds(),
		})
	})

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	return e
}

// this implementation is obviously wrong, it was temporary but for documenting reasons, we keep it as V1
func FibonizeRecursiveV1(num int) int64 {
	if num == 0 {
		return 0
	}

	if num == 1 {
		return 1
	}

	return FibonizeRecursiveV1(num-2) + FibonizeRecursiveV1(num-1)
}

// this implementation returns an array with the previous and the current number
// this version is still way slower than `FibonizeLoop`, but that's due to the call stack and all the necessary resources a new call to a function needs
func FibonizeRecursiveV2(num int) [2]int64 {
	if num == 0 {
		return [2]int64{0, 0}
	}

	if num == 1 {
		return [2]int64{0, 1}
	}

	fibonized := FibonizeRecursiveV2(num - 1)

	return [2]int64{fibonized[1], fibonized[0] + fibonized[1]}
}

// a simple loop that goes up to the num (meaning it runs N times), but the most efficient I came up with
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

func GetParsedNumber(c echo.Context) int {
	num := c.Param("num")
	parsedNum, err := strconv.Atoi(num)

	if err != nil {
		log.Fatal(err)
	}

	return parsedNum
}
