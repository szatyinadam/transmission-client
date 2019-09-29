package main

import (
	"./mock"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReceive(t *testing.T) {
	t.Run("Receive SQS message test", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSqsService, receiveParams, messageChannel := createParams(ctrl)

		Receive(mockSqsService, receiveParams, messageChannel)

		assert.Equal(t, 1, len(messageChannel))
		m := <-messageChannel
		assert.Equal(t, "test", *m.Body)
	})
}

func createParams(ctrl *gomock.Controller) (*mock_sqsiface.MockSQSAPI, *sqs.ReceiveMessageInput, chan *sqs.Message) {
	mockSqsService := mock_sqsiface.NewMockSQSAPI(ctrl)
	mockSqsService.
		EXPECT().
		ReceiveMessage(gomock.Any()).
		Return(createMessageOutput(), nil)
	queueUrl := "queue"
	receiveParams := &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(5),
		WaitTimeSeconds:     aws.Int64(20),
	}
	messageChannel := make(chan *sqs.Message, 5)
	return mockSqsService, receiveParams, messageChannel
}

func createMessage() *sqs.Message {
	body := "test"
	md5OfBody := "711a46d2bac61a9ec61e43b49bc4312a"
	messageId := "57b47cdd-3ee9-4101-8aef-2d87ab9a80a0"
	receiptHandle := "AQEB+L8B8EwJWJ0H4+dbJkVoQzqWF+CpmBpzIHhhC5EuLb4RHybpRirfb8ocqMii8qnlrpYnOam4b6OrKkcZmTw9UUrn0NI2tBUQsUITGJHaow9HBhImiWvghSEtwL1fYyJ3VcCvoaBB52CbsTZ9E+BwBB5FUc+3pX5F8KnYY5BQ2eMtn0PlhaurEYOU8/dbTSg3/IbriQ6C/+WG8eQCuMs379qVE38o3ZOYpj8jSGXJ7ZMnWDM224s6JFQ1BvdNNfdUKEOKLfbsQz8PIP4pthWc9A=="
	message := sqs.Message{
		Body:          &body,
		MD5OfBody:     &md5OfBody,
		MessageId:     &messageId,
		ReceiptHandle: &receiptHandle,
	}
	return &message
}

func createMessageOutput() *sqs.ReceiveMessageOutput {
	output := sqs.ReceiveMessageOutput{
		Messages: []*sqs.Message{createMessage()},
	}
	return &output
}

func TestDeleteReceivedMessages(t *testing.T) {
	t.Run("Delete received SQS message ", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockSqsService := mock_sqsiface.NewMockSQSAPI(ctrl)
		mockSqsService.
			EXPECT().
			DeleteMessage(gomock.Any()).
			Times(1).
			Return(nil, nil)
		config := Sqs{QueueUrl: "abc"}

		DeleteReceivedMessages(mockSqsService, &config, createMessage())
	})
}
