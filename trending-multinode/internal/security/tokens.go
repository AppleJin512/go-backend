package security

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/sirupsen/logrus"
)

func TokenGenerate(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func TokenMustGenerate(n int) string {
	token, err := TokenGenerate(n)
	if err != nil {
		logrus.Warnf("token generate error: %v", err)
	}
	return token
}
