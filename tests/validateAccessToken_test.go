package tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	sso "github.com/kuromii5/miku-notes-auth/generated"
	"github.com/kuromii5/miku-notes-auth/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateAccessToken(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	assert := assert.New(st.T)
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

	// Call ValidateAccessToken method with the access token
	validateATResp, err := st.AuthClient.ValidateAccessToken(ctx, &sso.ValidateATRequest{
		AccessToken: registerResp.GetAccessToken(),
	})
	require.NoError(err)
	assert.NotZero(validateATResp.GetUserId())
}

func TestValidateAccessToken_Fail(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	require := require.New(st.T)

	// Call ValidateAccessToken with invalid access token
	_, err := st.AuthClient.ValidateAccessToken(ctx, &sso.ValidateATRequest{
		AccessToken: "invalid-access-token",
	})
	require.Error(err)
}
