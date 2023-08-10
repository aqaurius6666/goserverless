//go:build wireinject
// +build wireinject

package main

import (
	"github.com/aqaurius6666/goserverless/internal/repository/dynamodb"
	"github.com/aqaurius6666/goserverless/internal/usecase"
	"github.com/google/wire"
)

type LambdaDeps struct {
	UserUseCase usecase.UserUseCase
}

type LambdaOpts struct {
	DdbTableName dynamodb.DdbTable
}

func BuildDeps(opts LambdaOpts) (LambdaDeps, error) {
	wire.Build(
		usecase.UserSet,
		dynamodb.UserSet,
		dynamodb.DdbSet,
		wire.Struct(new(LambdaDeps), "*"),
		wire.FieldsOf(new(LambdaOpts), "DdbTableName"),
	)
	return LambdaDeps{}, nil
}
