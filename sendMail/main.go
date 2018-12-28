package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler sends the response
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userData := &UserData{}
	json.Unmarshal([]byte(req.Body), userData)
	message, ok := userData.Validate()

	if !ok {
		return Respond(400, message)
	}

	u, err := userData.AddToDB()
	if err != nil {
		return Respond(400, err.Error())
	}

	message, ok = SendEmail(
		u,
		"Form Submission by "+u.Name,
		u.GenStringEmail(),
		u.GenHTMLEmail(),
	)

	if !ok {
		return Respond(500, message)
	}

	return Respond(200, message)
}

func main() {
	lambda.Start(Handler)
}
