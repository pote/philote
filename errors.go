package main

type InvalidTokenError struct {
	message string
}

func (e InvalidTokenError) Error() string {
	return e.message
}
