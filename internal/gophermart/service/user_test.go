package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/sergeizaitcev/gophermart/internal/gophermart/service"
	mocks_storage "github.com/sergeizaitcev/gophermart/internal/gophermart/service/mocks"
	"github.com/sergeizaitcev/gophermart/pkg/randutil"
	"github.com/sergeizaitcev/gophermart/pkg/sign"
)

var _ sign.Signer = (*signerStub)(nil)

type signerStub struct {
	value string
	err   error
}

func (s signerStub) Sign(payload string) (token string, err error) {
	return s.value, s.err
}

func (s signerStub) Parse(token string) (payload string, err error) {
	return s.value, s.err
}

var tUser = struct {
	id    uuid.UUID
	login string
	pass  string
}{
	id:    uuid.New(),
	login: "user",
	pass:  "password",
}

func TestAuth(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	mockStorage := mocks_storage.NewMockAuth(ctrl)

	t.Run("signup", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			t.Parallel()

			signer := signerStub{value: "token"}
			auth := service.NewAuth(mockStorage, signer)

			mockStorage.EXPECT().
				CreateUser(ctx, tUser.login, tUser.pass).
				Return(tUser.id, (error)(nil))

			token, err := auth.SignUp(ctx, tUser.login, tUser.pass)
			require.NoError(t, err)
			require.Equal(t, signer.value, token)
		})

		t.Run("error", func(t *testing.T) {
			t.Parallel()

			signer := signerStub{err: errors.New("error")}
			auth := service.NewAuth(mockStorage, signer)

			mockStorage.EXPECT().
				CreateUser(ctx, tUser.login, tUser.pass).
				Return(uuid.UUID{}, errors.New("error"))

			_, err := auth.SignUp(ctx, tUser.login, tUser.pass)
			require.Error(t, err)

			mockStorage.EXPECT().
				CreateUser(ctx, tUser.login, tUser.pass).
				Return(tUser.id, (error)(nil))

			_, err = auth.SignUp(ctx, tUser.login, tUser.pass)
			require.Error(t, err)
		})

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			auth := service.NewAuth(nil, nil)

			testCases := []struct {
				login string
				pass  string
			}{
				{"", randutil.String(10)},
				{randutil.String(10), ""},
				{randutil.String(10), randutil.String(100)},
			}

			for _, tc := range testCases {
				_, err := auth.SignUp(ctx, tc.login, tc.pass)
				require.Error(t, err)
			}
		})
	})

	t.Run("signin", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			t.Parallel()

			signer := signerStub{value: "token"}
			auth := service.NewAuth(mockStorage, signer)

			mockStorage.EXPECT().
				GetUser(ctx, tUser.login, tUser.pass).
				Return(tUser.id, (error)(nil))

			token, err := auth.SignIn(ctx, tUser.login, tUser.pass)
			require.NoError(t, err)
			require.Equal(t, signer.value, token)
		})

		t.Run("error", func(t *testing.T) {
			t.Parallel()

			signer := signerStub{err: errors.New("error")}
			auth := service.NewAuth(mockStorage, signer)

			mockStorage.EXPECT().
				GetUser(ctx, tUser.login, tUser.pass).
				Return(uuid.UUID{}, errors.New("error"))

			_, err := auth.SignIn(ctx, tUser.login, tUser.pass)
			require.Error(t, err)

			mockStorage.EXPECT().
				GetUser(ctx, tUser.login, tUser.pass).
				Return(tUser.id, (error)(nil))

			_, err = auth.SignIn(ctx, tUser.login, tUser.pass)
			require.Error(t, err)
		})

		t.Run("invalid", func(t *testing.T) {
			t.Parallel()

			auth := service.NewAuth(nil, nil)

			testCases := []struct {
				login string
				pass  string
			}{
				{"", randutil.String(10)},
				{randutil.String(10), ""},
				{randutil.String(10), randutil.String(100)},
			}

			for _, tc := range testCases {
				_, err := auth.SignIn(ctx, tc.login, tc.pass)
				require.Error(t, err)
			}
		})
	})

	t.Run("verify", func(t *testing.T) {
		t.Run("ok", func(t *testing.T) {
			t.Parallel()

			signer := signerStub{value: tUser.id.String()}
			auth := service.NewAuth(nil, signer)

			id, err := auth.Verify(ctx, "token")
			require.NoError(t, err)
			require.Equal(t, tUser.id.String(), id)
		})

		t.Run("expire", func(t *testing.T) {
			t.Parallel()

			signer := signerStub{err: sign.ErrTokenExpired}
			auth := service.NewAuth(nil, signer)

			_, err := auth.Verify(ctx, "token")
			require.ErrorIs(t, err, sign.ErrTokenExpired)
		})

		t.Run("error", func(t *testing.T) {
			t.Parallel()

			signer := signerStub{err: errors.New("error")}
			auth := service.NewAuth(nil, signer)

			_, err := auth.Verify(ctx, "token")
			require.Error(t, err)
		})
	})
}
