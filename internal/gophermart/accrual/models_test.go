package accrual_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/accrual"
)

func TestOrderInfo_MarshalJSON(t *testing.T) {
	testCases := []struct {
		info accrual.OrderInfo
		want string
	}{
		{
			info: accrual.OrderInfo{},
			want: `{"order":"","status":"UNKNOWN","accrual":0}`,
		},
		{
			info: accrual.OrderInfo{
				Order:   "1234567890",
				Status:  accrual.StatusRegistered,
				Accrual: 12303,
			},
			want: `{"order":"1234567890","status":"REGISTERED","accrual":123.03}`,
		},
	}

	for _, tc := range testCases {
		b, err := json.Marshal(tc.info)
		require.NoError(t, err)
		require.Equal(t, tc.want, string(b))
	}
}

func TestOrderInfo_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		json      string
		wantInfo  accrual.OrderInfo
		wantError bool
	}{
		{
			json: `{"order":"1234567890","status":"REGISTERED","accrual":123.03}`,
			wantInfo: accrual.OrderInfo{
				Order:   "1234567890",
				Status:  accrual.StatusRegistered,
				Accrual: 12303,
			},
			wantError: false,
		},
		{
			json:      `{}`,
			wantInfo:  accrual.OrderInfo{},
			wantError: false,
		},
		{
			json:      `{"status":1}`,
			wantInfo:  accrual.OrderInfo{},
			wantError: true,
		},
		{
			json:      `{"status":"SOMESTATUS"}`,
			wantInfo:  accrual.OrderInfo{},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		var got accrual.OrderInfo
		err := json.Unmarshal([]byte(tc.json), &got)

		if tc.wantError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tc.wantInfo, got)
		}
	}
}
