package main

func ReturningInt32() int32 {
	return 32
}

func ReturningStringSlice() []string {
	return make([]string, 0)
}

func ReturningIntArray() [3]int32 {
	return [...]int32{1, 0, 1}
}
