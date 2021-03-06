package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"log"
	"net/http"
)

func main() {
	config := ReadConfig()
	sqsService := createSqsService(&config.Sqs)
	messageChannel := make(chan *sqs.Message, config.Sqs.MaxMessages)

	go PollSqs(messageChannel, sqsService, &config.Sqs)
	processMessage(messageChannel, config, sqsService)

	log.Printf("Listening on queue: %s", config.Sqs.QueueUrl)
}

func processMessage(messageChannel chan *sqs.Message, config *Config, sqsService sqsiface.SQSAPI) {
	transmissionClient := TransmissionClient{
		config:     &config.Transmission,
		httpClient: &http.Client{},
	}
	for message := range messageChannel {
		log.Println(message)
		torrents, err := transmissionClient.GetTorrents()
		if err != nil {
			log.Fatal(err)
		}
		for index, name := range torrents {
			log.Printf("%d %s", index, name)
		}
		DeleteReceivedMessages(sqsService, &config.Sqs, message)
	}
}

func createSqsService(sqsConfig *SqsConfig) sqsiface.SQSAPI {
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
