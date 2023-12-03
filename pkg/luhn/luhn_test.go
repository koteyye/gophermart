package luhn_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/pkg/luhn"
)

func TestCheck(t *testing.T) {
	testCases := []struct {
		value string
		want  bool
	}{
		{"49927398716", true},
		{"49927398717", false},
		{"1234567812345678", false},
		{"1234567812345670", true},
	}

	for _, tc := range testCases {
		got := luhn.Check(tc.value)
		require.Equal(t, tc.want, got, tc.value)
	}
}
