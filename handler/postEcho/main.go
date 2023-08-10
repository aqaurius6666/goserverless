package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aqaurius6666/goserverless/internal/repository/dynamodb"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type CtxKey string

var LambdaDepsCtxKey CtxKey = "deps"

type requestBody struct {
	Name string `json:"name"`
}

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var (
		body requestBody
		err  error
	)
	deps := ctx.Value(LambdaDepsCtxKey).(LambdaDeps)
	if err = json.Unmarshal([]byte(e.Body), &body); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}

	err = deps.UserUseCase.CreateUser(ctx, body.Name)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "ok",
	}, nil
}

func main() {
	ctx := context.Background()
	deps, err := BuildDeps(LambdaOpts{
		DdbTableName: dynamodb.DdbTable(os.Getenv("DYNAMODB_TABLE_NAME")),
	})
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, LambdaDepsCtxKey, deps)
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}
