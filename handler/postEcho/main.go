package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type event struct {
	Name string `json:"name"`
}

func handler(ctx context.Context, e event) {

}

func main() {
	lambda.Start(handler)
}
