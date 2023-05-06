package encrypt

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestDes(t *testing.T) {
	// create a new Des struct
	des := NewDes()

	// generate a secret key
	key := []byte("2fa6c1e9")
	// create a message to be encrypted
	message := "this is a message"

	// encrypt
	cipherText := des.DESEncrypt([]byte(message), key)

	// decrypt
	plainText := des.DESDecrypt(cipherText, key)

	// compare
	assert.Equal(t, string(plainText), message)

}

func Test3Des(t *testing.T) {
	// create a new Des struct
	des := NewDes()

	// generate a secret key
	key := []byte("2fa6c1e92fa6c1e92fa6c1e9")
	// create a message to be encrypted
	message := "this is a message"

	// encrypt
	cipherText := des.DES3Encrypt([]byte(message), key)

	// decrypt
	plainText := des.DES3Decrypt(cipherText, key)

	// compare
	assert.Equal(t, string(plainText), message)

}
