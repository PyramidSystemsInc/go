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
      CallerReference: createCallerReference(),
      Comment: aws.String(str.Concat("Identity for ", domainName)),
    },
  })
  errors.LogIfError(err)
  originAccessId := str.Concat("origin-access-identity/cloudfront/", *OAIResult.CloudFrontOriginAccessIdentity.Id)
  distroResult, err := cloudfrontClient.CreateDistribution(&cloudfront.CreateDistributionInput{
    DistributionConfig: &cloudfront.DistributionConfig{
      Aliases: &cloudfront.Aliases{
        Items: []*string {
          aws.String(domainName),
        },
        Quantity: aws.Int64(1),
      },
      CallerReference: createCallerReference(),
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

func DisableDistribution(alias string, awsSession *session.Session) {
  cloudfrontClient := cloudfront.New(awsSession)
  distributionId := getDistributionIdUsingAlias(alias, cloudfrontClient)
  if distributionId != "" {
    result, err := cloudfrontClient.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
      Id: aws.String(distributionId),
    })
    errors.LogIfError(err)
    distributionConfig := result.DistributionConfig
    distributionConfig.SetEnabled(false)
    eTag := getDistributionETag(distributionId, cloudfrontClient)
    _, err = cloudfrontClient.UpdateDistribution(&cloudfront.UpdateDistributionInput {
      DistributionConfig: distributionConfig,
      Id: aws.String(distributionId),
      IfMatch: aws.String(eTag),
    })
    errors.LogIfError(err)
  }
}

func TagDistribution(distributionFqdn string, key string, value string, awsSession *session.Session) {
  cloudfrontClient := cloudfront.New(awsSession)
  arn, err := getArn(distributionFqdn, cloudfrontClient)
  errors.QuitIfError(err)
  _, err = cloudfrontClient.TagResource(&cloudfront.TagResourceInput{
    Resource: aws.String(arn),
    Tags: &cloudfront.Tags{
      Items: []*cloudfront.Tag{
        &cloudfront.Tag{
          Key: aws.String(key),
          Value: aws.String(value),
        },
      },
    },
  })
  errors.QuitIfError(err)
}

func getArn(distributionFqdn string, cloudfrontClient *cloudfront.CloudFront) (string, error) {
  distributions, err := cloudfrontClient.ListDistributions(&cloudfront.ListDistributionsInput{
    MaxItems: aws.Int64(500),
  })
  errors.QuitIfError(err)
  distributionSummaries := distributions.DistributionList.Items
  for _, distribution := range distributionSummaries {
    if *distribution.DomainName == distributionFqdn {
      return *distribution.ARN, nil
    }
  }
  return "", errors.New(str.Concat("Distribution not found with the provided domain name: ", distributionFqdn))
}

func getDistributionIdUsingAlias(targetAlias string, cloudfrontClient *cloudfront.CloudFront) string {
  distributions, err := cloudfrontClient.ListDistributions(&cloudfront.ListDistributionsInput{
    MaxItems: aws.Int64(500),
  })
  errors.QuitIfError(err)
  distributionSummaries := distributions.DistributionList.Items
  for _, distributionSummary := range distributionSummaries {
    if distributionSummary.Aliases != nil {
      for _, alias := range distributionSummary.Aliases.Items {
        if *alias == targetAlias {
          return *distributionSummary.Id
        }
      }
    }
  }
  return ""
}

func createCallerReference() *string {
  return aws.String(time.Now().String())
}

func getDistributionETag(id string, cloudfrontClient *cloudfront.CloudFront) string {
  distribution, err := cloudfrontClient.GetDistribution(&cloudfront.GetDistributionInput{
    Id: aws.String(id),
  })
  errors.QuitIfError(err)
  return *distribution.ETag
}
/*
func getOriginAccessIdentityETagByAlias(targetAlias string, cloudfrontClient *cloudfront.CloudFront) (string, error) {
  originAccessIdentities, err := cloudfrontClient.ListCloudFrontOriginAccessIdentities(&cloudfront.ListCloudFrontOriginAccessIdentitiesInput{
    MaxItems: aws.Int64(500),
  })
  var id string
  for _, originAccessIdentity := range originAccessIdentities.CloudFrontOriginAccessIdentityList.Items {
    if strings.HasSuffix(*originAccessIdentity.Comment, targetAlias) {
      id = *originAccessIdentity.Id
    }
  }
  errors.QuitIfError(err)
  if id != "" {
    originAccessIdentity, err := cloudfrontClient.GetCloudFrontOriginAccessIdentity(&cloudfront.GetCloudFrontOriginAccessIdentityInput {
      Id: aws.String(id),
    })
    errors.QuitIfError(err)
    return *originAccessIdentity.ETag, nil
  }
  return "", errors.New(str.Concat("Origin access identity not found using alias: ", targetAlias))
}
*/
