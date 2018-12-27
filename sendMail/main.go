package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Lambda Response type
type Response events.APIGatewayProxyResponse

// Handler sends the response
func Handler() (Response, error) {
	return Response{Body: "Hello", StatusCode: 200}, nil
}

func main() {
	lambda.Start(Handler)
}
