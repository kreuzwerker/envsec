# Envsec

[![Build Status](https://travis-ci.org/kreuzwerker/envsec.svg?branch=master)](https://travis-ci.org/kreuzwerker/envsec)

Envsec (`es`) encrypts and decrypts environment variables using [AWS KMS](https://aws.amazon.com/kms/). When encrypting it passes the variable values to KMS, let's the service encrypt them and prefixes the variables with a configurable prefix (default: `ENVSEC_`). When decrypting, it `exec`utes a given process and passes the decrypted environment variables (without the prefix) to the new process.

The usage of KMS allows authorized operators to encrypt configuration secrets and submit them to version control, [ECS task definitions](http://docs.aws.amazon.com/AmazonECS/latest/developerguide/task_defintions.html) and other sources of configuration truths while the decryption operation can be bound to different principals, e.g. the role of an EC2 instance's instance profile.

## Example

This example uses `es` in combination with [envplate](https://github.com/kreuzwerker/envplate) to replace confidential variable references in the following config file:

```
secret1 = "${SECRET1}"
secret2 = "${SECRET2}"
```

The demo script `examples/test.sh` now

1. exports the secrets in plaintext to the environment
* calls `es` with the ARN of an AWS KMS key which exports the now encrypted variables as key-value-pairs, prepends an `export` to each key-value pair and writes the result to a temporary file called `test.envs`
* unsets the plaintext secrets from the environment
* exports the encrypted environment variables into the environment
* calls `es` which decrypts the environment variables, `exec`s envplate which replaces the variable references with the correct values and eventually `exec`s `cat test.config`

```
BINARY=../build/darwin_amd64/es-bin

export SECRET1=opensesame
export SECRET2=foreyesonly

$BINARY enc --arn=arn:aws:kms:eu-west-1:1234:key/e17e953c-d06e-4a4e-916b-54c681ca80d4 SECRET1 SECRET2 | sed  "s/^/export /g" > test.envs

unset SECRET1
unset SECRET2

eval $(cat test.envs)

$BINARY dec --region eu-west-1 -- /usr/local/bin/ep -- cat test.envs
```

The result is the following config file:

```
secret1 = "opensesame"
secret2 = "foreyesonly"
```

### Integration test and CloudFormation template

The integration tests in `kms_test.go` utilize a test environment described using a [tacks](https://github.com/kreuzwerker/tacks) template. At it's core (starting at line 14) the file `example/cloudformation.yml` contains:

1. two roles (`encrypt` and `decrypt`)
* a user that can assume those roles
* a keypair for that user
* a KMS key with a key policy that ties encryption and decryption rights to the two roles; in order to be able to delete the key together with the stack, full access is given to the account the key was created in

In a production setup the roles (or rather the underlying policies) would be tied to an instance profile of an EC2 instance (for decrypting) and potentially to the group of developers / operators relevant to this project (for encrypting).
