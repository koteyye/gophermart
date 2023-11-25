package user

import "context"

// keyUserID определяет ключ для передачи UserID через контекст.
var keyUserID struct{}

// WithContext устанавливает UserID в контекст.
func WithContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, keyUserID, userID)
}

// FromContext возвращает UserID из контеста.
func FromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(keyUserID).(string)
	return userID, ok
}
