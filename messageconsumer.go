package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

func Consume(sqsConfig *Sqs) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(sqsConfig.Region),
		Credentials: credentials.NewSharedCredentials(sqsConfig.CredentialPath, sqsConfig.CredentialProfile),
		MaxRetries:  aws.Int(5),
	})
	if err != nil {
		log.Fatal(err)
	}

	svc := sqs.New(sess)

	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sqsConfig.QueueUrl),
		MaxNumberOfMessages: aws.Int64(3),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(20),
	}
	receiveResp, err := svc.ReceiveMessage(receiveParams)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Received message - %s", receiveResp)
}
