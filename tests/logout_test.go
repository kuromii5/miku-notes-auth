package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	sso "github.com/kuromii5/miku-notes-auth/generated"
	"github.com/kuromii5/miku-notes-auth/tests/suite"
	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	require := require.New(st.T)

	// Register and log in to get access token
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, 8)
	fingerprint := "fingerprint"

	registerResp, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:       email,
		Password:    pass,
		Fingerprint: fingerprint,
	})
	require.NoError(err)
	require.NotEmpty(registerResp.GetAccessToken())
	require.NotEmpty(registerResp.GetRefreshToken())

	// Call Logout method with the access token
	_, err = st.AuthClient.Logout(ctx, &sso.LogoutRequest{
		AccessToken: registerResp.GetAccessToken(),
		Fingerprint: fingerprint,
	})
	require.NoError(err)
}

func TestLogout_Fail(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	require := require.New(st.T)

	// Call Logout with invalid access token
	_, err := st.AuthClient.Logout(ctx, &sso.LogoutRequest{
		AccessToken: "invalid-access-token",
		Fingerprint: "fingerprint",
	})
	require.Error(err)
}
