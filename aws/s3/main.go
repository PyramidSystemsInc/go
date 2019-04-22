package s3

import (
  "fmt"
  "os"
  "strings"

  "github.com/PyramidSystemsInc/go/aws/util"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/logger"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/s3"
)

// MakeBucket The allowed values for the `access` parameter can be found here: https://docs.aws.amazon.com/AmazonS3/latest/dev/acl-overview.html#canned-acl
func MakeBucket(bucketName string, access string, region string, awsSession *session.Session) {
  s3Client := s3.New(awsSession)
  _, err := s3Client.CreateBucket(&s3.CreateBucketInput{
    ACL:    aws.String(access),
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
	DeleteAllObjectVersions(bucketName, awsSession)
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
      Quiet:   aws.Bool(true),
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
          Key:   aws.String(key),
          Value: aws.String(value),
        },
      },
    },
  })
  errors.QuitIfError(err)
}

func getBucketName(arnOrName string) string {
  if util.IsArn(arnOrName) {
    return arnOrName[strings.LastIndex(arnOrName, ":::")+3 : len(arnOrName)]
  } else {
    return arnOrName
  }
}

// EncryptBucket turns on encryption on the S3 bucket
func EncryptBucket(bucket, key string) {
  sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  svc := s3.New(sess)

  defEnc := &s3.ServerSideEncryptionByDefault{KMSMasterKeyID: aws.String(key), SSEAlgorithm: aws.String(s3.ServerSideEncryptionAwsKms)}
  rule := &s3.ServerSideEncryptionRule{ApplyServerSideEncryptionByDefault: defEnc}
  rules := []*s3.ServerSideEncryptionRule{rule}
  serverConfig := &s3.ServerSideEncryptionConfiguration{Rules: rules}
  input := &s3.PutBucketEncryptionInput{Bucket: aws.String(bucket), ServerSideEncryptionConfiguration: serverConfig}
  _, err := svc.PutBucketEncryption(input)
  if err != nil {
    fmt.Println("Got an error adding default KMS encryption to bucket", bucket)
    fmt.Println(err.Error())
    os.Exit(1)
  }

  logger.Info("Bucket " + bucket + " now has KMS encryption by default")
}

// EnableVersioning turns on version on the S3 bucket
func EnableVersioning(bucket string) {
  sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  svc := s3.New(sess)
  input := &s3.PutBucketVersioningInput{
    Bucket: aws.String(bucket),
    VersioningConfiguration: &s3.VersioningConfiguration{
      MFADelete: aws.String("Disabled"),
      Status:    aws.String("Enabled"),
    },
  }

  _, err := svc.PutBucketVersioning(input)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
    }
    return
  }
}

func DisableVersioning(bucket string) {
  sess := session.Must(session.NewSessionWithOptions(session.Options{
    SharedConfigState: session.SharedConfigEnable,
  }))

  svc := s3.New(sess)
  input := &s3.PutBucketVersioningInput{
    Bucket: aws.String(bucket),
    VersioningConfiguration: &s3.VersioningConfiguration{
      MFADelete: aws.String("Disabled"),
      Status:    aws.String("Suspended"),
    },
  }

  _, err := svc.PutBucketVersioning(input)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      default:
        fmt.Println(aerr.Error())
      }
    } else {
      // Print the error, cast err to awserr.Error to get the Code and
      // Message from an error.
      fmt.Println(err.Error())
    }
    return
  }
}

// DisableVersioning turns on versioning on the S3 bucket
// In the AWS console the bucket will be mark as 'disabled',
// in the AWS documentation the status is referred to as 'suspended'
func DisableVersioning(bucket string) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := s3.New(sess)
	input := &s3.PutBucketVersioningInput{
		Bucket: aws.String(bucket),
		VersioningConfiguration: &s3.VersioningConfiguration{
			MFADelete: aws.String("Disabled"),
			Status:    aws.String("Suspended"),
		},
	}

	_, err := svc.PutBucketVersioning(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
}

// DeleteAllObjectVersions gets an array of all the bucket object versions, iterates over them, and deletes them.
func DeleteAllObjectVersions(bucket string, awsSession *session.Session) {
	objectVersions := GetObjectVersions(bucket, awsSession)
	length := len(objectVersions.Versions)

	for i := 0; i < length; i++ {
		id := *objectVersions.Versions[i].VersionId
		key := *objectVersions.Versions[i].Key
		DeleteObjectVersion(id, bucket, key, awsSession)
	}
}

// DeleteObjectVersion deletes the specific version of an S3 bucket object.
func DeleteObjectVersion(id, bucket, key string, awsSession *session.Session) {
	svc := s3.New(awsSession)

	input := &s3.DeleteObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key), // S3 key (not encryption key)
		VersionId: aws.String(id),
	}

	_, err := svc.DeleteObject(input)

	if err != nil {
		fmt.Println(err.Error())
	}
}

// GetObjectVersions retuns the list of version for an S3 bucket.
func GetObjectVersions(bucket string, awsSession *session.Session) (result *s3.ListObjectVersionsOutput) {
	svc := s3.New(awsSession)
	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.ListObjectVersions(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	return result
}

// DeleteAllDeleteMarkers retrives the delete markers and deletes them.
func DeleteAllDeleteMarkers(bucket string, awsSession *session.Session) {
	deleteMarkers := GetObjectVersions(bucket, awsSession)
	length := len(deleteMarkers.DeleteMarkers)

	for i := 0; i < length; i++ {
		id := *deleteMarkers.DeleteMarkers[i].VersionId
		key := *deleteMarkers.DeleteMarkers[i].Key
		DeleteObjectVersion(id, bucket, key, awsSession)
	}
}
