package server

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryption(t *testing.T) {
	word := []byte("hello, world")
	passphrase := []byte("passphrase")
	encrypted, err := encryptBytes(compressByte(word), passphrase)
	assert.Nil(t, err)
	fmt.Printf("encrypted: %s", encrypted)
	decrypted, err := decryptBytes(encrypted, passphrase)
	assert.Nil(t, err)
	assert.Equal(t, word, decompressByte(decrypted))
}
