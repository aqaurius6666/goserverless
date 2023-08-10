package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type CtxKey string

var LambdaDepsCtxKey CtxKey = "deps"

type event struct {
	Name string `json:"name"`
}

func handler(ctx context.Context, e event) error {
	deps := ctx.Value(LambdaDepsCtxKey).(LambdaDeps)
	user, err := deps.UserUseCase.GetUserById(ctx, "test")
	if err != nil {
		return err
	}
	_ = user
	return nil
}

func main() {
	ctx := context.Background()
	deps, err := BuildDeps(LambdaOpts{
		DdbTableName: "dev-my-reading-table",
	})
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, LambdaDepsCtxKey, deps)
	lambda.StartWithOptions(handler, lambda.WithContext(ctx))
}
