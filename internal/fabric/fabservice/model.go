package fabservice

// PublicKey address: publicKey
type PublicKey struct {
	PK string `json:"pk"`
}

// AddressFile address receive files
type AddressFile struct {
	FileID []string `json:"file_id,omitempty"`
	Files  []File   `json:"files,omitempty"`
}

// File file info
type File struct {
	FileID   string `json:"file_id"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize string `json:"file_size"`
	FileDate string `json:"file_date"`
	FileHash string `json:"file_hash"`
}
