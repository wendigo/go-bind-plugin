package main

import (
	"fmt"
	"net/http"

	http2 "github.com/wendigo/go-bind-plugin/internal/test_fixtures/complex_plugin/http"
)

// DoWork is only exported for testing purposes
func DoWork(h *http.Header) *http.Header {
	return h
}

// PrintHello is only exported for testing purposes
func PrintHello() string {
	return "Hello world!"
}

// PrintHello2 is only exported for testing purposes
func PrintHello2(in int) string {
	return fmt.Sprintf("Hello %d", in)
}

// DoWork3 is only exported for testing purposes
func DoWork3() *http.Header {
	return &http.Header{}
}

// DoWork4 is only exported for testing purposes
func DoWork4() http.Header {
	return http.Header{}
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
func DoWorkOnChan2(x <-chan http2.Work) chan<- http.Header {
	return make(chan<- http.Header)
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

// DoWork2 is only exported for testing purposes
func DoWork2(work http.Header, work2 http.Header) string {
	return "Hello"
}

// X is only exported for testing purposes
var X = http.Header{"Name": []string{"Value"}}

// Y is only exported for testing purposes
var Y = &http.Header{"Name": []string{"Value2"}}

// Z is only exported for testing purposes
var Z = DoWork4

// V is only exported for testing purposes
var V = DoWork2
