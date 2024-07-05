package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	sso "github.com/kuromii5/sso-auth/generated"
	"github.com/kuromii5/sso-auth/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccessToken(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	assert := assert.New(st.T)
	require := require.New(st.T)

	// Register and log in to get refresh token
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

	// Call GetAccessToken method with the refresh token
	getATResp, err := st.AuthClient.GetAccessToken(ctx, &sso.GetATRequest{
		RefreshToken: registerResp.GetRefreshToken(),
		Fingerprint:  fingerprint,
	})
	require.NoError(err)
	assert.NotEmpty(getATResp.GetAccessToken())
}

func TestGetAccessToken_Fail(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	require := require.New(st.T)

	// Call GetAccessToken with invalid refresh token
	_, err := st.AuthClient.GetAccessToken(ctx, &sso.GetATRequest{
		RefreshToken: "invalid-refresh-token",
		Fingerprint:  "fingerprint",
	})
	require.Error(err)
}
