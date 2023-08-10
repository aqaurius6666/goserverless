package entity

type User struct {
	Pk   string `dynamodbav:"pk" json:"pk"`
	Sk   string `dynamodbav:"sk" json:"sk"`
	Name string `dynamodbav:"name" json:"name"`
}
