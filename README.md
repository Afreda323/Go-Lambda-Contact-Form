# Serverless Email Service (Work in progress)

A simple lambda that takes in a name, message and email address.  Sends an email to your desired mailbox.

Unique emails are saved to DynamoDB, emails can only be submitted once per week per address.

### Usage

Install Serverless, [follow the installation instructions here](https://serverless.com/blog/anatomy-of-a-serverless-app/#setup).  Requires at least version v1.26 or later as that's the version that comes with Golang support.  You'll also need to [install Go](https://golang.org/doc/install) and it's dependency manager, [dep](https://github.com/golang/dep).

When both of those tasks are done, cd into your `GOPATH` (more than likely ~/go/src/) and clone this project into that folder.  Then cd into the resulting folder and compile the source with `make`:

To deploy run

```serverless deploy```

There is only one enpoint ```/sendMail```.  

It requires a body consisting of:
```json
{
    "name": "String",
    "email": "String",
    "message": "String"
}
```

### Todo

- Add CI/CD
- Write tests