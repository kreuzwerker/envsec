package envsec

const empty = ""

// Verbose enables verbose logging
var Verbose bool

// Method describes a method for encrypting and decrypting secrets
type Method interface {
	Decrypt(ciphertext string) (string, error)
	Encrypt(plaintext string) (string, error)
}
