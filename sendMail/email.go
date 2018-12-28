package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SendEmail - Connect to AWS and send the email to the box.
func SendEmail(ud *UserData, subject string, text string, html string) (string, bool) {
	CharSet := "UTF-8"

	// Create conection to ses server
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := ses.New(sess)

	// Create instance of email inputs
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

	// dispatch email
	result, err := svc.SendEmail(input)

	if err != nil {
		LogEmailError(err)
		return "Something went wrong", false
	}

	fmt.Println("Email Result:", result)

	return "Success", true
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
