package main

import (
	"encoding/json"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// UserData - the details of the user filling out the form
type UserData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Validate that all values are present
func (ud UserData) Validate() (string, bool) {
	if len(ud.Name) < 3 {
		return "Invalid Name", false
	}
	if !emailRegex.MatchString(ud.Email) {
		return "Invalid Email", false
	}
	if len(ud.Message) < 10 {
		return "Invalid Message", false
	}
	return "Success", true
}

// Response - API gateway default res type
type Response events.APIGatewayProxyResponse

// Message - Response message to api call
type Message struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// Handler sends the response
func Handler(req events.APIGatewayProxyRequest) (Response, error) {
	userData := &UserData{}
	json.Unmarshal([]byte(req.Body), userData)
	message, ok := userData.Validate()

	if !ok {
		return Respond(400, message)
	}

	return Respond(200, message)
}

// Respond to the caller
func Respond(status int, message string) (Response, error) {
	resMessage := &Message{}

	resMessage.Status = status
	resMessage.Message = message

	json, err := json.Marshal(resMessage)
	if err != nil {
		return Response{Body: "Something went wrong", StatusCode: 500}, err
	}

	return Response{Body: string(json), StatusCode: 400}, nil
}

func main() {
	lambda.Start(Handler)
}
