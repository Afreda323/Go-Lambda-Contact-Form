package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
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

// GenStringEmail - Generate the email template for a user (String)
func (ud UserData) GenStringEmail() string {
	_, ok := ud.Validate()
	if !ok {
		return "Invalid User"
	}

	return fmt.Sprintf(`
		Your form was submitted.
		Name: %s
		Email: %s
		Message: %s
	`, ud.Name, ud.Email, ud.Message)
}

// GenHTMLEmail - Generate the email template for a user (HTML)
func (ud UserData) GenHTMLEmail() string {
	_, ok := ud.Validate()
	if !ok {
		return "Invalid User"
	}

	return fmt.Sprintf(`
		<h1>Your form was submitted.</h1>
		<hr />
		<p>
			<b>Name:</b> %s <br />
			<b>Email:</b> %s <br />
			<b>Message:</b> %s <br />
		</p>
	`, ud.Name, ud.Email, ud.Message)
}

// Response - API gateway default res type
type Response events.APIGatewayProxyResponse

// Message - Response message to api call
type Message struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// LogEmailError - If error occurs when sending email, log it to the console
// Codes taken directly from the docs
func LogEmailError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case ses.ErrCodeMessageRejected:
			fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
		case ses.ErrCodeMailFromDomainNotVerifiedException:
			fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
		case ses.ErrCodeConfigurationSetDoesNotExistException:
			fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}

// SendEmail - Connect to AWS and send the email to the box.
func SendEmail(ud *UserData, subject string, text string, html string) (string, bool) {
	CharSet := "UTF-8"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := ses.New(sess)

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(os.Getenv("DESIRED_RECIPIENT")),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(text),
				},
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(html),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(os.Getenv("DESIRED_RECIPIENT")),
	}

	result, err := svc.SendEmail(input)

	if err != nil {
		LogEmailError(err)
		return "Something went wrong", false
	}

	fmt.Println("Email Result:", result)

	return "Success", true
}

// Handler sends the response
func Handler(req events.APIGatewayProxyRequest) (Response, error) {
	userData := &UserData{}
	json.Unmarshal([]byte(req.Body), userData)
	message, ok := userData.Validate()

	if !ok {
		return Respond(400, message)
	}

	message, ok = SendEmail(
		userData,
		"Form Submission by "+userData.Name,
		userData.GenStringEmail(),
		userData.GenHTMLEmail(),
	)

	if !ok {
		return Respond(500, message)
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
