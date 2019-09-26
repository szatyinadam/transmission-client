package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

func Consume(sqsConfig *Sqs) {
	svc := createSqs(sqsConfig)
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(sqsConfig.QueueUrl),
		MaxNumberOfMessages: aws.Int64(3),
		VisibilityTimeout:   aws.Int64(30),
		WaitTimeSeconds:     aws.Int64(20),
	}
	for {
		receiveResp, err := svc.ReceiveMessage(receiveParams)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Received message - %s", receiveResp)
		if receiveResp != nil {
			deleteReceivedMessages(svc, sqsConfig, receiveResp)
		}
	}
}

func createSqs(sqsConfig *Sqs) *sqs.SQS {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(sqsConfig.Region),
		Credentials: credentials.NewSharedCredentials(sqsConfig.CredentialPath, sqsConfig.CredentialProfile),
		MaxRetries:  aws.Int(5),
	})
	if err != nil {
		log.Fatal(err)
	}
	svc := sqs.New(sess)
	return svc
}

func deleteReceivedMessages(svc *sqs.SQS, sqsConfig *Sqs, receivedMessages *sqs.ReceiveMessageOutput) {
	for _, message := range receivedMessages.Messages {
		deleteParams := &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(sqsConfig.QueueUrl),
			ReceiptHandle: message.ReceiptHandle,
		}
		_, err := svc.DeleteMessage(deleteParams)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Delete message - Message ID: %s has beed deleted.", *message.MessageId)
	}
}
