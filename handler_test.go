package envsec

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yawn/envmap"
)

// the upcase method encrypts / decrypts by upper / lower casing
type upcase struct {
}

func (u upcase) Decrypt(ciphertext string) (string, error) {
	return strings.ToLower(ciphertext), nil
}

func (u upcase) Encrypt(plaintext string) (string, error) {
	return strings.ToUpper(plaintext), nil
}

func TestHandler(t *testing.T) {

	assert := assert.New(t)

	h := Handler{
		Method: upcase{},
		Prefix: "FOO_",
	}

	enc := h.Encrypt(map[string]string{
		"one": "a",
		"two": "b",
	})

	res := envmap.ToMap(enc)

	assert.Equal(res["FOO_one"], "A")
	assert.Equal(res["FOO_two"], "B")

	dec := h.Decrypt(envmap.ToMap(res.ToEnv()))

	res = envmap.ToMap(dec)

	assert.Equal(res["one"], "a")
	assert.Equal(res["two"], "b")

}
