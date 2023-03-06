package cmd

import (
	"github.com/fentec-project/gofe/abe"
	"github.com/stretchr/testify/assert"
	"sfs-go/internal/file"
	"sfs-go/internal/tools"
	"testing"
)

func TestMAABE(t *testing.T) {
	// create new MAABE struct with Global Parameters
	maabe := abe.NewMAABE()

	// create three authorities, each with two attributes
	attribs1 := []string{"auth1:at1", "auth1:at2"}
	attribs2 := []string{"auth2:at1", "auth2:at2"}
	attribs3 := []string{"auth3:at1", "auth3:at2"}

	auth, err := maabe.NewMAABEAuth("auth", attribs1)

	PubKeys1String := file.ReadWithFile("D:\\workspace\\go_workspace\\src\\sfs-go\\auth1")
	auth1, err := auth.DecodeAuth(tools.StringToByte(PubKeys1String))
	//auth1, err := maabe.NewMAABEAuth("auth1", attribs1)
	if err != nil {
		t.Fatalf("Failed generation authority %s: %v\n", "auth1", err)
	}

	PubKeys2String := file.ReadWithFile("D:\\workspace\\go_workspace\\src\\sfs-go\\auth2")
	auth2, err := auth.DecodeAuth(tools.StringToByte(PubKeys2String))
	//auth2, err := maabe.NewMAABEAuth("auth2", attribs2)
	if err != nil {
		t.Fatalf("Failed generation authority %s: %v\n", "auth2", err)
	}
	PubKeys3String := file.ReadWithFile("D:\\workspace\\go_workspace\\src\\sfs-go\\auth3")
	auth3, err := auth.DecodeAuth(tools.StringToByte(PubKeys3String))
	//auth3, err := maabe.NewMAABEAuth("auth3", attribs3)
	if err != nil {
		t.Fatalf("Failed generation authority %s: %v\n", "auth3", err)
	}

	// create a msp struct out of the boolean formula
	msp, err := abe.BooleanToMSP("((auth1:at1 AND auth2:at1) OR (auth1:at2 AND auth2:at2)) OR (auth3:at1 AND auth3:at2)", false)
	if err != nil {
		t.Fatalf("Failed to generate the policy: %v\n", err)
	}

	// define the set of all public keys we use
	pks := []*abe.MAABEPubKey{auth1.PubKeys(), auth2.PubKeys(), auth3.PubKeys()}

	// choose a message
	msg := "Attack at dawn!"

	// encrypt the message with the decryption policy in msp
	ct, err := maabe.Encrypt(msg, msp, pks)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v\n", err)
	}

	// choose a single user's Global ID
	gid := "gid1"

	// authority 1 issues keys to user
	keys1, err := auth1.GenerateAttribKeys(gid, attribs1)
	if err != nil {
		t.Fatalf("Failed to generate attribute keys: %v\n", err)
	}
	key11, key12 := keys1[0], keys1[1]

	// authority 2 issues keys to user
	keys2, err := auth2.GenerateAttribKeys(gid, attribs2)
	if err != nil {
		t.Fatalf("Failed to generate attribute keys: %v\n", err)
	}
	key21, key22 := keys2[0], keys2[1]

	// authority 3 issues keys to user
	keys3, err := auth3.GenerateAttribKeys(gid, attribs3)
	if err != nil {
		t.Fatalf("Failed to generate attribute keys: %v\n", err)
	}
	key31, key32 := keys3[0], keys3[1]

	// user tries to decrypt with different key combos
	ks1 := []*abe.MAABEKey{key11, key21, key31} // ok
	ks2 := []*abe.MAABEKey{key12, key22, key32} // ok
	ks3 := []*abe.MAABEKey{key11, key22}        // not ok
	ks4 := []*abe.MAABEKey{key12, key21}        // not ok
	ks5 := []*abe.MAABEKey{key31, key32}        // ok

	// try to decrypt all messages
	msg1, err := maabe.Decrypt(ct, ks1)
	if err != nil {
		t.Fatalf("Error decrypting with keyset 1: %v\n", err)
	}
	assert.Equal(t, msg, msg1)

	msg2, err := maabe.Decrypt(ct, ks2)
	if err != nil {
		t.Fatalf("Error decrypting with keyset 2: %v\n", err)
	}
	assert.Equal(t, msg, msg2)

	_, err = maabe.Decrypt(ct, ks3)
	assert.Error(t, err)

	_, err = maabe.Decrypt(ct, ks4)
	assert.Error(t, err)

	msg5, err := maabe.Decrypt(ct, ks5)
	if err != nil {
		t.Fatalf("Error decrypting with keyset 5: %v\n", err)
	}
	assert.Equal(t, msg, msg5)

	// generate keys with a different GID
	gid2 := "gid2"
	// authority 1 issues keys to user
	foreignKeys, err := auth1.GenerateAttribKeys(gid2, []string{"auth1:at1"})
	if err != nil {
		t.Fatalf("Failed to generate attribute key for %s: %v\n", "auth1:at1", err)
	}
	foreignKey11 := foreignKeys[0]
	// join two users who have sufficient attributes together, but not on their
	// own
	ks6 := []*abe.MAABEKey{foreignKey11, key21}
	// try and decrypt
	_, err = maabe.Decrypt(ct, ks6)
	assert.Error(t, err)

}
