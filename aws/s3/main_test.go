package s3

import (
	"fmt"
	"log"
	"testing"

	pacaws "github.com/PyramidSystemsInc/go/aws"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// TestEncryptBucket creates and s3 bucket, encrypts, and then tests the successful enabling
// of an S3 bucket by calling GetBucketEncryption. If an error is returned the bucket is not encrypted.
func TestEncryptBucket(t *testing.T) {
	//create bucket
	session := pacaws.CreateAwsSession("us-east-2")

	svc := s3.New(session)
	name := "saveferris"

	input := &s3.CreateBucketInput{
		ACL:    aws.String("private"),
		Bucket: aws.String(name),
	}

	_, err := svc.CreateBucket(input)

	if err != nil {
		log.Fatal("can't create bucket", err)
	}

	//encrypt bucket
	key := "fc8181c8-a0ca-4a95-9bd7-8673b179dee5"
	EncryptBucket(name, key)
	einput := &s3.GetBucketEncryptionInput{Bucket: aws.String(name)}

	//test bucket is encrypted

	result, err := svc.GetBucketEncryption(einput)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	dinput := &s3.DeleteBucketInput{Bucket: aws.String(name)}
	svc.DeleteBucket(dinput)
}
