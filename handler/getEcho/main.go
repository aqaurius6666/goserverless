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

func handler(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	deps := ctx.Value(LambdaDepsCtxKey).(LambdaDeps)
	id, ok := e.QueryStringParameters["id"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "missing id",
		}, nil
	}
	user, err := deps.UserUseCase.GetUserById(ctx, id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}
	res, err := json.Marshal(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(res),
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
