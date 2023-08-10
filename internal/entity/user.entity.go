package entity

type User struct {
	Pk   string `dynamodbav:"pk"`
	Sk   string `dynamodbav:"sk"`
	Name string `dynamodbav:"name"`
}
