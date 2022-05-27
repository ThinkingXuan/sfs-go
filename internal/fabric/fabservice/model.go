package fabservice

// PublicKey address: publicKey
type PublicKey struct {
	PK string `json:"pk"`
}

// AddressFile address receive files
type AddressFile struct {
	FileEncrypt []EncryptEntity `json:"file_encrypt,omitempty"`
	Files       []File          `json:"files,omitempty"`
}

type EncryptEntity struct {
	FileID            string `json:"file_id"`
	FileEncryptCipher string `json:"file_encrypt_cipher"`

	FileRekey  string `json:"file_rekey"`
	NewCapsule string `json:"new_capsule"`
}

// File file info
type File struct {
	FileID   string `json:"file_id,omitempty"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize string `json:"file_size"`
	FileDate string `json:"file_date"`
	FileHash string `json:"file_hash"`
}
