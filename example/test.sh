BINARY=../build/darwin_amd64/es-bin

export SECRET1=opensesame
export SECRET2=foreyesonly

$BINARY enc --arn=arn:aws:kms:eu-west-1:1234:key/e17e953c-d06e-4a4e-916b-54c681ca80d4 SECRET1 SECRET2 | sed  "s/^/export /g" > test.envs

unset SECRET1
unset SECRET2

eval $(cat test.envs)

$BINARY dec -- /usr/local/bin/ep -v test.config -- /bin/cat test.config
