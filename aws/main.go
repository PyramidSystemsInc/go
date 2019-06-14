package aws

import (
	"github.com/PyramidSystemsInc/go/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func CreateAwsSession(region string) *session.Session {
  awsSession, err := session.NewSession(&aws.Config{
    Region: aws.String(region),
  })
  errors.QuitIfError(err)
  _, err = awsSession.Config.Credentials.Get()
  errors.QuitIfError(err)
  return awsSession
}

func GetAccessKey() (string) {
	return getSharedCredentials().AccessKeyID
}

func GetSecretKey() (string) {
	return getSharedCredentials().SecretAccessKey
}

func getSharedCredentials() credentials.Value {
	sharedCreds, err := credentials.NewSharedCredentials("", "").Get()
	errors.QuitIfError(err)
	return sharedCreds
}
