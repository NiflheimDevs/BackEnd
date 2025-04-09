package pkg

import "golang.org/x/crypto/bcrypt"

type SecretSauce struct {
}

func NewSecretSauce() *SecretSauce {
	return &SecretSauce{}
}

func (ss *SecretSauce) MakeSauce(ingredient string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(ingredient), bcrypt.DefaultCost)
}

func (ss *SecretSauce) SauceReferee(sauce []byte, ingredient string) error {
	return bcrypt.CompareHashAndPassword(sauce, []byte(ingredient))
}
