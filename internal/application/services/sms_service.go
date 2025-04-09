package services

type SmsService interface {
	GenerateOTP() string
	SendOTP(phonenumber string, code string)
}
