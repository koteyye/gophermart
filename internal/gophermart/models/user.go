package models

// UserID тип для указания UserID в контексте
type UserID string

// CtxUserID константа для добавления/получения значения UserID в/из контекста
const CtxUserID UserID = "userID"