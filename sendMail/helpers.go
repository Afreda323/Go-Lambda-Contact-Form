package main

import (
	"encoding/json"
	"os"

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

	resp := events.APIGatewayProxyResponse{Headers: make(map[string]string)}
	resp.Headers["Access-Control-Allow-Origin"] = os.Getenv("ALLOWED_DOMAIN")
	resp.Headers["Access-Control-Allow-Credentials"] = "true"
	resp.StatusCode = status

	json, err := json.Marshal(resMessage)
	if err != nil {
		resp.Body = "Something went wrong"
		return resp, nil
	}

	resp.Body = string(json)
	return resp, nil
}
