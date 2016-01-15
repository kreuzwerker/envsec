package envsec

import (
	"log"
	"strings"

	"github.com/yawn/envmap"
)

// Handler implements the parallel execution of a Method
// over a specific prefix
type Handler struct {
	Method Method
	Prefix string
}

func (h Handler) Decrypt(env envmap.Envmap) (result []string) {
	return h.apply(env, h.decrypt)
}

func (h Handler) Encrypt(env envmap.Envmap) (result []string) {
	return h.apply(env, h.encrypt)
}

func (h Handler) apply(env envmap.Envmap, f func(k, v string, results chan string)) (result []string) {

	var results = make(chan string, len(env))

	for k, v := range env {
		go f(k, v, results)
	}

	for i := 0; i < len(env); i++ {
		result = append(result, <-results)
	}

	return

}

func (h Handler) decrypt(k, v string, results chan string) {

	var err error

	if strings.HasPrefix(k, h.Prefix) {

		// TODO: use slice instead
		k = strings.Replace(k, h.Prefix, "", 1)
		v, err = h.Method.Decrypt(v)

		if err != nil {
			log.Fatalf("Failed to decrypt secure variable %q: %v", k, err)
		} else if Verbose {
			log.Printf("Decrypted secure variable %q", k)
		}

	}

	results <- envmap.Join(k, v)

}

func (h Handler) encrypt(k, v string, results chan string) {

	var err error

	if v == empty {
		log.Fatalf("Empty or non-existing environment variable %q", k)
	}

	v, err = h.Method.Encrypt(v)

	if err != nil {
		log.Fatalf("Failed to encrypt secure variable %q: %v", k, err)
	} else if Verbose {
		log.Printf("Encrypted secure variable %q", k)
	}

	results <- strings.Join([]string{h.Prefix, envmap.Join(k, v)}, "")

}
