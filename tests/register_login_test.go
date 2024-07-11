package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/mock/gomock"
	sso "github.com/kuromii5/miku-notes-auth/generated"
	"github.com/kuromii5/miku-notes-auth/internal/auth"
	"github.com/kuromii5/miku-notes-auth/internal/service"
	"github.com/kuromii5/miku-notes-auth/tests/suite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterLoginHappyPath(t *testing.T) {
	// ----------------------------------------- REGISTER AND LOGIN TEST ---------------------------------
	ctx, st := suite.NewSuite(t)
	assert := assert.New(st.T)
	require := require.New(st.T)

	// Generate fake email and password
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, 8)
	fingerprint := "fingerprint"

	// Call Register method
	authResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:       email,
		Password:    pass,
		Fingerprint: fingerprint,
	})
	require.NoError(err)

	// ----------------------------------------- TOKEN VALIDATION ----------------------------------------

	// Get access token
	accessToken := authResponse.GetAccessToken()
	require.NotEmpty(accessToken)

	// Validate access token using token manager instance
	userIDFromAccess, err := st.TokenManager.ValidateAccessToken(ctx, accessToken)
	require.NoError(err)

	// get current time and JWT claims
	currentTime := time.Now()
	JWT, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(st.Cfg.Tokens.Secret), nil
	})
	require.NoError(err)

	claims, ok := JWT.Claims.(jwt.MapClaims)
	require.True(ok)

	// Validate the claims
	// id
	subject, err := claims.GetSubject()
	require.NoError(err)
	parsedUserID, err := strconv.Atoi(subject)
	require.NoError(err)
	assert.Equal(userIDFromAccess, int32(parsedUserID))

	// expiry time
	expiredAt, err := claims.GetExpirationTime()
	require.NoError(err)
	assert.InDelta(currentTime.Add(st.Cfg.Tokens.AccessTTL).Unix(), float64(expiredAt.Unix()), float64(time.Second*5)) // 5 sec error range

	// Get refresh token
	refreshToken := authResponse.GetRefreshToken()
	require.NotEmpty(refreshToken)

	// Setup mock expectation for ValidateRefreshToken
	mockUserGetter := st.Mocks.UserGetter
	mockUserGetter.EXPECT().UserID(gomock.Any(), refreshToken, fingerprint).Return(subject, nil)

	// Validate refresh token using token manager instance
	userIDFromRefresh, err := st.TokenManager.ValidateRefreshToken(ctx, refreshToken, fingerprint)
	require.NoError(err)
	assert.Equal(userIDFromRefresh, int32(parsedUserID))
}

func TestRegisterLogin_DoubleRegistration(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	assert := assert.New(st.T)
	require := require.New(st.T)

	// Generate fake email and password
	email := gofakeit.Email()
	pass := gofakeit.Password(true, true, true, true, false, 8)
	fingerprint := "fingerprint"

	registerResponse, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:       email,
		Password:    pass,
		Fingerprint: fingerprint,
	})
	require.NoError(err)
	assert.NotEmpty(registerResponse.GetAccessToken())
	assert.NotEmpty(registerResponse.GetRefreshToken())

	// register with same credentials for the second time
	registerResponse, err = st.AuthClient.Register(ctx, &sso.RegisterRequest{
		Email:       email,
		Password:    pass,
		Fingerprint: fingerprint,
	})
	require.Error(err)
	assert.Empty(registerResponse.GetAccessToken())
	assert.Empty(registerResponse.GetRefreshToken())
	assert.ErrorContains(err, "user already exists")
}

func TestRegisterLogin_FailCases(t *testing.T) {
	ctx, st := suite.NewSuite(t)
	require := require.New(st.T)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Invalid Email Format",
			email:       "invalid-email",
			password:    "password123",
			expectedErr: auth.ErrInvalidEmail.Error(),
		},
		{
			name:        "Short Password",
			email:       "test@example.com",
			password:    "short",
			expectedErr: auth.ErrShortPassword.Error(),
		},
		{
			name:        "Long Email",
			email:       fmt.Sprintf("%s@example.com", string(make([]byte, 255))),
			password:    "password123",
			expectedErr: auth.ErrInvalidEmail.Error(),
		},
		{
			name:        "Long Password",
			email:       "test@example.com",
			password:    string(make([]byte, 65)),
			expectedErr: auth.ErrLongPassword.Error(),
		},
		{
			name:        "Missing Email",
			email:       "",
			password:    "password123",
			expectedErr: auth.ErrRequired.Error(),
		},
		{
			name:        "Missing Password",
			email:       "test@example.com",
			password:    "",
			expectedErr: auth.ErrRequired.Error(),
		},
		{
			name:        "Invalid Login Credentials",
			email:       "nonexistent@example.com",
			password:    "wrongpassword",
			expectedErr: service.ErrUserExists.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &sso.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(err)
			require.Contains(err.Error(), tt.expectedErr)
		})
	}

	// Test invalid login credentials
	_, err := st.AuthClient.Login(ctx, &sso.LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "dumbass",
	})
	require.Error(err)
	require.ErrorContains(err, service.ErrInvalidCreds.Error())
}
