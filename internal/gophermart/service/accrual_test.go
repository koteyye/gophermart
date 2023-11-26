package service_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
)

func TestOrderInfo_MarshalJSON(t *testing.T) {
	testCases := []struct {
		info service.AccrualOrderInfo
		want string
	}{
		{
			info: service.AccrualOrderInfo{},
			want: `{"order":"","status":"UNKNOWN","accrual":0}`,
		},
		{
			info: service.AccrualOrderInfo{
				Order:   "1234567890",
				Status:  service.AccrualOrderStatusRegistered,
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
		wantInfo  service.AccrualOrderInfo
		wantError bool
	}{
		{
			json: `{"order":"1234567890","status":"REGISTERED","accrual":123.03}`,
			wantInfo: service.AccrualOrderInfo{
				Order:   "1234567890",
				Status:  service.AccrualOrderStatusRegistered,
				Accrual: 12303,
			},
			wantError: false,
		},
		{
			json:      `{}`,
			wantInfo:  service.AccrualOrderInfo{},
			wantError: false,
		},
		{
			json:      `{"status":1}`,
			wantInfo:  service.AccrualOrderInfo{},
			wantError: true,
		},
		{
			json:      `{"status":"SOMESTATUS"}`,
			wantInfo:  service.AccrualOrderInfo{},
			wantError: false,
		},
	}

	for _, tc := range testCases {
		var got service.AccrualOrderInfo
		err := json.Unmarshal([]byte(tc.json), &got)

		if tc.wantError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
			require.Equal(t, tc.wantInfo, got)
		}
	}
}
