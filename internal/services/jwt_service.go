package services

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type JWT struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewJWT(Const *bootstrap.Constants) *JWT {
	return &JWT{
		PrivateKey: loadPrivateKey(Const.JWTKeysPath + "/privateKey.pem"),
		PublicKey:  loadPublicKey(Const.JWTKeysPath + "/publicKey.pem"),
	}
}

func loadPrivateKey(keyPath string) *rsa.PrivateKey {
	privateKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}
	return privateKey
}

func loadPublicKey(keyPath string) *rsa.PublicKey {
	publicKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}
	return publicKey
}

func (j *JWT) GenerateToken(userID int) (string, string) {
	accessTokenClaims := jwt.MapClaims{
		"iss": "bidlancer",
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
		"iat": time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(j.PrivateKey)
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
	refreshTokenString, err := refreshToken.SignedString(j.PrivateKey)
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

func (j *JWT) VerifyToken(tokenString string) jwt.MapClaims {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			panic(fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return j.PublicKey, nil
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
