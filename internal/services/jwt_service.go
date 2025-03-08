package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	jwt_keys "github.com/niflheimdevs/backend/internal/jwt"
)

type JWTToken struct{}

func NewJWTToken() *JWTToken {
	return &JWTToken{}
}

func (jt *JWTToken) GenerateToken(userID int) (string, string) {
	jwtKeys := jwt_keys.GetJWTKeys()

	accessTokenClaims := jwt.MapClaims{
		"iss": "bidlancer",
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
		"iat": time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtKeys.PrivateKey)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.AUTH_GENERATE_TOKEN_ERROR,
			},
		})
	}

	refreshTokenClaims := jwt.MapClaims{
		"iss": "bidlancer",
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKeys.PrivateKey)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.AUTH_GENERATE_TOKEN_ERROR,
			},
		})
	}

	return accessTokenString, refreshTokenString
}

func (jt *JWTToken) VerifyToken(tokenString string) jwt.MapClaims {
	jwtKeys := jwt_keys.GetJWTKeys()
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return jwtKeys.PublicKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	}
	panic(exceptions.Exception{
		Tag: enums.VALIDATION_ERROR,
		Errors: []enums.SpecificError{
			enums.AUTH_ACCESS_DENIED,
		},
	})
}
