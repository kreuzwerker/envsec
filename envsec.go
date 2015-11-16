package envsec

const empty = ""

var Verbose bool

type Method interface {
	Decrypt(ciphertext string) (string, error)
	Encrypt(plaintext string) (string, error)
}
