package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	strip "github.com/grokify/html-strip-tags-go"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// UserData - the details of the user filling out the form
type UserData struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Message     string    `json:"message"`
	CreatedDate time.Time `json:"created_date"`
}

// Validate that all values are present
func (ud UserData) Validate() (string, bool) {
	if len(ud.Name) < 3 || len(ud.Message) > 100 {
		return "Invalid Name", false
	}
	if !emailRegex.MatchString(ud.Email) || len(ud.Message) > 300 {
		return "Invalid Email", false
	}
	if len(ud.Message) < 10 || len(ud.Message) > 300 {
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
	`, strip.StripTags(ud.Name), strip.StripTags(ud.Email), strip.StripTags(ud.Message))
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
	`, strip.StripTags(ud.Name), strip.StripTags(ud.Email), strip.StripTags(ud.Message))
}

// GetByEmail - Check the dynamoDB for the user by email
func (ud UserData) GetByEmail() (*UserData, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := dynamodb.New(sess)

	fromDB := &UserData{}

	fmt.Println("Looking Up By Email:", ud.Email+" "+ud.Name)
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(ud.Email),
			},
			"name": {
				S: aws.String(ud.Name),
			},
		},
	})

	if err != nil {
		fmt.Println("Error looking up email:")
		fmt.Println(err.Error())
		return fromDB, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, fromDB)
	if err != nil {
		fmt.Println(err.Error())
		return fromDB, err
	}

	return fromDB, nil
}

// AddToDB - Check the dynamoDB for the user by email
func (ud UserData) AddToDB() (*UserData, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	svc := dynamodb.New(sess)

	itemToSave := &ud
	itemToSave.CreatedDate = time.Now()

	u, err := ud.GetByEmail()
	if err != nil {
		return itemToSave, err
	}

	if time.Now().Sub(u.CreatedDate).Seconds() <= 604800 {
		return u, errors.New("You have already submitted a request in the last week")
	}

	av, err := dynamodbattribute.MarshalMap(itemToSave)
	if err != nil {
		fmt.Println("Error marshalling map:")
		fmt.Println(err.Error())
		return itemToSave, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}
	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
		return itemToSave, err
	}

	return itemToSave, nil
}
