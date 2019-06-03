package aws

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/PyramidSystemsInc/go/commands"
  "github.com/PyramidSystemsInc/go/errors"
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

func GetAccessKey() (string, error) {
	key, err := commands.Run("aws configure get aws_access_key_id", "")
	return key, err
}

func GetSecretKey() (string, error) {
	key, err := commands.Run("aws configure get aws_secret_access_key", "")
	return key, err
}
