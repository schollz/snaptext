package server

import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	math_rand "math/rand"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func sha256sum(s string) string {
	h := sha256.New()
	h.Write([]byte("snaptext salt"))
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func encryptBytes(plaintext []byte, passphrase []byte) (encrypted []byte, err error) {
	encryptedS, salt, iv := encrypt(plaintext, passphrase)
	data := []string{encryptedS, salt, iv}
	encrypted, err = json.Marshal(data)
	return
}

func decryptBytes(encrypted []byte, passphrase []byte) (decrypted []byte, err error) {
	var data []string
	err = json.Unmarshal(encrypted, &data)
	if err != nil {
		return
	}
	if len(data) != 3 {
		err = errors.New("data corrupted")
		return
	}
	plaintext, err := hex.DecodeString(data[0])
	if err != nil {
		return
	}
	decrypted, err = decrypt(plaintext, passphrase, data[1], data[2])
	return
}

func encrypt(plaintext []byte, passphrase []byte) (encrypted string, salt string, iv string) {
	key, saltBytes := deriveKey(passphrase, nil)
	ivBytes := make([]byte, 12)
	// http://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf
	// Section 8.2
	rand.Read(ivBytes)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	encrypted = hex.EncodeToString(aesgcm.Seal(nil, ivBytes, plaintext, nil))
	salt = hex.EncodeToString(saltBytes)
	iv = hex.EncodeToString(ivBytes)
	return
}

func decrypt(data []byte, passphrase []byte, salt string, iv string) (plaintext []byte, err error) {
	saltBytes, _ := hex.DecodeString(salt)
	ivBytes, _ := hex.DecodeString(iv)
	key, _ := deriveKey(passphrase, saltBytes)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	plaintext, err = aesgcm.Open(nil, ivBytes, data, nil)
	return
}

func deriveKey(passphrase []byte, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 8)
		// http://www.ietf.org/rfc/rfc2898.txt
		// Salt.
		rand.Read(salt)
	}
	return pbkdf2.Key(passphrase, salt, 1000, 32, sha256.New), salt
}

// compressByte returns a compressed byte slice.
func compressByte(src []byte) []byte {
	compressedData := new(bytes.Buffer)
	compress(src, compressedData, 9)
	return compressedData.Bytes()
}

// decompressByte returns a decompressed byte slice.
func decompressByte(src []byte) []byte {
	compressedData := bytes.NewBuffer(src)
	deCompressedData := new(bytes.Buffer)
	decompress(compressedData, deCompressedData)
	return deCompressedData.Bytes()
}

// compress uses flate to compress a byte slice to a corresponding level
func compress(src []byte, dest io.Writer, level int) {
	compressor, _ := flate.NewWriter(dest, level)
	compressor.Write(src)
	compressor.Close()
}

// compress uses flate to decompress an io.Reader
func decompress(src io.Reader, dest io.Writer) {
	decompressor := flate.NewReader(src)
	io.Copy(dest, decompressor)
	decompressor.Close()
}

// src is seeds the random generator for generating random strings
var src = math_rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// RandStringBytesMaskImprSrc prints a random string
func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
