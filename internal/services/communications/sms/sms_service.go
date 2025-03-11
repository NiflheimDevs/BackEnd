package sms

import (
	"fmt"
	"math/rand"
)

func GenerateOTP() string {
	return fmt.Sprintf("%05d", rand.Int()%100000)

}

// TODO: connect to actual service
func SendOTP(phonenumber string, code string) {

}
