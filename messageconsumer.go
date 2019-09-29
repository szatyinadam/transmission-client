package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"log"
)

func PollSqs(chanel chan<- *sqs.Message, sqsService sqsiface.SQSAPI, sqsConfig *Sqs) {
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            &sqsConfig.QueueUrl,
		MaxNumberOfMessages: aws.Int64(sqsConfig.MaxMessages),
		WaitTimeSeconds:     aws.Int64(20),
	}
	for {
		Receive(sqsService, receiveParams, chanel)
	}
}

func Receive(sqsService sqsiface.SQSAPI, receiveParams *sqs.ReceiveMessageInput, chanel chan<- *sqs.Message) {
	output, err := sqsService.ReceiveMessage(receiveParams)
	if err != nil {
		log.Fatalf("Failed to fetch SQS message %v", err)
	}
	for _, message := range output.Messages {
		chanel <- message
	}
}

func DeleteReceivedMessages(sqsService sqsiface.SQSAPI, sqsConfig *Sqs, message *sqs.Message) {
	deleteParams := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(sqsConfig.QueueUrl),
		ReceiptHandle: message.ReceiptHandle,
	}
	_, err := sqsService.DeleteMessage(deleteParams)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Delete message - Message ID: %s has beed deleted.", *message.MessageId)
}
