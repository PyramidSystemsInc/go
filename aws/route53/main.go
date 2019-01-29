package route53

import (
  "time"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/route53"
  "github.com/PyramidSystemsInc/go/errors"
  "github.com/PyramidSystemsInc/go/str"
)

func CreateHostedZone(domainName string, awsSession *session.Session) []string {
  route53Client := route53.New(awsSession)
  result, err := route53Client.CreateHostedZone(&route53.CreateHostedZoneInput{
    CallerReference: aws.String(time.Now().String()),
    Name: aws.String(domainName),
  })
  errors.LogIfError(err)
  nameServers := make([]string, 0)
  for _, nameServer := range result.DelegationSet.NameServers {
    nameServers = append(nameServers, *nameServer)
  }
  return nameServers
}

func ChangeRecord(domainName string, recordType string, recordName string, records []string, ttl int64, awsSession *session.Session) {
  route53Client := route53.New(awsSession)
  hostedZoneId, err := findDomainNameId(domainName, route53Client)
  errors.LogIfError(err)
  resourceRecords := make([]*route53.ResourceRecord, 0)
  for _, record := range records {
    resourceRecords = append(resourceRecords, &route53.ResourceRecord{
      Value: aws.String(record),
    })
  }
  _, err = route53Client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
    ChangeBatch: &route53.ChangeBatch{
      Changes: []*route53.Change{
        {
          Action: aws.String("UPSERT"),
          ResourceRecordSet: &route53.ResourceRecordSet{
            Name: aws.String(recordName),
            ResourceRecords: resourceRecords,
            TTL: aws.Int64(ttl),
            Type: aws.String(recordType),
          },
        },
      },
    },
    HostedZoneId: aws.String(hostedZoneId),
  })
  errors.QuitIfError(err)
}

func domainNamesMatch(domainNameA string, domainNameB string) bool {
  return domainNameA == domainNameB || domainNameA == str.Concat(domainNameB, ".") || str.Concat(domainNameA, ".") == domainNameB
}

func findDomainNameId(domainName string, route53Client *route53.Route53) (string, error) {
  result, err := route53Client.ListHostedZonesByName(&route53.ListHostedZonesByNameInput{
    DNSName: aws.String(domainName),
    MaxItems: aws.String("1"),
  })
  errors.QuitIfError(err)
  if domainNamesMatch(*result.HostedZones[0].Name, domainName) {
    return *result.HostedZones[0].Id, nil
  } else {
    return "", errors.New("Domain name provided could not be found")
  }
}
