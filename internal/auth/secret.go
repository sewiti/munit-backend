package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"io/ioutil"
	"os"

	"github.com/apex/log"
)

const (
	secretLen      = 32
	secretFilePerm = os.FileMode(0600)
)

func LoadSecret(file string) error {
	secret, err := readSecret(file)
	if err != nil && !os.IsNotExist(err) {
		return err // error
	}
	if len(secret) != secretLen {
		// ErrNotFound or bad length
		if len(secret) != secretLen {
			log.WithField("file", file).Warn("bad secret length")
		}
		log.WithField("file", file).Info("creating secret")
		secret, err = createSecret(secretLen)
		if err != nil {
			return err
		}
		if err = writeSecret(file, secret); err != nil {
			log.WithError(err).WithField("file", file).Warn("unable to save secret")
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
	return ioutil.ReadAll(f)
}

func writeSecret(file string, secret []byte) error {
	return os.WriteFile(file, secret, secretFilePerm)
}

func createSecret(len int) ([]byte, error) {
	secret := make([]byte, secretLen)
	_, err := rand.Read(secret)
	return secret, err
}
