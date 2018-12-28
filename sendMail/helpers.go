package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

// Message - Response message to api call
type Message struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Respond to the caller
func Respond(status int, message string) (events.APIGatewayProxyResponse, error) {
	resMessage := &Message{}

	resMessage.Status = status
	resMessage.Message = message

	json, err := json.Marshal(resMessage)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Something went wrong",
			StatusCode: status,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(json),
		StatusCode: status,
	}, nil
}
