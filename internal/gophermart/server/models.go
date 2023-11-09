package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const tokenType = "Bearer"

type token struct {
	Type    string `json:"type"`
	Token   string `json:"token"`
	Expires int    `json:"expires,omitempty"`
}

func newToken(value string) token {
	return token{
		Type:  tokenType,
		Token: value,
	}
}

func parseToken(s string) (string, error) {
	split := strings.SplitN(s, " ", 2)
	if len(split) != 2 || split[0] != tokenType {
		return "", errors.New("unsupported token")
	}
	return split[1], nil
}

type user struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func decodeUser(r io.Reader) (user, error) {
	var u user

	err := json.NewDecoder(r).Decode(&u)
	if err != nil {
		return user{}, fmt.Errorf("decoding the user: %w", err)
	}
	if u.Login == "" || u.Password == "" {
		return user{}, errors.New("login or password is empty")
	}

	return u, nil
}
