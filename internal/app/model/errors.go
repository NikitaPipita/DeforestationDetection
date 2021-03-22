package model

import "errors"

var (
	ErrIncorrectRole  = errors.New("incorrect user role")
	ErrIncorrectState = errors.New("incorrect iot state")
	ErrIncorrectType  = errors.New("incorrect iot type")
)
