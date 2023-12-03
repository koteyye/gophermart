package sign

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ErrTokenExpired возвращается, если токен просрочен.
var ErrTokenExpired = errors.New("token is expired")

// Option устанавливает не обязательные свойства для Signer.
type Option func(*jwtSigner)

// WithTTL устанавливает время жизни токена.
func WithTTL(d time.Duration) Option {
	return func(s *jwtSigner) {
		s.ttl = d
	}
}

type claims struct {
	jwt.RegisteredClaims
	Payload string `json:"payload"`
}

// Signer представляет интерфейс подписанта полезной нагрузки.
type Signer interface {
	// Sign подписывает полезную нагрузку и возвращает токен.
	Sign(payload string) (token string, err error)

	// Parse парсит токен и возвращает полезную нагрузку.
	Parse(token string) (payload string, err error)
}

// jwtSigner определеяет подписанта полезной нагрузки, основанного на JWT.
type jwtSigner struct {
	secret []byte
	ttl    time.Duration
}

// New возвращает новый экземпляр Signer.
func New(secret []byte, opts ...Option) Signer {
	s := &jwtSigner{
		secret: secret,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Sign реализует интерфейс Signer.
func (s *jwtSigner) Sign(payload string) (token string, err error) {
	claims := claims{Payload: payload}
	if s.ttl > 0 {
		expirationTime := time.Now().Add(s.ttl)
		claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = t.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("signing the payload: %w", err)
	}

	return token, nil
}

// Parse реализует интерфейс Signer.
func (s *jwtSigner) Parse(token string) (payload string, err error) {
	claims, err := s.parseToken(token)
	if err != nil {
		return "", fmt.Errorf("token parsing: %w", err)
	}
	return claims.Payload, nil
}

func (s *jwtSigner) parseToken(token string) (claims, error) {
	var c claims

	_, err := jwt.ParseWithClaims(token, &c, s.secretKey)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			err = ErrTokenExpired
		}
		return claims{}, err
	}

	return c, nil
}

func (s *jwtSigner) secretKey(t *jwt.Token) (any, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
	}
	return s.secret, nil
}
