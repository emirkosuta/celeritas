package jwt

import (
	"fmt"
	"time"

	"github.com/emirkosuta/celeritas/jwt/dto"
	"github.com/emirkosuta/celeritas/rsa"
	"github.com/golang-jwt/jwt"
)

const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

type JwtToken struct {
	JwtTokenTimeExp        time.Duration
	JwtRefreshTokenTimeExp time.Duration
	RSAPrivate             string
	RSAPublic              string
}

func (t *JwtToken) createAndSignToken(tokenType string, iat int64, claims jwt.MapClaims, expirationTime time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	claims["exp"] = expirationTime.Unix()
	claims["token_type"] = tokenType
	claims["iat"] = iat
	token.Claims = claims

	privateRsa, err := rsa.ReadPrivateKeyFromEnv(t.RSAPrivate)
	if err != nil {
		return "", fmt.Errorf("error reading private RSA key from environment: %w", err)
	}
	tokenString, err := token.SignedString(privateRsa)
	if err != nil {
		return "", fmt.Errorf("error signing the %s: %w", tokenType, err)
	}
	return tokenString, nil
}

// Sign is a method for generating JWT tokens and refresh tokens. This method accepts a map of claims
// that will be included in the token. The following claims are added or modified within the method:
// 'exp' for the expiration time of the token,
// 'iat' for the issuance time of the token,
// 'token_type' indicating whether the token is an access token or a refresh token.
//
// The method requires the 'id' claim to be present in the input. If it's not, the method returns an error.
// If the 'exp' claim is not provided in the input, the method sets it to the current time plus JWT_ACCESS_TOKEN_EXPIRY env variable for the access token
// and JWT_REFRESH_TOKEN_EXPIRY env variable for the refresh token.
//
// The generated JWT tokens are signed using the RS256 algorithm and the RSA private key read from the environment.
//
// The method returns a Token object that includes the access token and refresh token, both of type 'Bearer'.
// If any error occurs during the process, it is returned by the method.
func (t *JwtToken) Sign(claims jwt.MapClaims) (*dto.Token, error) {

	if claims["id"] == nil {
		return nil, fmt.Errorf("missing 'id' claim")
	}
	timeNow := time.Now()

	// Create auth token
	authTokenString, err := t.createAndSignToken(AccessToken, timeNow.Unix(), claims, timeNow.Add(t.JwtTokenTimeExp))
	if err != nil {
		return nil, err
	}

	// Create refresh token
	refreshTokenString, err := t.createAndSignToken(RefreshToken, timeNow.Unix(), claims, timeNow.Add(t.JwtRefreshTokenTimeExp))
	if err != nil {
		return nil, err
	}

	return &dto.Token{
		Type:  "Bearer",
		Token: authTokenString,
		RefreshToken: dto.RefreshToken{
			Value: refreshTokenString,
			Iat:   fmt.Sprintf("%d", timeNow.Unix()),
		},
	}, nil
}
