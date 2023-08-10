package usecase

import (
	"context"

	"github.com/aqaurius6666/goserverless/internal/entity"
	"github.com/google/wire"
)

var UserSet = wire.NewSet(
	wire.Struct(new(UserUseCaseImpl), "*"),
	wire.Bind(new(UserUseCase), new(*UserUseCaseImpl)),
)

type UserUseCase interface {
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	CreateUser(ctx context.Context, id string) error
}

type UserRepository interface {
	GetUserById(ctx context.Context, id string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
}

type UserUseCaseImpl struct {
	UserRepo UserRepository
}

func (u *UserUseCaseImpl) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	return u.UserRepo.GetUserById(ctx, id)
}

func (u *UserUseCaseImpl) CreateUser(ctx context.Context, id string) error {
	return u.UserRepo.CreateUser(ctx, &entity.User{
		Sk:   id,
		Pk:   id,
		Name: id,
	})
}
