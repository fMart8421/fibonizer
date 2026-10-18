package main

import (
	"log"
	"math"
	"math/big"
	"net/http"
	"os"
	"sort"
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

		log.Printf("\n@@@@\nRecursive Fibonizer v1\nNumber:%d\nTime: %d\n@@@@\n", fibo, duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "recursive/v1",
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

		log.Printf("\n@@@@\nRecursive Fibonizer v2\nNumber:%d\nTime: %d\n@@@@\n", fibo, duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "recursive/v2",
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

		log.Printf("\n@@@@\nLoop Fibonizer\nNumber:%d\nTime: %d\n@@@@\n", fibo, duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "loop",
			N:          num,
			Result:     strconv.FormatInt(fibo, 10),
			DurationNs: duration.Nanoseconds(),
		})
	})

	e.GET("/binet/:num", func(c echo.Context) error {
		num := GetParsedNumber(c)

		if num < 0 {
			return c.JSON(http.StatusBadRequest, FiboError{
				Method:  "loop",
				N:       num,
				Message: "Number cannot be negative",
			})
		}

		start := time.Now()
		fibo := FibonizeBinet(num)
		duration := time.Since(start)

		log.Printf("\n@@@@\nBinet Fibonizer\nNumber:%d\nTime: %d\n@@@@\n", fibo, duration)

		return c.JSON(http.StatusOK, FiboResponse{
			Method:     "binet",
			N:          num,
			Result:     fibo.String(),
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

// Binet's formula is the best way to calculate very big numbers since it's O(n)=log2(n)
//
// the formula is:
// (a,b) = (F(k), F(k+1))
// c = a*(2b-a) -> F(2k)
// d = a^2+b^2 -> F(2k+1)
//
// if n is even -> new pair = (c, d)
// if n is odd -> new pair = (d, c+d)

func FibonizeBinet(num int) *big.Int {
	var halves = []int{}

	numFloat64 := float64(num)

	for i := numFloat64; i > 1; i = math.Floor(i / 2) {
		halves = append(halves, int(i))
	}

	sort.Slice(halves, func(i, j int) bool {
		return j-i > 0
	})

	a, b := big.NewInt(0), big.NewInt(1)
	c, d := big.NewInt(0), big.NewInt(1)
	twoBigInt := big.NewInt(2)

	for n := range halves {
		// the steps in the next expression are
		// c = 2b
		// c = c - a = 2b - a
		// c = c * a = a(2b - a)
		c.Mul(twoBigInt, b).Sub(c, a).Mul(c, a)

		// this one is easier (because we need to save an intermediate value)
		d.Exp(a, twoBigInt, nil) // d=a^2

		tempBSquared := big.NewInt(0).Exp(b, twoBigInt, nil) // save the value of b^2

		d.Add(d, tempBSquared) // d = d + b^2 = a^2 + b^2

		if n%2 == 0 {
			a, b = c, d
		} else {
			a = d
			b.Add(c, d)
		}
	}

	if num%2 == 0 {
		return b
	}

	return a
}

func GetParsedNumber(c echo.Context) int {
	num := c.Param("num")
	parsedNum, err := strconv.Atoi(num)

	if err != nil {
		log.Fatal(err)
	}

	return parsedNum
}

// gets the cumulative sum for the initialN, given the previousValues, until maxN
func getCumulativeSum(initialN int, previousValues [2]*big.Int, maxN int) *big.Int {
	sum, previous := previousValues[1], previousValues[0] // F(N-1), F(N-2)
	previousSum := big.NewInt(0)

	for i := initialN; i <= maxN; i++ {
		previousSum = sum      // save F(i-1)
		sum.Add(sum, previous) // F(i) = F(i-1) + F(i-2)
		previous = previousSum // next iteration, F(N-2) will be this iteration's F(i-1)
	}

	return sum
}
