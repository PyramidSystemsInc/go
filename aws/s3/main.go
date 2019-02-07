package s3

import (
  "strings"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/PyramidSystemsInc/go/aws/util"
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

func DeleteBucket(bucketNameOrArn string, awsSession *session.Session) {
  bucketName := getBucketName(bucketNameOrArn)
  s3Client := s3.New(awsSession)
  _, err := s3Client.DeleteBucket(&s3.DeleteBucketInput{
    Bucket: aws.String(bucketName),
  })
  errors.QuitIfError(err)
}

func EmptyBucket(bucketNameOrArn string, awsSession *session.Session) {
  bucketName := getBucketName(bucketNameOrArn)
  s3Client := s3.New(awsSession)
  bucketObjects, err := s3Client.ListObjects(&s3.ListObjectsInput{
    Bucket: aws.String(bucketName),
  })
  errors.QuitIfError(err)
  bucketContents := bucketObjects.Contents
  objectIdentifiers := make([]*s3.ObjectIdentifier, 0)
  for _, file := range bucketContents {
    objectIdentifiers = append(objectIdentifiers, &s3.ObjectIdentifier{
      Key: file.Key,
    })
  }
  _, err = s3Client.DeleteObjects(&s3.DeleteObjectsInput{
    Bucket: aws.String(bucketName),
    Delete: &s3.Delete{
      Objects: objectIdentifiers,
      Quiet: aws.Bool(true),
    },
  })
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

func TagBucket(bucketName string, key string, value string, awsSession *session.Session) {
  s3Client := s3.New(awsSession)
  _, err := s3Client.PutBucketTagging(&s3.PutBucketTaggingInput{
    Bucket: aws.String(bucketName),
    Tagging: &s3.Tagging{
      TagSet: []*s3.Tag{
        &s3.Tag{
          Key: aws.String(key),
          Value: aws.String(value),
        },
      },
    },
  })
  errors.QuitIfError(err)
}

func getBucketName(arnOrName string) string {
  if util.IsArn(arnOrName) {
    return arnOrName[strings.LastIndex(arnOrName, ":::") + 3:len(arnOrName)]
  } else {
    return arnOrName
  }
}
