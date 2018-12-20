package aws

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
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
