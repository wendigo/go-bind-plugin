package http

// Work is only exported for testing purposes
type Work struct {
	Work string
}

// DoWork2 is only exported for testing purposes
func DoWork2(w Work) string {
	return "I'm doing my work: " + w.Work
}
