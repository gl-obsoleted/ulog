package ushare

import (
	"crypto/rand"
)

const GSaltLength = 128

type CryptoSalt []byte

func NewSalt() (CryptoSalt, error) {
	saltBuf := make([]byte, GSaltLength)
	if _, err := rand.Read(saltBuf); err != nil {
		return nil, err
	}

	return saltBuf, nil
}
