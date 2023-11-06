package sign_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/pkg/randutil"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

func TestSigner(t *testing.T) {
	testCases := []struct {
		name      string
		secret    []byte
		options   []sign.Option
		payload   string
		wantError bool
	}{
		{
			name:      "unlim",
			secret:    randutil.Bytes(20),
			payload:   randutil.String(128),
			wantError: false,
		},
		{
			name:      "ttl",
			secret:    randutil.Bytes(20),
			options:   []sign.Option{sign.WithTTL(time.Minute)},
			payload:   randutil.String(128),
			wantError: false,
		},
		{
			name:      "expired",
			secret:    randutil.Bytes(20),
			options:   []sign.Option{sign.WithTTL(time.Nanosecond)},
			payload:   randutil.String(128),
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			signer := sign.New(tc.secret, tc.options...)

			token, err := signer.Sign(tc.payload)
			require.NoError(t, err)

			payload, err := signer.Parse(token)

			if tc.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.payload, payload)
			}
		})
	}
}
