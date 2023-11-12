package models

import "errors"

// ошибки storage
var (
	ErrDuplicate          = errors.New("duplicate value")
	ErrDuplicateOtherUser = errors.New("duplicate value from other user")
	ErrNotFound           = errors.New("value not found")
	ErrOther              = errors.New("other storage error")
	ErrBalanceBelowZero   = errors.New("balance can't be below zero")
	ErrInvalidPassword = errors.New("invalid password")
)