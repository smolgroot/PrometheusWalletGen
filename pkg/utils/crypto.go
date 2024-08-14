package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

// Crypto is an interface that defines methods for encryption and decryption.
type Crypto interface {
	Encrypt(string) string
	Decrypt(string) string
}

// XOR is a struct that represents the XOR encryption algorithm.
type XOR struct {
	K string // The encryption key for XOR algorithm.
}

// AES256 is a struct that represents the AES-256 encryption algorithm.
type AES256 struct {
	K string // The encryption key for AES-256 algorithm.
}

// Shasha20 is a struct that represents the Shasha-20 encryption algorithm.
type Shasha20 struct {
	K string // The encryption key for Shasha-20 algorithm.
}

func (a AES256) Encrypt(s string) string {
	config, err := LoadConfiguration("./config.json")
	Check(err)

	key, _ := hex.DecodeString(config.MasterKey)
	ct := []byte(s)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	Check(err)

	//Create a new GCM block cipher for AES
	aesGCM, err := cipher.NewGCM(block)
	Check(err)

	//Create a nonce. It should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatal(err)
	}

	//Encrypt the data using aesGCM.Seal
	et := aesGCM.Seal(nonce, nonce, ct, nil)
	return fmt.Sprintf("%x", et)
}

func (a AES256) Decrypt(s string) string {
	config, err := LoadConfiguration("./config.json")
	Check(err)
	key, _ := hex.DecodeString(config.MasterKey)
	enc, _ := hex.DecodeString(s)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	Check(err)

	//Create a new GCM block cipher for AES
	aesGCM, err := cipher.NewGCM(block)
	Check(err)

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	ct, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}

	return string(ct)
}

// Encrypt encrypts the given string using the XOR algorithm.
func (x XOR) Encrypt(s string) string {
	var o string
	for i := 0; i < len(s); i++ {
		o += string(s[i] ^ x.K[i%len(x.K)])
	}
	return o
}

func (x XOR) Decrypt(s string) string {
	return x.Encrypt(s)
}
