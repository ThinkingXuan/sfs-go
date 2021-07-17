package encrypt

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestDes(t *testing.T) {
	// generate a secret key
	key := []byte("2fa6c1e9")
	// create a message to be encrypted
	message := "this is a message"

	cipherText, err := DESEncrypt([]byte(message), key)
	if err != nil {
		t.Fatal("Aes encrypt failure:", err)
	}

	plainText, err := DESDecrypt(cipherText, key)
	if err != nil {
		t.Fatal("Aes decrypt failure:", err)
	}

	assert.Equal(t, string(plainText), message)

}
