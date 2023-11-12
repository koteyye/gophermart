package models

import "errors"

//Ошибки storage
var (
	ErrDuplicate          = errors.New("duplicate value")
	ErrNotFound           = errors.New("value not found")
	ErrOther              = errors.New("other storage error")
)