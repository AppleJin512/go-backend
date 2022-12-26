package security

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func PasswordGenerate(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func PasswordMustGenerate(password string) string {
	pwHash, err := PasswordGenerate(password)
	if err != nil {
		logrus.Warnf("password generate error: %v", err)
	}
	return pwHash
}

func PasswordValidate(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
