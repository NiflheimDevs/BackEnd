package jwt

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JWTKeys struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

var jwtKeys = &JWTKeys{}
var isKeysLoaded = false

func loadPrivateKey(keyPath string) {
	privateKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	jwtKeys.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}
}

func loadPublicKey(keyPath string) {
	publicKeyBytes, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}
	jwtKeys.PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
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
