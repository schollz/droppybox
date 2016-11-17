package gojot

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
)

// decryptString returns the decrypted string using a passphrase and
// GPG symmetric encryption
func DecryptString(decryptionString string, encryptionPassphraseString string) (string, error) {
	encryptionPassphrase := []byte(encryptionPassphraseString)
	decbuf := bytes.NewBuffer([]byte(decryptionString))
	result, err := armor.Decode(decbuf)
	if err != nil {
		return "", err
	}

	alreadyPrompted := false
	md, err := openpgp.ReadMessage(result.Body, nil, func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		if alreadyPrompted {
			return nil, errors.New("Could not decrypt using passphrase")
		} else {
			alreadyPrompted = true
		}
		return encryptionPassphrase, nil
	}, nil)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// decryptString returns the encrypted string using a passphrase and
// GPG symmetric encryption
func EncryptString(encryptionText string, encryptionPassphraseString string) string {
	encryptionPassphrase := []byte(encryptionPassphraseString)
	encbuf := bytes.NewBuffer(nil)
	w, err := armor.Encode(encbuf, "PGP SIGNATURE", nil)
	if err != nil {
		logger.Error("Error encrypting:" + err.Error())
	}

	plaintext, err := openpgp.SymmetricallyEncrypt(w, encryptionPassphrase, nil, nil)
	if err != nil {
		logger.Error("Error opengpg encrypting:" + err.Error())
	}
	message := []byte(encryptionText)
	_, err = plaintext.Write(message)

	plaintext.Close()
	w.Close()
	return strings.TrimSpace(encbuf.String())
}

// DecryptFile returns the decrypted contents of a GPG symmetric encrypted file
func DecryptFile(file string, passphrase string) error {
	fileContents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	decrypted, err := DecryptString(string(fileContents), passphrase)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, []byte(decrypted), 0644)
	return err
}

// EncryptFile creates an encrypted file with extension gpg
// and shreds old file
func EncryptFile(file string, passphrase string) error {
	logger.Debug("Encrypting %s", file)
	fileContents, _ := ioutil.ReadFile(file)
	encrypted := EncryptString(string(fileContents), passphrase)
	err := ioutil.WriteFile(file, []byte(encrypted), 0644)
	return err
}
