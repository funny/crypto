// Package aes256cbc is a helper to generate OpenSSL compatible encryption
// with autmatic IV derivation and storage. As long as the key is known all
// data can also get decrypted using OpenSSL CLI.
// Code from http://dequeue.blogspot.de/2014/11/decrypting-something-encrypted-with.html
package aes256cbc // import "github.com/funny/crypto/aes256cbc"

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// OpenSSL salt is always this string + 8 bytes of actual salt
var openSSLSaltHeader = []byte("Salted__")

type openSSLCreds [48]byte

// openSSLEvpBytesToKey follows the OpenSSL (undocumented?) convention for extracting the key and IV from passphrase.
// It uses the EVP_BytesToKey() method which is basically:
// D_i = HASH^count(D_(i-1) || password || salt) where || denotes concatentaion, until there are sufficient bytes available
// 48 bytes since we're expecting to handle AES-256, 32bytes for a key and 16bytes for the IV
func (c *openSSLCreds) Extract(password, salt []byte) (key, iv []byte) {
	m := c[:]
	buf := make([]byte, 0, 16+len(password)+len(salt))
	var prevSum [16]byte
	for i := 0; i < 3; i++ {
		n := 0
		if i > 0 {
			n = 16
		}
		buf = buf[:n+len(password)+len(salt)]
		copy(buf, prevSum[:])
		copy(buf[n:], password)
		copy(buf[n+len(password):], salt)
		prevSum = md5.Sum(buf)
		copy(m[i*16:], prevSum[:])
	}
	return c[:32], c[32:]
}

// DecryptString decrypts a string that was encrypted using OpenSSL and AES-256-CBC.
func DecryptString(passphrase, encryptedBase64String string) (string, error) {
	text, err := Decrypt([]byte(passphrase), []byte(encryptedBase64String))
	return string(text), err
}

// Decrypt decrypts a []byte that was encrypted using OpenSSL and AES-256-CBC.
func Decrypt(passphrase, encryptedBase64 []byte) ([]byte, error) {
	data := make([]byte, base64.StdEncoding.DecodedLen(len(encryptedBase64)))
	n, err := base64.StdEncoding.Decode(data, encryptedBase64)
	if err != nil {
		return nil, err
	}
	saltHeader := data[:aes.BlockSize]
	if !bytes.Equal(saltHeader[:8], openSSLSaltHeader) {
		return nil, fmt.Errorf("Does not appear to have been encrypted with OpenSSL, salt header missing.")
	}
	var creds openSSLCreds
	key, iv := creds.Extract(passphrase, saltHeader[8:])
	return decrypt(key, iv, data[:n])
}

func decrypt(key, iv, data []byte) ([]byte, error) {
	if len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("bad blocksize(%v), aes.BlockSize = %v\n", len(data), aes.BlockSize)
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCDecrypter(c, iv)
	cbc.CryptBlocks(data[aes.BlockSize:], data[aes.BlockSize:])
	out, err := pkcs7Unpad(data[aes.BlockSize:])
	if out == nil {
		return nil, err
	}
	return out, nil
}

// EncryptString encrypts a string in a manner compatible to OpenSSL encryption
// functions using AES-256-CBC as encryption algorithm
func EncryptString(passphrase, plaintextString string) (string, error) {
	text, err := Encrypt([]byte(passphrase), []byte(plaintextString))
	return string(text), err
}

// Encrypt encrypts a []byte in a manner compatible to OpenSSL encryption
// functions using AES-256-CBC as encryption algorithm
func Encrypt(passphrase, plaintext []byte) ([]byte, error) {
	var salt [8]byte // Generate an 8 byte salt
	_, err := io.ReadFull(rand.Reader, salt[:])
	if err != nil {
		return nil, err
	}

	data := make([]byte, len(plaintext)+aes.BlockSize)
	copy(data[0:], openSSLSaltHeader)
	copy(data[8:], salt[:])
	copy(data[aes.BlockSize:], plaintext)

	var creds openSSLCreds
	key, iv := creds.Extract(passphrase, salt[:])
	enc, err := encrypt(key, iv, data)
	if err != nil {
		return nil, err
	}

	encryptedBase64 := make([]byte, base64.StdEncoding.EncodedLen(len(enc)))
	base64.StdEncoding.Encode(encryptedBase64, enc)
	return encryptedBase64, nil
}

func encrypt(key, iv, data []byte) ([]byte, error) {
	padded, err := pkcs7Pad(data)
	if err != nil {
		return nil, err
	}
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	cbc := cipher.NewCBCEncrypter(c, iv)
	cbc.CryptBlocks(padded[aes.BlockSize:], padded[aes.BlockSize:])
	return padded, nil
}

var padPatterns [aes.BlockSize][]byte

func init() {
	for i := 0; i < len(padPatterns); i++ {
		padPatterns[i] = bytes.Repeat([]byte{byte(i)}, i)
	}
}

// pkcs7Pad appends padding.
func pkcs7Pad(data []byte) ([]byte, error) {
	padlen := 1
	for ((len(data) + padlen) % aes.BlockSize) != 0 {
		padlen = padlen + 1
	}
	return append(data, padPatterns[padlen]...), nil
}

// pkcs7Unpad returns slice of the original data without padding.
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data)%aes.BlockSize != 0 || len(data) == 0 {
		return nil, fmt.Errorf("invalid data len %d", len(data))
	}
	padlen := int(data[len(data)-1])
	if padlen > aes.BlockSize || padlen == 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	if !bytes.Equal(padPatterns[padlen], data[len(data)-padlen:]) {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-padlen], nil
}
