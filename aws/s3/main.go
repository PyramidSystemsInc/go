package s3

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/PyramidSystemsInc/go/errors"
)

// The allowed values for the `access` parameter can be found here: https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#canned-acl
func MakeBucket(bucketName string, access string, region string, awsSession *session.Session) {
  s3Client := s3.New(awsSession)
  _, err := s3Client.CreateBucket(&s3.CreateBucketInput{
    ACL: aws.String(access),
    Bucket: aws.String(bucketName),
    CreateBucketConfiguration: &s3.CreateBucketConfiguration{
      LocationConstraint: aws.String(region),
    },
    ObjectLockEnabledForBucket: aws.Bool(false),
  })
  errors.QuitIfError(err)
}

func EnableWebsiteHosting(bucketName string, awsSession *session.Session) {
  s3Client := s3.New(awsSession)
  documentName := "index.html"
  _, err := s3Client.PutBucketWebsite(&s3.PutBucketWebsiteInput{
    Bucket: aws.String(bucketName),
    WebsiteConfiguration: &s3.WebsiteConfiguration{
      ErrorDocument: &s3.ErrorDocument{
        Key: aws.String(documentName),
      },
      IndexDocument: &s3.IndexDocument{
        Suffix: aws.String(documentName),
      },
    },
  })
  errors.QuitIfError(err)
}
