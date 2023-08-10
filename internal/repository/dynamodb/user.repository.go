package dynamodb

import (
	"context"

	"github.com/aqaurius6666/goserverless/internal/entity"
	"github.com/aqaurius6666/goserverless/internal/usecase"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
)

var UserSet = wire.NewSet(
	wire.Struct(new(DdbUserRepository), "*"),
	wire.Bind(new(usecase.UserRepository), new(*DdbUserRepository)),
)

var _ usecase.UserRepository = (*DdbUserRepository)(nil)

type DdbUserRepository struct {
	TableName DdbTable
	Client    *dynamodb.Client
}

func (d *DdbUserRepository) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	return &entity.User{Name: "test"}, nil
}
