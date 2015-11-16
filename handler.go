package envsec

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Handler struct {
	Method Method
	Prefix string
}

func (h Handler) Decrypt(env []string) (result []string) {

	var results = make(chan string, len(env))

	for _, e := range env {
		go h.decrypt(e, results)
	}

	for i := 0; i < len(env); i++ {
		result = append(result, <-results)
	}

	return

}

func (h Handler) Encrypt(keys []string) (result []string) {

	var results = make(chan string, len(keys))

	for _, e := range keys {
		go h.encrypt(e, results)
	}

	for i := 0; i < len(keys); i++ {
		result = append(result, <-results)
	}

	return

}

func (h Handler) decrypt(env string, results chan string) {

	var err error

	s := strings.Split(env, "=")

	key := s[0]
	val := strings.Join(s[1:], "=")

	if strings.HasPrefix(key, h.Prefix) {

		key = strings.Replace(key, h.Prefix, "", 1)
		val, err = h.Method.Decrypt(val)

		if err != nil {
			log.Fatalf("Failed to decrypt secure variable %q: %v", key, err)
		} else if Verbose {
			log.Printf("Decrypted secure variable %q", key)
		}

	}

	results <- fmt.Sprintf("%s=%s", key, val)

}

func (h Handler) encrypt(key string, results chan string) {

	var (
		err error
		val = os.Getenv(key)
	)

	if val == empty {
		log.Fatalf("Empty or non-existing environemnt variable %q", key)
	}

	val, err = h.Method.Encrypt(val)

	if err != nil {
		log.Fatalf("Failed to encrypt secure variable %q: %v", key, err)
	} else if Verbose {
		log.Printf("Encrypted secure variable %q", key)
	}

	results <- fmt.Sprintf("%s%s=%s", h.Prefix, key, val)

}
