package models

// CtxKey определяет ключ для передачи значения через контекст.
type CtxKey uint8

const (
	// KeyUserID константа для добавления/получения значения UserID в/из контекста
	KeyUserID CtxKey = iota + 1
)
