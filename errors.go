package main

type InvalidSocketTokenError struct {
	message string
}

func (e InvalidSocketTokenError) Error() string {
	return e.message
}
