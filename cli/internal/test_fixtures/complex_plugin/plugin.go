package main

import "C"

import (
	"fmt"
	"net/http"

	http2 "github.com/wendigo/go-bind-plugin/cli/internal/test_fixtures/complex_plugin/http"
)

// PrintHello is only exported for testing purposes
func PrintHello() string {
	return "Hello world!"
}

// PrintHello2 is only exported for testing purposes
func PrintHello2(in int) string {
	return fmt.Sprintf("Hello %d", in)
}

// DoWorkInt is only exported for testing purposes
func DoWorkInt(x map[string]int) []int32 {
	return []int32{}
}

// DoWorkOnChan is only exported for testing purposes
func DoWorkOnChan(x <-chan int) chan<- int32 {
	return make(chan<- int32)
}

// DoWorkOnChan2 is only exported for testing purposes
func DoWorkOnChan2(x <-chan http2.Work) chan<- http2.Work {
	return make(chan<- http2.Work)
}

// DoWorkIntArray is only exported for testing purposes
func DoWorkIntArray(x map[string]int) [5]int32 {
	return [5]int32{0, 1, 2, 3, 4}
}

// DoWorkIntArrayVariadic is only exported for testing purposes
func DoWorkIntArrayVariadic(x ...map[string]int) [5]int32 {
	return [5]int32{0, 1, 2, 3, 4}
}

// DoWorkingString is only exported for testing purposes
func DoWorkingString(x string) http2.Work {
	return http2.Work{Work: x}
}

// DoWorkArray is only exported for testing purposes
func DoWorkArray(in []*http2.Work) http2.Work {
	return http2.Work{Work: "Hello"}
}

// DoWorkMap is only exported for testing purposes
func DoWorkMap(m map[string]*http2.Work) *http2.Work {
	return nil
}

// DoWorkVariadic is only exported for testing purposes
func DoWorkVariadic(m ...string) bool {
	return false
}

// DoWorkOnTwoNamedTypes is only exported for testing purposes
func DoWorkOnTwoNamedTypes(l http.Header, r http2.Work) http.Dir {
	return ""
}

// X is only exported for testing purposes
var X = http2.Work{Work: "Hello world!"}

// Y is only exported for testing purposes
var Y = &http2.Work{Work: "Hello"}
