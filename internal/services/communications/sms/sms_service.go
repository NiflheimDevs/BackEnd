package sms

import (
	"fmt"
	"math/rand"
)

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Int()%1000000)

}

func SendOTP(phonenumber string) {

}
