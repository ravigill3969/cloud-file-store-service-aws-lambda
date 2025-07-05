package main

import (
	"go-lambda/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.ImageHandler)
}
