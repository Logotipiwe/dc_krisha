package domain

import "errors"

var (
	ParserNotFoundErr = errors.New("parser not found")
	LimitExceededErr  = errors.New("chat limit exceeded")
)
