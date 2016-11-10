package main

import "fmt"

// ReturningInt32 is only exported for testing purposes
func ReturningInt32() int32 {
	return 32
}

// ReturningStringSlice is only exported for testing purposes
func ReturningStringSlice() []string {
	return []string{"hello", "world"}
}

// ReturningIntArray is only exported for testing purposes
func ReturningIntArray() [3]int32 {
	return [...]int32{1, 0, 1}
}

// NonReturningFunction is only exported for testing purposes
func NonReturningFunction() {
	fmt.Println("I'm not returning anything")
}

// X should not be exported in the interface
var X = "Should not be exported"
