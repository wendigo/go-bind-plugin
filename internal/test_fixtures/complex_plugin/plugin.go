package main

import (
	"fmt"
	"net/http"

	http2 "github.com/wendigo/plugin_test/plug/http"
)

func DoWork(h *http.Header) *http.Header {
	return h
}

func PrintHello() string {
	return "Hello world!"
}

func PrintHello2(in int) string {
	return fmt.Sprintf("Hello %d", in)
}

func DoWork3() *http.Header {
	return &http.Header{}
}

func DoWork4() http.Header {
	return http.Header{}
}

func DoWorkInt(x map[string]int) []int32 {
	return []int32{}
}

func DoWorkOnChan(x <-chan int) chan<- int32 {
	return make(chan<- int32)
}

func DoWorkOnChan2(x <-chan http2.Work) chan<- http.Header {
	return make(chan<- http.Header)
}

func DoWorkIntArray(x map[string]int) [5]int32 {
	return [5]int32{0, 1, 2, 3, 4}
}

func DoWorkIntArrayVariadic(x ...map[string]int) [5]int32 {
	return [5]int32{0, 1, 2, 3, 4}
}

func DoWorkingString(x string) http2.Work {
	return http2.Work{Work: x}
}

func DoWorkArray(in []*http2.Work) http2.Work {
	return http2.Work{Work: "Hello"}
}

func DoWorkMap(m map[string]*http2.Work) *http2.Work {
	return nil
}

func DoWork2(work http.Header, work2 http.Header) string {
	return "Hello"
}

var X http.Header = http.Header{"Name": []string{"Value"}}
var Y *http.Header = &http.Header{"Name": []string{"Value2"}}

var Z func() http.Header = DoWork4
var V func(work http.Header, work2 http.Header) string = DoWork2
