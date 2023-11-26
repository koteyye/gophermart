package strutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sergeizaitcev/gophermart/pkg/strutil"
)

func TestOnlyDigits(t *testing.T) {
	testCases := []struct {
		value string
		want  bool
	}{
		{"", false},
		{"_123", false},
		{"12_3", false},
		{"123_", false},
		{"abc", false},
		{"123", true},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.want, strutil.OnlyDigits(tc.value), tc.value)
	}
}
