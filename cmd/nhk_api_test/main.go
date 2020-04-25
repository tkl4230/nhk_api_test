package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	nhk "github.com/tkl4230/nhk_api_test"
)

func webhook() {
	nhk.Webhook()
}

func main() {
	lambda.Start(webhook)
}
