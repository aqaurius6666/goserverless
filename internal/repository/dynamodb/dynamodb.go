package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/wire"
)

type DdbTable string

var DdbSet = wire.NewSet(
	NewDynamodbClient,
)

func NewDynamodbClient() (*dynamodb.Client, error) {
	return dynamodb.New(dynamodb.Options{}), nil
}
