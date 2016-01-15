package envsec

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
)

const secret = "this is a very secret message"

var (
	client      *sts.STS
	decryptRole = os.Getenv("ES_DECRYPT_ROLE")
	encryptRole = os.Getenv("ES_ENCRYPT_ROLE")
)

func assumeRole(sessionName, arn string, f func()) {

	const (
		awsAccessKeyID     = "AWS_ACCESS_KEY_ID"
		awsSecretAccessKey = "AWS_SECRET_ACCESS_KEY"
		awsSessionToken    = "AWS_SESSION_TOKEN"
	)

	var (
		oldAwsAccessKeyID     = os.Getenv(awsAccessKeyID)
		oldAwsSecretAccessKey = os.Getenv(awsSecretAccessKey)
	)

	client = sts.New(session.New())

	req := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(900),
		RoleArn:         aws.String(arn),
		RoleSessionName: aws.String(sessionName),
	}

	res, err := client.AssumeRole(req)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Assuming role %q\n", arn)

	creds := res.Credentials

	// utilizing the default credentials detection chain here - in a deployment
	// situation this would be an instance profile
	os.Setenv(awsAccessKeyID, *creds.AccessKeyId)
	os.Setenv(awsSecretAccessKey, *creds.SecretAccessKey)
	os.Setenv(awsSessionToken, *creds.SessionToken)

	f()

	os.Setenv(awsAccessKeyID, oldAwsAccessKeyID)
	os.Setenv(awsSecretAccessKey, oldAwsSecretAccessKey)
	os.Unsetenv(awsSessionToken)

}

func kmsMethod() Method {

	method, err := NewKMSMethod("eu-west-1", os.Getenv("ES_KEY_ID"))

	if err != nil {
		panic(err)
	}

	return method

}

func TestKMSEncryptAndDecrypt(t *testing.T) {

	assert := assert.New(t)

	var (
		a1, b1, a2, b2 string
		err            error
	)

	assumeRole("encrypt", encryptRole, func() {

		method := kmsMethod()

		a1, err = method.Encrypt(secret)
		assert.NoError(err)
		assert.NotEqual(secret, a1)

		b1, err = method.Decrypt(a1)
		assert.Error(err)
		assert.Empty(b1)

	})

	assumeRole("decrypt", decryptRole, func() {

		method := kmsMethod()

		assert.NotEmpty(a1)

		a2, err = method.Encrypt(secret)
		assert.Error(err)
		assert.Empty(a2)

		b2, err = method.Decrypt(a1)
		assert.NoError(err)
		assert.Equal(secret, b2)

	})

}
