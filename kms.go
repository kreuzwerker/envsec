package envsec

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type KMSMethod struct {
	client *kms.KMS
	keyId  string
}

func (c *KMSMethod) Decrypt(ciphertext string) (string, error) {

	decoded, err := base64.StdEncoding.DecodeString(ciphertext)

	if err != nil {
		return empty, err
	}

	resp, err := c.client.Decrypt(&kms.DecryptInput{
		CiphertextBlob: decoded,
	})

	if err != nil {
		return empty, err
	}

	return string(resp.Plaintext), nil

}

func (c *KMSMethod) Encrypt(plaintext string) (string, error) {

	resp, err := c.client.Encrypt(&kms.EncryptInput{
		Plaintext: []byte(plaintext),
		KeyId:     aws.String(c.keyId),
	})

	if err != nil {
		return empty, err
	}

	encoded := base64.StdEncoding.EncodeToString(resp.CiphertextBlob)

	return encoded, nil

}

func NewKMSMethod(region, keyId string) (*KMSMethod, error) {

	method := &KMSMethod{
		client: kms.New(session.New(), &aws.Config{Region: aws.String(region)}),
		keyId:  keyId,
	}

	return method, nil

}
