package auth

import (
	"bytes"
	crand "crypto/rand"
	"io"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64MiB
	argonThreads = 4
	argonKeyLen  = 32

	tokenExpiry = 7 * 24 * time.Hour // 7d
)

func MakeSalt(rand io.Reader) ([]byte, error) {
	salt := make([]byte, argonKeyLen)
	_, err := crand.Read(salt)
	return salt, err
}

func HashPasswd(password, salt []byte) []byte {
	return argon2.IDKey(password, salt, argonTime, argonMemory, argonThreads, argonKeyLen)
}

func VerifyPasswd(passwdHash, password, salt []byte) bool {
	return bytes.Equal(passwdHash, HashPasswd(password, salt))
}
