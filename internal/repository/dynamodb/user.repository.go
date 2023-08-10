package dynamodb

import (
	"context"
	"errors"

	"github.com/aqaurius6666/goserverless/internal/entity"
	"github.com/aqaurius6666/goserverless/internal/usecase"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
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

// https://github.com/awsdocs/aws-doc-sdk-examples/blob/main/gov2/dynamodb/actions/table_basics.go#L238
func (d *DdbUserRepository) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	output, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(string(d.TableName)),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: id},
			"sk": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if output.Item == nil {
		return nil, errors.New("not found")
	}
	var user entity.User
	if err = attributevalue.UnmarshalMap(output.Item, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *DdbUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}
	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(string(d.TableName)),
		Item:      av,
	})
	return err
}
