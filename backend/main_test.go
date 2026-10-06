package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFibonizeLoop(t *testing.T) {
	var tests = []struct {
		n    int
		want int64
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{10, 55},
		{20, 6765},
		{50, 12586269025},
		{92, 7540113804746346429},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("F(%d)", tt.n), func(t *testing.T) {
			ans := FibonizeLoop(tt.n)
			if ans != tt.want {
				t.Errorf("FibonizeLoop(%d) = %d; want %d", tt.n, ans, tt.want)
			}
		})
	}
}

func TestFibonizeRecursive(t *testing.T) {
	var tests = []struct {
		n    int
		want [2]int64
	}{
		{0, [2]int64{0, 0}},
		{1, [2]int64{0, 1}},
		{2, [2]int64{0, 1}},
		{10, [2]int64{34, 55}},
		{20, [2]int64{4181, 6765}},
		{30, [2]int64{514229, 832040}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("F(%d)", tt.n), func(t *testing.T) {
			ans := FibonizeRecursiveV2(tt.n)
			if ans != tt.want {
				t.Errorf("FibonizeRecursive(%d) = %d; want %d", tt.n, ans, tt.want)
			}
		})
	}
}

func TestImplementationsAgree(t *testing.T) {
	for n := 0; n <= 25; n++ {
		loop := FibonizeLoop(n)
		recursive := FibonizeRecursiveV2(n)
		if loop != recursive[1] {
			t.Errorf("F(%d): loop = %d, recursive = %d", n, loop, recursive)
		}
	}
}

func TestFibonizeEndpoints(t *testing.T) {
	e := newServer()

	for _, method := range []string{"loop", "recursive"} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+method+"/10", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d; want %d", rec.Code, http.StatusOK)
			}

			var res FiboResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
				t.Fatalf("invalid JSON %q: %v", rec.Body.String(), err)
			}

			if res.Method != method || res.N != 10 || res.Result != "55" {
				t.Errorf("got %+v; want method=%s n=10 result=55", res, method)
			}
			if res.DurationNs < 0 {
				t.Errorf("durationNs = %d; want >= 0", res.DurationNs)
			}
		})
	}
}
