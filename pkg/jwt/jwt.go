package jwt_client

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	_defaultAccessDuration  = time.Hour
	_defaultRefreshDuration = time.Hour * 24 * 7
)

type TokensDetails struct {
	Access     string
	AccessExp  int64
	Refresh    string
	RefreshExp int64
	UserId     int
}

type JWTClient struct {
	secret          []byte
	accessDuration  time.Duration
	refreshDuration time.Duration
}

func New(secret string) *JWTClient {
	return NewWithDurations(
		secret,
		_defaultAccessDuration,
		_defaultRefreshDuration,
	)
}

func NewWithDurations(secret string, accessDuration, refreshDuration time.Duration) *JWTClient {
	return &JWTClient{
		secret:          []byte(secret),
		accessDuration:  accessDuration,
		refreshDuration: refreshDuration,
	}
}

func (j *JWTClient) GenerateTokenString(tokenType string, exp int64, userId int) (string, error) {
	claims := jwt.MapClaims{}
	claims["token_type"] = tokenType
	claims["exp"] = exp
	claims["user_id"] = userId
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTClient) GenerateTokenDetails(userId int) (TokensDetails, error) {
	td := TokensDetails{
		AccessExp:  time.Now().Add(j.accessDuration).Unix(),
		RefreshExp: time.Now().Add(j.refreshDuration).Unix(),
		UserId:     userId,
	}

	// Access Token
	access, err := j.GenerateTokenString("access", td.AccessExp, td.UserId)
	if err != nil {
		return TokensDetails{}, err
	}
	td.Access = access

	// Refresh Token
	refresh, err := j.GenerateTokenString("refresh", td.RefreshExp, td.UserId)
	if err != nil {
		return TokensDetails{}, err
	}
	td.Refresh = refresh

	return td, nil
}

func (j *JWTClient) VerifyToken(t string) error {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return j.secret, nil
	})

	if err != nil {
		return err
	}

	clm, ok := token.Claims.(jwt.Claims)
	if !ok && !token.Valid {
		fmt.Println(clm)
		return nil
	}

	return nil
}
