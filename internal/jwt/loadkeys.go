package jwt_keys

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type JWTKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var jwtKeys = &JWTKeys{}
var isKeysLoaded = false

func loadPrivateKey(keyPath string) {
	privateKeyBytes, err := os.ReadFile(keyPath)
	errors := exceptions.Exception{}
	errors.Tag = enums.INTERNAL_ERROR
	if err != nil {
		errors.AddError(enums.MISSING_FILE)
		panic(errors)
	}
	jwtKeys.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		errors.AddError(enums.CAST_ERROR)
		panic(errors)
	}
}

func loadPublicKey(keyPath string) {
	publicKeyBytes, err := os.ReadFile(keyPath)
	errors := exceptions.Exception{}
	errors.Tag = enums.INTERNAL_ERROR
	if err != nil {
		errors.AddError(enums.MISSING_FILE)
		panic(err)
	}
	jwtKeys.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		errors.AddError(enums.CAST_ERROR)
		panic(err)
	}
}

func SetupJWTKeys(jwtKeysPath string) {
	if !isKeysLoaded {
		privateKeyPath := jwtKeysPath + "/privateKey.pem"
		loadPrivateKey(privateKeyPath)
		publicKeyPath := jwtKeysPath + "/publicKey.pem"
		loadPublicKey(publicKeyPath)
		isKeysLoaded = true
	}
}

func GetJWTKeys() *JWTKeys {
	if jwtKeys.PrivateKey == nil || jwtKeys.PublicKey == nil {
		panic("JWT keys are not loaded")
	}
	return jwtKeys
}
