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
	cipherText, err := des.DESEncrypt([]byte(message), key)
	if err != nil {
		t.Fatal("Aes encrypt failure:", err)
	}

	// decrypt
	plainText, err := des.DESDecrypt(cipherText, key)
	if err != nil {
		t.Fatal("Aes decrypt failure:", err)
	}

	// compare
	assert.Equal(t, string(plainText), message)

}
