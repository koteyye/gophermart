package models

// UserID тип для указания UserID в контексте
type UserID uint8

// CtxUserID константа для добавления/получения значения UserID в/из контекста
const CtxUserID UserID = iota+1
