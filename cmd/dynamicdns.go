package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func checkCurrentIp() (string, error) {
	resp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", fmt.Errorf("getting the current ip: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	currentIp := strings.TrimSpace(string(body))
	if net.ParseIP(strings.TrimSpace(string(body))) == nil {
		return "", fmt.Errorf("content blocked")
	}

	return currentIp, nil
}

func dynamicDNS() {
	currentIp, err := checkCurrentIp()
	if err != nil {
		fmt.Println("Error checking ip:", err)
	}

	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("Error setting up session:", err)
	}
	svc := route53.New(sess)

	id := "Z0573558CQKPQJAUG4FW"
	recordSet, err := svc.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId: &id,
	})
	if err != nil {
		fmt.Println("Error getting hosted zone id:", err)
	}

	recordIp := *recordSet.ResourceRecordSets[0].ResourceRecords[0].Value
	recordIpWWW := *recordSet.ResourceRecordSets[3].ResourceRecords[0].Value

	if currentIp != recordIp || recordIpWWW != recordIp {
		fmt.Printf("%s, %s\n", currentIp, recordIp)

		recordSet.ResourceRecordSets[0].ResourceRecords[0].Value = &currentIp
		recordSet.ResourceRecordSets[3].ResourceRecords[0].Value = &currentIp

		change1 := route53.Change{}
		change1.SetAction("UPSERT")
		change1.SetResourceRecordSet(recordSet.ResourceRecordSets[0])

		change2 := route53.Change{}
		change2.SetAction("UPSERT")
		change2.SetResourceRecordSet(recordSet.ResourceRecordSets[3])

		changeBatch := route53.ChangeBatch{Changes: []*route53.Change{&change1, &change2}}

		input := route53.ChangeResourceRecordSetsInput{
			ChangeBatch:  &changeBatch,
			HostedZoneId: &id,
		}

		recordSetOut, err := svc.ChangeResourceRecordSets(&input)
		if err != nil {
			fmt.Println("Error changing resource record set")
		}
		_ = recordSetOut
	}
}

func main() {
	for {
		dynamicDNS()
		time.Sleep(time.Minute * 5)
	}
}
