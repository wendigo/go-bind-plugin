package http

type Work struct {
	Work string
}

func DoWork2(w Work) string {
	return "I'm doing my work: " + w.Work
}
