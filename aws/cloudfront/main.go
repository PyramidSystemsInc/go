package cloudfront

import (
  "time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/cloudfront"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/str"
)

func CreateDistributionFromS3Bucket(domainName string, awsSession *session.Session) string {
  cloudfrontClient := cloudfront.New(awsSession)
  OAIResult, err := cloudfrontClient.CreateCloudFrontOriginAccessIdentity(&cloudfront.CreateCloudFrontOriginAccessIdentityInput{
    CloudFrontOriginAccessIdentityConfig: &cloudfront.OriginAccessIdentityConfig{
      CallerReference: aws.String(time.Now().String()),
      Comment: aws.String(str.Concat("Identity for ", domainName)),
    },
  })
  errors.LogIfError(err)
  originAccessId := str.Concat("origin-access-identity/cloudfront/", *OAIResult.CloudFrontOriginAccessIdentity.Id)
  distroResult, err := cloudfrontClient.CreateDistribution(&cloudfront.CreateDistributionInput{
    DistributionConfig: &cloudfront.DistributionConfig{
      CallerReference: aws.String(time.Now().String()),
      Comment: aws.String(str.Concat("Distribution for ", domainName)),
      CustomErrorResponses: &cloudfront.CustomErrorResponses{
        Items: []*cloudfront.CustomErrorResponse{
          &cloudfront.CustomErrorResponse{
            ErrorCachingMinTTL: aws.Int64(86400),
            ErrorCode: aws.Int64(403),
            ResponseCode: aws.String("200"),
            ResponsePagePath: aws.String("/index.html"),
          },
          &cloudfront.CustomErrorResponse{
            ErrorCachingMinTTL: aws.Int64(86400),
            ErrorCode: aws.Int64(404),
            ResponseCode: aws.String("200"),
            ResponsePagePath: aws.String("/index.html"),
          },
        },
        Quantity: aws.Int64(2),
      },
      DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
        AllowedMethods: &cloudfront.AllowedMethods{
          Items: []*string{
            aws.String("GET"),
            aws.String("HEAD"),
            aws.String("OPTIONS"),
          },
          Quantity: aws.Int64(3),
        },
        DefaultTTL: aws.Int64(300),
        ForwardedValues: &cloudfront.ForwardedValues{
          Cookies: &cloudfront.CookiePreference{
            Forward: aws.String("none"),
          },
          QueryString: aws.Bool(false),
        },
        MinTTL: aws.Int64(0),
        TargetOriginId: aws.String(domainName),
        TrustedSigners: &cloudfront.TrustedSigners{
          Enabled: aws.Bool(false),
          Quantity: aws.Int64(0),
        },
        ViewerProtocolPolicy: aws.String("allow-all"),
      },
      DefaultRootObject: aws.String("index.html"),
      Enabled: aws.Bool(true),
      Origins: &cloudfront.Origins{
        Items: []*cloudfront.Origin{
          &cloudfront.Origin{
            DomainName: aws.String(str.Concat(domainName, ".s3.amazonaws.com")),
            Id: aws.String(domainName),
            S3OriginConfig: &cloudfront.S3OriginConfig{
              OriginAccessIdentity: aws.String(originAccessId),
            },
          },
        },
        Quantity: aws.Int64(1),
      },
    },
  })
  errors.QuitIfError(err)
  return *distroResult.Distribution.DomainName
}
