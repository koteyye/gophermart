package monetary_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/pkg/monetary"
)

func TestFormat(t *testing.T) {
	testCases := []struct {
		value float64
		want  monetary.Unit
	}{
		{0, 0},
		{-0, 0},
		{+0, 0},
		{math.NaN(), 0},
		{math.Inf(-1), 0},
		{math.Inf(+1), 0},
		{1, 1e2},
		{1e5, 1e7},
	}

	for _, tc := range testCases {
		got := monetary.Format(tc.value)
		require.Equal(t, tc.want, got)
	}
}

func TestUnit(t *testing.T) {
	a := monetary.Format(9.154)  // 9.15
	b := monetary.Format(10.847) // 10.85

	c := a + b

	require.Equal(t, float64(20), c.Float64())
}

func TestUnit_JSON(t *testing.T) {
	toJSON := func(u monetary.Unit) []byte {
		data, _ := json.Marshal(u)
		return data
	}

	testCases := []struct {
		json      []byte
		wantUnit  monetary.Unit
		wantError bool
	}{
		{
			json:      []byte("true"),
			wantUnit:  monetary.Unit(0),
			wantError: true,
		},
		{
			json:      []byte("null"),
			wantUnit:  monetary.Unit(0),
			wantError: false,
		},
		{
			json:      toJSON(1000),
			wantUnit:  1000,
			wantError: false,
		},
		{
			json:      toJSON(1045),
			wantUnit:  monetary.Unit(1045),
			wantError: false,
		},
		{
			json:      toJSON(1099),
			wantUnit:  monetary.Unit(1099),
			wantError: false,
		},
		{
			json:      toJSON(1),
			wantUnit:  monetary.Unit(1),
			wantError: false,
		},
		{
			json:      toJSON(0),
			wantUnit:  monetary.Unit(0),
			wantError: false,
		},
	}

	for _, tc := range testCases {
		var unit monetary.Unit
		err := json.Unmarshal(tc.json, &unit)

		if tc.wantError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tc.wantUnit, unit, tc.wantUnit.String())
		}
	}
}
