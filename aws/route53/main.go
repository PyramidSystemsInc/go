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

func DeleteHostedZone(domainName string, awsSession *session.Session) {
  route53Client := route53.New(awsSession)
  hostedZoneId, _ := findDomainNameId(domainName, route53Client)
  if hostedZoneId != "" {
    listResult, err := route53Client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
      HostedZoneId: aws.String(hostedZoneId),
    })
    errors.LogIfError(err)
    records := listResult.ResourceRecordSets
    var batchChanges []*route53.Change
    for _, record := range records {
      if *record.Type != "SOA" && *record.Type != "NS" {
        batchChanges = append(batchChanges, &route53.Change{
          Action: aws.String("DELETE"),
          ResourceRecordSet: record,
        })
      }
    }
    if len(batchChanges) > 0 {
      _, err = route53Client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
        ChangeBatch: &route53.ChangeBatch{
          Changes: batchChanges,
          Comment: aws.String("Deleted record(s) as part of call to PyramidSystemsInc/go/aws/route53/DeleteHostedZone"),
        },
        HostedZoneId: aws.String(hostedZoneId),
      })
      errors.LogIfError(err)
    }
    _, err = route53Client.DeleteHostedZone(&route53.DeleteHostedZoneInput{
      Id: aws.String(hostedZoneId),
    })
    errors.LogIfError(err)
  }
}

func DeleteRecord(domainName string, recordName string, awsSession *session.Session) {
  route53Client := route53.New(awsSession)
  hostedZoneId, _ := findDomainNameId(domainName, route53Client)
  if hostedZoneId != "" {
    listResult, err := route53Client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
      HostedZoneId: aws.String(hostedZoneId),
    })
    errors.LogIfError(err)
    records := listResult.ResourceRecordSets
    var batchChanges []*route53.Change
    for _, record := range records {
      if *record.Name == recordName {
        batchChanges = append(batchChanges, &route53.Change{
          Action: aws.String("DELETE"),
          ResourceRecordSet: record,
        })
      }
    }
    if len(batchChanges) > 0 {
      _, err = route53Client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
        ChangeBatch: &route53.ChangeBatch{
          Changes: batchChanges,
          Comment: aws.String("Deleted record(s) as part of call to PyramidSystemsInc/go/aws/route53/DeleteRecord"),
        },
        HostedZoneId: aws.String(hostedZoneId),
      })
      errors.LogIfError(err)
    }
  }
}

func TagHostedZone(domainName string, key string, value string, awsSession *session.Session) {
  route53Client := route53.New(awsSession)
  id, err := findDomainNameId(domainName, route53Client)
  errors.QuitIfError(err)
  _, err = route53Client.ChangeTagsForResource(&route53.ChangeTagsForResourceInput{
    AddTags: []*route53.Tag{
      &route53.Tag{
        Key: aws.String(key),
        Value: aws.String(value),
      },
    },
    ResourceId: aws.String(id),
    ResourceType: aws.String("hostedzone"),
  })
  errors.LogIfError(err)
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
