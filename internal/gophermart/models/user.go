package models

// UserID тип для указания UserID в контексте
type CtxKey uint8

// CtxUserID константа для добавления/получения значения UserID в/из контекста
const CtxUserID CtxKey = iota+1
