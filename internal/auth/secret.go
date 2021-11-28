package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"os"

	"github.com/apex/log"
)

const (
	secretLen      = ed25519.SeedSize
	secretFilePerm = os.FileMode(0600)
)

func LoadSecret(file string) error {
	secret, err := readSecret(file)
	if err != nil {
		if !os.IsNotExist(err) {
			log.WithError(err).WithField("file", file).Warn("unable to read secret")
		}
	} else if len(secret) != secretLen {
		log.WithField("file", file).Warnf("secret invalid length: got %d, expected %d", len(secret), secretLen)
	}

	if err != nil || len(secret) != secretLen {
		log.WithField("file", file).Info("creating secret")
		secret, err = createSecret(secretLen)
		if err != nil {
			return err // will log.Fatal outside
		}
		if err = writeSecret(file, secret); err != nil { // non fatal
			log.WithError(err).WithField("file", file).Warn("unable to write secret")
		}
	}

	secretKey = ed25519.NewKeyFromSeed(secret)
	return nil
}

func readSecret(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Mode().Perm() != secretFilePerm {
		log.WithField("file", file).Warnf("secret file perm not %v", secretFilePerm)
	}

	enc := base64.NewDecoder(base64.StdEncoding, f)
	secret := make([]byte, secretLen)
	_, err = enc.Read(secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func writeSecret(file string, secret []byte) error {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, secretFilePerm)
	if err != nil {
		return err
	}
	enc := base64.NewEncoder(base64.RawStdEncoding, f)
	_, err = enc.Write(secret)
	if errCl := enc.Close(); errCl != nil && err == nil {
		err = errCl
	}
	if errCl := f.Close(); errCl != nil && err == nil {
		err = errCl
	}
	return err
}

func createSecret(len int) ([]byte, error) {
	secret := make([]byte, secretLen)
	_, err := rand.Read(secret)
	return secret, err
}
