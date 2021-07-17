package encrypt

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestAes(t *testing.T) {
	// generate a secret key
	key := []byte{0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
		0xBA, 0x37, 0x2F, 0x02, 0xC3, 0x92, 0x1F, 0x7D,
		0x7A, 0x3D, 0x5F, 0x06, 0x41, 0x9B, 0x3F, 0x2D,
	}
	// create a message to be encrypted
	message := "this is a message"

	cipherText, err := AESEncrypt([]byte(message), key)
	if err != nil {
		t.Fatal("Aes encrypt failure")
	}

	plainText, err := AESDecrypt(cipherText, key)
	if err != nil {
		t.Fatal("Aes decrypt failure")
	}

	assert.Equal(t, string(plainText), message)

}
